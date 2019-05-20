package client

import (
	"context"
	"fmt"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/distsql"
	"github.com/pingcap/tidb/tablecodec"
	"github.com/pingcap/tidb/util/ranger"
	"testing"
)

func TestGetRegions(t *testing.T) {
	c := getClient(t)
	tableInfo, _ := c.GetTableInfo("test", "a")

	// for record
	startKey, _ := tablecodec.GetTableHandleKeyRange(tableInfo.ID)
	region, perr, _ := c.PdClient.GetRegion(context.Background(), startKey)
	fmt.Printf("%v, %v,", region, perr)
	regionInfa, _ := GetRegions(TiDBConfig{Server: "127.0.0.1", Port: 10080}, "test", "a")
	fmt.Printf("%v", regionInfa)
	tableRegion, err := c.GetTableRegion(tableInfo.ID)
	fmt.Printf("%v, %v", tableRegion, err)
}

func TestGetTableInfo(t *testing.T) {
	dbInfo, err := getClient(t).GetTableInfo("mysql", "tidb")
	fmt.Printf("%v, %v", dbInfo, err)
}

func TestSendIndexScan(t *testing.T) {
	client := getClient(t)
	dbInfo, err := client.GetTableInfo("test", "a")
	if err != nil {
		t.Fatal("get db info failed")
	}
	SendCopIndexScanRequest(context.Background(), dbInfo, client)
}

func TestSendTableScan(t *testing.T) {
	client := getClient(t)
	dbInfo, err := client.GetTableInfo("test", "a")
	if err != nil {
		t.Fatal("get db info failed")
	}
	SendCopTableScanRequest(context.Background(), dbInfo, client)
}

func Test_range(t *testing.T) {
	full := ranger.FullRange()
	dbInfo, _ := getClient(t).GetTableInfo("mysql", "user")
	kv := distsql.TableRangesToKVRanges(dbInfo.ID, full, nil)
	fmt.Printf("%v", kv)
}

func getClient(t *testing.T) *ClusterClient {
	client, err := NewClient([]string{"127.0.0.1:2379"}, config.Security{})
	if err != nil {
		t.Fatal("create client failed")
	}
	return client
}
