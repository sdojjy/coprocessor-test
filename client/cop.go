package client

import (
	"context"
	"encoding/json"
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
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tipb/go-tipb"
)

func SendCopRequest(ctx context.Context, dbInfo *model.TableInfo, c *ClusterClient) error {

	regionInfa, _ := GetRegions(TiDBConfig{Server: "127.0.0.1", Port: 10080}, "test", "a")

	//regionInfa, _ := c.GetTableRegion("test", "a")

	for _, region := range regionInfa.RecordRegions {

		_, peer, _ := c.PdClient.GetRegionByID(ctx, region.ID)

		bo := tikv.NewBackoffer(context.Background(), 20000)
		keyLocation, _ := c.RegionCache.LocateRegionByID(bo, region.ID)

		rpcContext, _ := c.RegionCache.GetRPCContext(bo, keyLocation.Region)

		columnInfo := model.ColumnsToProto(dbInfo.Columns, dbInfo.PKIsHandle)
		getColumnsTypes(dbInfo.Columns)
		dag := &tipb.DAGRequest{
			StartTs:       math.MaxInt64,
			OutputOffsets: []uint32{0, 1},
			//Flags:226,
			//TimeZoneName: "Asia/Chongqing",
			//TimeZoneOffset:28800,
			//EncodeType:tipb.EncodeType_TypeDefault,
		}
		//executor1 := tipb.Executor{
		//	Tp: tipb.ExecType_TypeTableScan,
		//	TblScan: &tipb.TableScan{
		//		TableId: dbInfo.ID,
		//		Columns: columnInfo,
		//		Desc:    false,
		//	},
		//}

		executor1 := tipb.Executor{
			Tp: tipb.ExecType_TypeIndexScan,
			IdxScan: &tipb.IndexScan{
				TableId: dbInfo.ID,
				IndexId: regionInfa.Indices[0].ID,
				Columns: columnInfo,
				Desc:    false,
			},
		}

		ss, _ := session.CreateSession(c.Storage)

		expr, _ := expression.ParseSimpleExprWithTableInfo(ss, "id=2 or id=1", dbInfo)
		exprPb := expression.NewPBConverter(c.Storage.GetClient(), &stmtctx.StatementContext{InSelectStmt: true}).ExprToPB(expr)
		executor2 := tipb.Executor{
			Tp: tipb.ExecType_TypeSelection,
			Selection: &tipb.Selection{
				Conditions: []*tipb.Expr{exprPb},
			},
		}
		dag.Executors = append(dag.Executors, &executor1, &executor2)

		full := ranger.FullIntRange(false)
		//_, first := splitRanges(full, true, false)
		//rangeO := GetKeyRange(c, region)
		keyRange := distsql.TableRangesToKVRanges(dbInfo.ID, full, nil)
		copRange := &copRanges{mid: keyRange}
		fmt.Println(copRange.String())
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
			parseResponse(resp, getColumnsTypes(dbInfo.Columns))
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

func getDBInfo(addr string, dbName, tableName string) (*model.TableInfo, error) {
	url := fmt.Sprintf("http://%s/schema/%s", addr, fmt.Sprintf("%s/%s", dbName, tableName))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil && err == nil {
			err = errClose
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// Print response body directly if status is not ok.
		fmt.Println(string(body))
		return nil, err
	}

	db := &model.TableInfo{}
	json.Unmarshal(body, db)
	return db, nil
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
