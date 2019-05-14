package client

import (
	"fmt"
	"github.com/pingcap/kvproto/pkg/coprocessor"
	"github.com/pingcap/kvproto/pkg/kvrpcpb"
	"github.com/pingcap/kvproto/pkg/tikvpb"
	"github.com/pingcap/tidb/store/tikv/tikvrpc"

	"github.com/pingcap/tidb/distsql"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/sessionctx/variable"
	"github.com/pingcap/tipb/go-tipb"
)

func send() {
	req := &tikvrpc.Request{
		Type: tikvrpc.CmdCop,
		Cop: &coprocessor.Request{
			Tp:     kv.ReqTypeDAG,
			Data:   nil,
			Ranges: nil,
		},
		Context: kvrpcpb.Context{
			IsolationLevel: nil,
			Priority:       nil,
			NotFillCache:   true,
			HandleTime:     true,
			ScanDetail:     true,
		},
	}
	fmt.Println(req)

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

	request, err := (&distsql.RequestBuilder{}).SetKeyRanges(nil).
		SetDAGRequest(dag).
		SetDesc(false).
		SetKeepOrder(false).
		SetFromSessionVars(variable.NewSessionVars()).
		Build()
	fmt.Printf("%v%v", request, err)

	client := tikvpb.NewTikvClient(nil)
	fmt.Printf("%v", client)
}
