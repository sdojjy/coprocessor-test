package diff

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"github.com/fatih/color"
	"github.com/pingcap/log"
	"github.com/pingcap/parser"
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
)

var failedQuery = map[string][]string{}
var currentCase string

var total = 0

func compareQuery(query string) {
	total++
	log.Info("comparing query result", zap.String("query", query))
	tidbResult, err1 := getQueryResult(db1, query)
	mysqlResult, err2 := getQueryResult(db2, query)
	if err1 != nil && err2 != nil {
		log.Error("execute failed on both db", zap.String("query", query), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
		failedQuery[currentCase] = append(failedQuery[currentCase], query)
	} else if err1 != nil || err2 != nil {
		log.Error("execute failed on both one db", zap.String("query", query), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
		failedQuery[currentCase] = append(failedQuery[currentCase], query)
	} else {
		// now compare the results
		equals := false
		if strict {
			equals = mysqlResult.strictCompare(tidbResult)
		} else {
			equals = mysqlResult.nonOrderCompare(tidbResult)
		}

		if equals {
			log.Info("done with no difference", zap.Int("result size", len(tidbResult.data)))
		} else {
			failedQuery[currentCase] = append(failedQuery[currentCase], query)
		}
	}
}

type SqlQueryResult struct {
	data        [][][]byte
	header      []string
	columnTypes []*sql.ColumnType
}

// readable query result like mysql shell client
func (result *SqlQueryResult) String() string {
	if result.data == nil || result.header == nil {
		return "no result"
	}

	// Calculate the max column length
	var colLength []int
	for _, c := range result.header {
		colLength = append(colLength, len(c))
	}
	for _, row := range result.data {
		for n, col := range row {
			if l := len(col); colLength[n] < l {
				colLength[n] = l
			}
		}
	}
	// The total length
	var total = len(result.header) - 1
	for index := range colLength {
		colLength[index] += 2 // Value will wrap with space
		total += colLength[index]
	}

	var lines []string
	var push = func(line string) {
		lines = append(lines, line)
	}

	// Write table header
	var header string
	for index, col := range result.header {
		length := colLength[index]
		padding := length - 1 - len(col)
		if index == 0 {
			header += "|"
		}
		header += " " + col + strings.Repeat(" ", padding) + "|"
	}
	splitLine := "+" + strings.Repeat("-", total) + "+"
	push(splitLine)
	push(header)
	push(splitLine)

	// Write rows data
	for _, row := range result.data {
		var line string
		for index, col := range row {
			length := colLength[index]
			padding := length - 1 - len(col)
			if index == 0 {
				line += "|"
			}
			line += " " + string(col) + strings.Repeat(" ", padding) + "|"
		}
		push(line)
	}
	push(splitLine)
	return strings.Join(lines, "\n")
}

//return true if two result is equals with order
func (result *SqlQueryResult) strictCompare(tidbResult *SqlQueryResult) bool {
	queryData1 := result.data
	queryData2 := tidbResult.data
	if len(queryData1) != len(queryData2) {
		log.Info("result length not equals", zap.Int("db1", len(queryData1)), zap.Int("db2", len(queryData2)))
		return false
	}

	for rowIndex, row := range queryData1 {
		if !compareRow(rowIndex, row, queryData2[rowIndex]) {
			printColorDiff(result.String(), tidbResult.String())
			return false
		}
	}
	return true
}

// compare two query result without ordered
func (result *SqlQueryResult) nonOrderCompare(result2 *SqlQueryResult) bool {
	queryData1 := result.data
	queryData2 := result2.data
	if len(queryData1) != len(queryData2) {
		log.Info("result length not equals", zap.Int("db1", len(queryData1)), zap.Int("db2", len(queryData2)))
		return false
	}
	var checkedRowArray = make([]bool, len(queryData1))
	for rowIndex, row := range queryData1 {
		hasOneEquals := false
		for checkIndex, checked := range checkedRowArray {
			if !checked {
				equals := compareRow(rowIndex, row, queryData2[checkIndex])
				if equals {
					checkedRowArray[checkIndex] = true
					hasOneEquals = true
					break
				}
			}
		}
		if !hasOneEquals {
			printColorDiff(result.String(), result2.String())
		}
	}
	return true
}

// compare two result row
func compareRow(rowIndex int, row [][]byte, row2 [][]byte) bool {
	//var line string
	for colIndex, col := range row {
		if len(row) != len(row2) {
			log.Info("result column length not equals", zap.Int("db1", len(row)), zap.Int("db2", len(row2)))
			return false
		}

		cv1 := string(col)
		cv2 := string(row2[colIndex])
		//driver not support column type
		if cv1 != cv2 {
			//maybe it's json
			if (strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) || (strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) {
				if !jsonEquals(cv1, cv2) {
					log.Info("result json value not equals", zap.Int("row", rowIndex+1), zap.Int("col", colIndex+1), zap.String("cv1", cv1), zap.String("cv2", cv2))
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}

func getQueryResult(db *sql.DB, query string) (*SqlQueryResult, error) {
	result, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	cols, err := result.Columns()
	if err != nil {
		return nil, err
	}
	types, err := result.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var allRows [][][]byte
	for result.Next() {
		var columns = make([][]byte, len(cols))
		var pointer = make([]interface{}, len(cols))
		for i := range columns {
			pointer[i] = &columns[i]
		}
		err := result.Scan(pointer...)
		if err != nil {
			return nil, err
		}
		allRows = append(allRows, columns)
	}
	queryResult := SqlQueryResult{data: allRows, header: cols, columnTypes: types}
	return &queryResult, nil
}

func getTestCaseSqlFile(dir string) []string {
	return loadItemsFromDir(dir, false)
}

func getAllTestCaseDir(dir string) []string {
	return loadItemsFromDir(dir, true)
}

// load all items from a directory, sub directory or file base on the loadDirectory parameter
func loadItemsFromDir(dir string, loadDirectory bool) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("read dir failed", zap.String("directory", dir), zap.String("err", fmt.Sprintf("%v", err)))
		os.Exit(-1)
	}
	var filesPaths []string
	for _, f := range files {
		if !loadDirectory && !f.IsDir() {
			if f.Name() != "dml.sql" {
				filesPaths = append(filesPaths, path.Join(dir, f.Name()))
			} else {
				log.Error("ignore dml.sql file")
			}
		} else if loadDirectory && f.IsDir() {
			filesPaths = append(filesPaths, path.Join(dir, f.Name()))
		} else {
			log.Info("ignore items", zap.String("name", f.Name()))
		}
	}
	//sort it
	sort.Strings(filesPaths)
	return filesPaths
}

var (
	db1 *sql.DB
	db2 *sql.DB
)

func executeFile(file string, f func(string)) {
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
		f(stmt.Text())
	}
}

func compareCase(file string) {
	executeFile(file, compareQuery)
}

func prepare(file string) {
	executeFile(file, prepareData)
}

func prepareData(statement string) {

	result1, err1 := (*sql.DB)(db1).Exec(statement)
	result2, err2 := (*sql.DB)(db2).Exec(statement)
	if err1 != nil && err2 != nil {
		log.Error("execute failed on both db", zap.String("query", statement), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
	} else if err1 != nil || err2 != nil {
		log.Error("execute failed on both one db", zap.String("query", statement), zap.String("err1", fmt.Sprintf("%v", err1)), zap.String("err2", fmt.Sprintf("%v", err2)))
	} else {
		row1, _ := result1.RowsAffected()
		row2, _ := result2.RowsAffected()
		if row1 != row2 {
			log.Error("the affected row not the same", zap.String("rows", fmt.Sprintf("[%d,%d]", row1, row2)))
		}
	}
}

var strict = true

func SqlDiff(caseDir string, tidb, mysql *sql.DB, strictMode bool) {
	strict = strictMode
	db1 = tidb
	db2 = mysql

	for _, dir := range getAllTestCaseDir(caseDir) {
		log.Info("prepare test case data")
		prepare(path.Join(dir, "dml.sql"))
		log.Info("prepare data done")

		var failedCaseSql []string
		failedQuery[dir] = failedCaseSql
		cases := getTestCaseSqlFile(dir)
		currentCase = dir
		for _, caseFile := range cases {
			log.Info("executing sql case", zap.String("case=", caseFile))
			//split=""
			compareCase(caseFile)
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

func jsonEquals(s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(o1, o2)
}

func printColorDiff(expect, actual string) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	patch := diffmatchpatch.New()
	diff := patch.DiffMain(expect, actual, false)
	var newExpectedContent, newActualResult bytes.Buffer
	for _, d := range diff {
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			newExpectedContent.WriteString(d.Text)
			newActualResult.WriteString(d.Text)
		case diffmatchpatch.DiffDelete:
			newExpectedContent.WriteString(red(d.Text))
		case diffmatchpatch.DiffInsert:
			newActualResult.WriteString(green(d.Text))
		}
	}
	fmt.Printf("Expected Result:\n%s\nActual Result:\n%s\n", newExpectedContent.String(), newActualResult.String())
}
