package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pingcap/kvproto/pkg/coprocessor"
	"github.com/pingcap/kvproto/pkg/tikvpb"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/tidb/distsql"
	"github.com/pingcap/tidb/sessionctx/variable"
	"io/ioutil"
	"net/http"

	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tipb/go-tipb"
)

func send() {
	//req := &tikvrpc.Request{
	//	Type: tikvrpc.CmdCop,
	//	//Cop:,
	//	Context: kvrpcpb.Context{
	//		IsolationLevel: nil,
	//		Priority:       nil,
	//		NotFillCache:   true,
	//		HandleTime:     true,
	//		ScanDetail:     true,
	//	},
	//}
	//fmt.Println(req)

	tInfo := getDBInfo("127.0.0.1", "3306", "test", "student")
	column := tipb.ColumnInfo{}
	var columns []*tipb.ColumnInfo
	columns = append(columns, &column)

	dag := &tipb.DAGRequest{}
	executor := tipb.Executor{
		Tp: tipb.ExecType_TypeTableScan,
		TblScan: &tipb.TableScan{
			TableId: 1,
			Columns: columns,
			Desc:    false,
		},
	}
	dag.Executors = append(dag.Executors, &executor)

	data, err := dag.Marshal()
	cop := &coprocessor.Request{
		Tp:     kv.ReqTypeDAG,
		Data:   data,
		Ranges: nil,
	}
	request, err := (&distsql.RequestBuilder{}).SetKeyRanges(nil).
		SetDAGRequest(dag).
		SetDesc(false).
		SetKeepOrder(false).
		SetFromSessionVars(variable.NewSessionVars()).
		Build()
	fmt.Printf("%v%v", request, err)

	client := tikvpb.NewTikvClient(nil)
	fmt.Printf("%v", client)
	client.Coprocessor(context.Background(), cop)
}

func getDBInfo(tidbAddress, tidbPort, dbName, tableName string) *model.TableInfo {
	url := fmt.Sprintf("http://%s:%d/schema/%s", tidbAddress, tidbPort, fmt.Sprintf("%s/%s", dbName, tableName))
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil && err == nil {
			err = errClose
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		// Print response body directly if status is not ok.
		fmt.Println(string(body))
		panic(err)
	}

	db := &model.TableInfo{}
	json.Unmarshal(body, db)
	return db
}
