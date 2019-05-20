package client

import (
	"context"
	"fmt"
	"github.com/pingcap/kvproto/pkg/coprocessor"
	"github.com/pingcap/kvproto/pkg/kvrpcpb"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/distsql"
	"github.com/pingcap/tidb/expression"
	"github.com/pingcap/tidb/session"
	"github.com/pingcap/tidb/sessionctx/stmtctx"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/pingcap/tidb/store/tikv/tikvrpc"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/chunk"
	"github.com/pingcap/tidb/util/codec"
	"github.com/pingcap/tidb/util/ranger"
	"github.com/prometheus/common/log"
	"math"
	"time"

	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tipb/go-tipb"
)

func SendCopIndexScanRequest(ctx context.Context, tableInfo *model.TableInfo, c *ClusterClient) error {
	regionInfa, _ := GetRegions(TiDBConfig{Server: "127.0.0.1", Port: 10080}, "test", "a")

	for _, region := range regionInfa.RecordRegions {

		_, peer, _ := c.PdClient.GetRegionByID(ctx, region.ID)
		log.Info("region id: ", region.ID)

		bo := tikv.NewBackoffer(context.Background(), 20000)
		keyLocation, _ := c.RegionCache.LocateRegionByID(bo, region.ID)

		rpcContext, _ := c.RegionCache.GetRPCContext(bo, keyLocation.Region)

		collectSummary := true
		dag := &tipb.DAGRequest{
			StartTs:       math.MaxInt64,
			OutputOffsets: []uint32{0},
			//Flags:226,
			//TimeZoneName: "Asia/Chongqing",
			//TimeZoneOffset:28800,
			//EncodeType:tipb.EncodeType_TypeDefault,
			CollectExecutionSummaries: &collectSummary,
		}

		var idxColx []*model.ColumnInfo
		idxColx = append(idxColx, model.FindColumnInfo(tableInfo.Columns, "int_idx"))
		executor1 := tipb.Executor{
			Tp: tipb.ExecType_TypeIndexScan,
			IdxScan: &tipb.IndexScan{
				TableId: tableInfo.ID,
				IndexId: regionInfa.Indices[0].ID,
				Columns: model.ColumnsToProto(idxColx, tableInfo.PKIsHandle),
				Desc:    false,
			},
		}

		dag.Executors = append(dag.Executors, &executor1)

		//_, first := splitRanges(full, true, false)
		//rangeO := GetKeyRange(c, region)

		full := ranger.FullIntRange(false)
		keyRange, _ := distsql.IndexRangesToKVRanges(&stmtctx.StatementContext{InSelectStmt: true}, tableInfo.ID, tableInfo.Indices[0].ID, full, nil)
		copRange := &copRanges{mid: keyRange}
		data, _ := dag.Marshal()

		request := &tikvrpc.Request{
			Type: tikvrpc.CmdCop,
			Cop: &coprocessor.Request{
				Tp:     kv.ReqTypeDAG,
				Data:   data,
				Ranges: copRange.toPBRanges(),
			},

			Context: kvrpcpb.Context{
				//IsolationLevel: pbIsolationLevel(worker.req.IsolationLevel),
				//Priority:       kvPriorityToCommandPri(worker.req.Priority),
				//NotFillCache:   worker.req.NotFillCache,
				HandleTime: true,
				ScanDetail: true,
			},
		}

		if e := tikvrpc.SetContext(request, rpcContext.Meta, rpcContext.Peer); e != nil {
			log.Fatal(e)
		}

		addr, err := c.loadStoreAddr(ctx, bo, peer.StoreId)
		if err != nil {
			log.Fatal(err)
		}
		tikvResp, err := c.RpcClient.SendRequest(ctx, addr, request, readTimeout)
		if err != nil {
			return err
		}
		resp := tikvResp.Cop
		if resp.RegionError != nil || resp.OtherError != "" {
			log.Fatal("coprocessor grpc call failed", resp.RegionError, resp.OtherError)
		}

		if resp.Data != nil {
			parseResponse(resp, []*types.FieldType{types.NewFieldType(mysql.TypeLong)})
		}

	}
	return nil
}

