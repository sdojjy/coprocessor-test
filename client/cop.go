package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pingcap/kvproto/pkg/coprocessor"
	"github.com/pingcap/kvproto/pkg/kvrpcpb"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/pingcap/tidb/store/tikv/tikvrpc"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/chunk"
	"github.com/pingcap/tidb/util/codec"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tipb/go-tipb"

	"github.com/pingcap/tidb/config"
)

func send(clusterClient *ClusterClient, tikvAddr string, dbInfo *model.TableInfo, copRequest *tikvrpc.Request) error {

	//tikvResp, err := clusterClient.RpcClient.SendRequest(context.Background(), tikvAddr, copRequest, time.Duration(time.Second * 3))
	//if err != nil {
	//	log.Fatal("coprocessor grpc call failed", err)
	//}
	//resp := tikvResp.Cop
	//if resp.RegionError != nil || resp.OtherError != "" {
	//	log.Fatal("coprocessor grpc call failed", resp.RegionError, resp.OtherError)
	//}
	//parseResponse(resp, getColumnsTypes(dbInfo.Columns))
	return nil
}

func newTikvClient(addr string) *ClusterClient {
	c, err := NewClient([]string{addr}, config.Security{})
	if err != nil {
		log.Error("Create client failed", fmt.Sprintf("err=%v", err))
		return nil
	}
	return c
}

func buildCopRequest(ctx context.Context, dbInfo *model.TableInfo, c *ClusterClient) (*tikvrpc.Request, error) {

	//c.ScanByRegion(ctx, dbInfo.ID, )

	regionInfa, _ := GetRegions(TiDBConfig{Server: "127.0.0.1", Port: 10080}, "mysql", "user")

	for _, region := range regionInfa.RegionIDs {

		r, peer, _ := c.PdClient.GetRegionByID(ctx, region)

		columnInfo := model.ColumnsToProto(dbInfo.Columns, dbInfo.PKIsHandle)
		getColumnsTypes(dbInfo.Columns)
		dag := &tipb.DAGRequest{}
		executor := tipb.Executor{
			Tp: tipb.ExecType_TypeTableScan,
			TblScan: &tipb.TableScan{
				TableId: dbInfo.ID,
				Columns: columnInfo,
				Desc:    false,
			},
		}
		dag.Executors = append(dag.Executors, &executor)

		//full := ranger.FullRange()
		rangeO := GetKeyRange(c, region)
		keyRange := []kv.KeyRange{*rangeO} //distsql.TableRangesToKVRanges(dbInfo.ID, full, nil)
		copRange := &copRanges{mid: keyRange}
		data, _ := dag.Marshal()

		var ctxt = kvrpcpb.Context{
			//IsolationLevel: pbIsolationLevel(worker.req.IsolationLevel),
			//Priority:       kvPriorityToCommandPri(worker.req.Priority),
			//NotFillCache:   worker.req.NotFillCache,
			HandleTime: true,
			ScanDetail: true,

			Peer:        peer,
			RegionId:    region,
			RegionEpoch: r.RegionEpoch,
		}

		cop := &coprocessor.Request{
			Context: &ctxt,
			Tp:      kv.ReqTypeDAG,
			Data:    data,
			Ranges:  copRange.toPBRanges(),
		}

		request := &tikvrpc.Request{
			Type:    tikvrpc.CmdCop,
			Cop:     cop,
			Context: ctxt,
		}

		bo := tikv.NewBackoffer(context.Background(), 20000)
		addr, err := c.loadStoreAddr(ctx, bo, peer.StoreId)
		if err != nil {
			log.Fatal(err)
		}
		tikvResp, err := c.RpcClient.SendRequest(ctx, addr, request, readTimeout)
		if err != nil {
			log.Fatal(err)
		}
		resp := tikvResp.Cop
		if resp.RegionError != nil || resp.OtherError != "" {
			log.Fatal("coprocessor grpc call failed", resp.RegionError, resp.OtherError)
		}

		if resp.Data != nil {
			parseResponse(resp, getColumnsTypes(dbInfo.Columns))
		}

	}
	return nil, nil
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
	for len(r.selectResp.Chunks[r.respChkIdx].RowsData) != 0 {
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
			row.GetString(i)
		case f.Tp == mysql.TypeLong:
			row.GetInt64(i)
		}
	}
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

	//ctx        sessionctx.Context
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
