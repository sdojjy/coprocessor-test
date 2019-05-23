package diffrunner

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sdojjy/coprocessor_test/sqldiff"
	"path"
	"strings"

	"github.com/pingcap/log"
	"github.com/pingcap/parser"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var failedQuery = map[string][]string{}
var currentCase string
var total = 0
var comparer sqldiff.SqlResultComparer

func SqlDiff(caseDir string, tidb, mysql *sql.DB, strictMode bool) {
	comparer = &sqldiff.StandardComparer{Strict: strictMode}

	for _, dir := range getAllTestCaseDir(caseDir) {
		log.Info("prepare test case data")
		prepare(tidb, mysql, path.Join(dir, "dml.sql"))
		log.Info("prepare data done")

		var failedCaseSql []string
		failedQuery[dir] = failedCaseSql
		cases := getTestCaseSqlFile(dir)
		currentCase = dir
		for _, caseFile := range cases {
			log.Info("executing sql case", zap.String("case=", caseFile))
			//split=""
			compareCase(tidb, mysql, caseFile)
			fmt.Println()
		}
	}

	var failedCaseCount int
	var failedSqlCount int
	for key, value := range failedQuery {
		if len(value) > 0 {
			failedCaseCount++
			fmt.Printf("\n%s failed queries:\n", key)
			for _, query := range failedQuery[key] {
				fmt.Printf("%s\n", strings.TrimSpace(query))
				failedSqlCount++
			}
		}
	}
	fmt.Printf("\nSummery: total=%d, success=%d, fail=%d, failed-sql-count=%d\n", len(failedQuery), len(failedQuery)-failedCaseCount, failedCaseCount, failedSqlCount)
	if failedCaseCount > 0 {
		os.Exit(-1)
	}
}

func executeFile(db1, db2 *sql.DB, file string, f func(*sql.DB, *sql.DB, string)) {
	sqlBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("open file failed", zap.String("file", file), zap.String("err", fmt.Sprintf("%v", err)))
	}
	p := parser.New()
	stmts, warns, err := p.Parse(string(sqlBytes), "", "")
	if err != nil {
		log.Fatal("parse schema file failed", zap.String("err", fmt.Sprintf("%v", err)))
		os.Exit(-1)
	}
	for _, w := range warns {
		log.Info("warn: " + w.Error())
	}

	for _, stmt := range stmts {
		f(db1, db2, stmt.Text())
	}
}

func compareCase(db1, db2 *sql.DB, file string) {
	executeFile(db1, db2, file, logComparer)
}

func prepare(db1, db2 *sql.DB, file string) {
	executeFile(db1, db2, file, prepareData)
}

func logComparer(db1, db2 *sql.DB, statement string) {
	total++
	ok, err1, err2 := comparer.CompareQuery(db1, db2, statement)
	if ok {
		log.Info("done with no difference", zap.String("query", statement))
	} else {
		failedQuery[currentCase] = append(failedQuery[currentCase], statement)
		if err1 != nil && err2 != nil {
			log.Error("execute failed on both db", zap.String("query", statement), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
		} else if err1 != nil || err2 != nil {
			log.Error("execute failed on both one db", zap.String("query", statement), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
		}
	}
}

func prepareData(db1, db2 *sql.DB, statement string) {
	result1, err1 := db1.Exec(statement)
	result2, err2 := db2.Exec(statement)
	if err1 != nil && err2 != nil {
		log.Error("execute failed on both db", zap.String("query", statement), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
	} else if err1 != nil || err2 != nil {
		log.Error("execute failed on both one db", zap.String("query", statement), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
	} else {
		row1, _ := result1.RowsAffected()
		row2, _ := result2.RowsAffected()
		if row1 != row2 {
			log.Warn("the affected row not the same", zap.String("query", statement), zap.String("rows", fmt.Sprintf("[%d,%d]", row1, row2)))
		}
	}
}