func parseResponse(resp *coprocessor.Response, fs []*types.FieldType) {
	var data []byte = resp.Data
	selectResp := new(tipb.SelectResponse)
	err := selectResp.Unmarshal(data)
	if err != nil {
		log.Fatal("parse response failed", err)
	}
	if selectResp.Error != nil {
		log.Fatal("query failed ", selectResp.Error)
	}
	location, _ := time.LoadLocation("")
	chk := chunk.New(fs, 1024, 4096)
	r := selectResult{selectResp: selectResp, fieldTypes: fs, location: location}
	for r.respChkIdx < len(r.selectResp.Chunks) && len(r.selectResp.Chunks[r.respChkIdx].RowsData) != 0 {
		r.readRowsData(chk)
		it := chunk.NewIterator4Chunk(chk)
		for row := it.Begin(); row != it.End(); row = it.Next() {
			err = decodeTableRow(row, r.fieldTypes)
			if err != nil {
				log.Fatal("decode failed", err)
			}
		}
		chk.Reset()
		r.respChkIdx++
	}

}

func decodeTableRow(row chunk.Row, fs []*types.FieldType) error {
	for i, f := range fs {
		switch {
		case f.Tp == mysql.TypeVarchar:
			fmt.Printf("%s\t", row.GetString(i))
		case f.Tp == mysql.TypeLong:
			fmt.Printf("%d\t", row.GetInt64(i))
		}
	}
	fmt.Println()
	return nil
}

