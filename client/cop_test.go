package client

import (
	"context"
	"fmt"
	"github.com/pingcap/tidb/distsql"
	"github.com/pingcap/tidb/util/ranger"
	"testing"
)

func Test_getDBInfo(t *testing.T) {
	dbInfo, err := getDBInfo("127.0.0.1:10080", "mysql", "tidb")
	fmt.Printf("%v, %v", dbInfo, err)
}

func Test_newTikvClient(t *testing.T) {
	tikvClient := newTikvClient("127.0.0.1:20160")
	fmt.Printf("%v", tikvClient)
}

func Test_send(t *testing.T) {

	tikvClient := newTikvClient("127.0.0.1:2379")
	dbInfo, _ := getDBInfo("127.0.0.1:10080", "mysql", "user")
	request, _ := buildCopRequest(context.Background(), dbInfo, tikvClient)
	send(tikvClient, "127.0.0.1:20160", dbInfo, request)
}

func Test_range(t *testing.T) {
	full := ranger.FullRange()
	dbInfo, _ := getDBInfo("127.0.0.1:10080", "mysql", "user")
	kv := distsql.TableRangesToKVRanges(dbInfo.ID, full, nil)
	fmt.Printf("%v", kv)
}
