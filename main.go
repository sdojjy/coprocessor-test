package main

import (
	"flag"
	"fmt"
	"github.com/pingcap/log"
	"github.com/sdojjy/coprocessor_test/diff"
	"github.com/sdojjy/coprocessor_test/randgen"
	"github.com/sdojjy/coprocessor_test/util"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var cliArgs struct {
	dsn1       string
	dsn2       string
	testDir    string
	strictMode bool

	genData    bool
	recordSize int
	dmlFile    string
}

func main() {
	flag.StringVar(&cliArgs.dsn1, "dsn1", "root:@tcp(127.0.0.1:4000)/test?charset=utf8", "mysql or tidb connect  string")
	flag.StringVar(&cliArgs.dsn2, "dsn2", "root:@tcp(127.0.0.1:3306)/test?charset=utf8", "mysql or tidb connect  string")
	flag.StringVar(&cliArgs.testDir, "test-dir", "./tests/", "the base directory of test cases")
	flag.BoolVar(&cliArgs.strictMode, "strict", true, "compare the sql result line by line with order")

	flag.BoolVar(&cliArgs.genData, "gen-data", false, "generate rand data base on the DML file")
	flag.StringVar(&cliArgs.dmlFile, "ddl", "./ddl.sql", "the create table statements file")
	flag.IntVar(&cliArgs.recordSize, "records", 1000, "record count")
	flag.Parse()

	if cliArgs.genData {
		if cliArgs.dmlFile == "" {
			log.Fatal("dml file path is required")
			os.Exit(-1)
		}
		sql, _ := ioutil.ReadFile(cliArgs.dmlFile)
		fmt.Println(randgen.GenData(sql, cliArgs.recordSize))
	} else {
		tidb, err1 := util.OpenDBWithRetry("mysql", cliArgs.dsn1)
		mysql, err2 := util.OpenDBWithRetry("mysql", cliArgs.dsn2)
		if err1 != nil || err2 != nil {
			log.Error("connect to db failed", zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
			os.Exit(-1)
		}
		diff.SqlDiff(cliArgs.testDir, tidb, mysql, cliArgs.strictMode)
	}
}