func SendCopTableScanRequest(ctx context.Context, tableInfo *model.TableInfo, c *ClusterClient) error {
	regionInfa, _ := GetRegions(TiDBConfig{Server: "127.0.0.1", Port: 10080}, "test", "a")

	for _, region := range regionInfa.RecordRegions {
		_, peer, _ := c.PdClient.GetRegionByID(ctx, region.ID)
		log.Info("region id: ", region.ID)

		bo := tikv.NewBackoffer(context.Background(), 20000)
		keyLocation, _ := c.RegionCache.LocateRegionByID(bo, region.ID)

		rpcContext, _ := c.RegionCache.GetRPCContext(bo, keyLocation.Region)

		collectSummary := true
		dag := &tipb.DAGRequest{
			StartTs:       math.MaxInt64,
			OutputOffsets: []uint32{0, 1, 2},
			//Flags:226,
			//TimeZoneName: "Asia/Chongqing",
			//TimeZoneOffset:28800,
			//EncodeType:tipb.EncodeType_TypeDefault,
			CollectExecutionSummaries: &collectSummary,
		}

		columnInfo := model.ColumnsToProto(tableInfo.Columns, tableInfo.PKIsHandle)
		executor1 := tipb.Executor{
			Tp: tipb.ExecType_TypeTableScan,
			TblScan: &tipb.TableScan{
				TableId: tableInfo.ID,
				Columns: columnInfo,
				Desc:    false,
			},
		}

		ss, _ := session.CreateSession(c.Storage)
		expr, _ := expression.ParseSimpleExprWithTableInfo(ss, "name='aaa' or name='name' ", tableInfo)
		exprPb := expression.NewPBConverter(c.Storage.GetClient(), &stmtctx.StatementContext{InSelectStmt: true}).ExprToPB(expr)
		executor2 := tipb.Executor{
			Tp: tipb.ExecType_TypeSelection,
			Selection: &tipb.Selection{
				Conditions: []*tipb.Expr{exprPb},
			},
		}

		//stmts, _, err :=parser.New().Parse("select count(id) from test.a", "", "")
		//parser.New().Parse(exprStr, "", "")
		//for _, warn := range warns {
		//	ctx.GetSessionVars().StmtCtx.AppendWarning(util.SyntaxWarn(warn))
		//}
		//if err != nil {
		//
		//}

		//groupBy := []expression.Expression{childCols[1]}
		//col := &expression.Column{Index: 0, RetType: types.NewFieldType(mysql.TypeLong)}
		//aggFunc := aggregation.NewAggFuncDesc(ss, "count", []expression.Expression{col}, false)
		//aggFuncs := []*aggregation.AggFuncDesc{aggFunc}
		//aggExprPb := aggregation.AggFuncToPBExpr(&stmtctx.StatementContext{InSelectStmt: true}, c.Storage.GetClient(), aggFunc)
		//pbConvert := expression.NewPBConverter(c.Storage.GetClient(), &stmtctx.StatementContext{InSelectStmt: true})
		//groupByExprPb := pbConvert.ExprToPB(col)
		//groupByExprPb := expression.GroupByItemToPB(&stmtctx.StatementContext{InSelectStmt: true}, c.Storage.GetClient(), col)

		//aggExpr, _ := expression.ParseSimpleExprWithTableInfo(ss, " count(id) ", tableInfo)
		//aggExprPb := expression.NewPBConverter(c.Storage.GetClient(),&stmtctx.StatementContext{InSelectStmt: true}).ExprToPB(aggExpr)
		//executor3 := tipb.Executor{
		//	Tp: tipb.ExecType_TypeAggregation,
		//	Aggregation: &tipb.Aggregation{
		//		AggFunc: []*tipb.Expr{aggExprPb},
		//		GroupBy: []*tipb.Expr{groupByExprPb},
		//	},
		//}
		//
		executor4 := tipb.Executor{
			Tp: tipb.ExecType_TypeLimit,
			Limit: &tipb.Limit{
				Limit: 1,
			},
		}

		//executor5 := tipb.Executor{
		//	Tp: tipb.ExecType_TypeTopN,
		//	TopN: &tipb.TopN{
		//		Limit: 2,
		//		//OrderBy:
		//	},
		//}

		dag.Executors = append(dag.Executors, &executor1, &executor2, &executor4)
		full := ranger.FullIntRange(false)
		keyRange := distsql.TableRangesToKVRanges(tableInfo.ID, full, nil)
		copRange := &copRanges{mid: keyRange}
		data, _ := dag.Marshal()

		request := &tikvrpc.Request{
			Type: tikvrpc.CmdCop,
			Cop: &coprocessor.Request{
				Tp:     kv.ReqTypeDAG,
				Data:   data,
				Ranges: copRange.toPBRanges(),
			},

			Context: kvrpcpb.Context{
				//IsolationLevel: pbIsolationLevel(worker.req.IsolationLevel),
				//Priority:       kvPriorityToCommandPri(worker.req.Priority),
				//NotFillCache:   worker.req.NotFillCache,
				HandleTime: true,
				ScanDetail: true,
			},
		}

		if e := tikvrpc.SetContext(request, rpcContext.Meta, rpcContext.Peer); e != nil {
			log.Fatal(e)
		}

		addr, err := c.loadStoreAddr(ctx, bo, peer.StoreId)
		if err != nil {
			log.Fatal(err)
		}
		tikvResp, err := c.RpcClient.SendRequest(ctx, addr, request, readTimeout)
		if err != nil {
			return err
		}
		resp := tikvResp.Cop
		if resp.RegionError != nil || resp.OtherError != "" {
			log.Fatal("coprocessor grpc call failed", resp.RegionError, resp.OtherError)
		}

		if resp.Data != nil {
			parseResponse(resp, getColumnsTypes(tableInfo.Columns))
		}
	}
	return nil
}

type selectResult struct {
	respChkIdx int
	fieldTypes []*types.FieldType

	selectResp *tipb.SelectResponse
	location   *time.Location
}

func (r *selectResult) readRowsData(chk *chunk.Chunk) (err error) {
	rowsData := r.selectResp.Chunks[r.respChkIdx].RowsData
	decoder := codec.NewDecoder(chk, r.location)
	for !chk.IsFull() && len(rowsData) > 0 {
		for i := 0; i < len(r.fieldTypes); i++ {
			rowsData, err = decoder.DecodeOne(rowsData, i, r.fieldTypes[i])
			if err != nil {
				return err
			}
		}
	}
	r.selectResp.Chunks[r.respChkIdx].RowsData = rowsData
	return nil
}

func getColumnsTypes(columns []*model.ColumnInfo) []*types.FieldType {
	colTypes := make([]*types.FieldType, 0, len(columns))
	for _, col := range columns {
		colTypes = append(colTypes, &col.FieldType)
	}
	return colTypes
}
