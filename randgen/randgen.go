package randgen

import (
	"bytes"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/sdojjy/coprocessor_test/util"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/pingcap/tidb/types/parser_driver"
)

func GenData(sql []byte, count int) string {
	gofakeit.Seed(0)
	p := parser.New()
	stmts, warns, err := p.Parse(string(sql), "", "")
	if err != nil {
		log.Fatalf("parse schema file failed: %v", err)
		os.Exit(-1)
	}
	for _, w := range warns {
		log.Println("warn: " + w.Error())
	}

	var insert bytes.Buffer
	insert.WriteString("SET time_zone='+00:00';\n")
	for _, stmt := range stmts {
		createTable, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			log.Printf("statement %T is not create table", stmt)
			continue
		}

		config := genConfig{dataFunc: make(map[string]*indexTuple)}
		for _, col := range createTable.Cols {
			fillFlags(col)
			config.dataFunc[col.Name.Name.O] = &indexTuple{max: len(dataTypeRandGenFunMap[col.Tp.Tp]), current: 0}
		}

		var columnDataMap = map[string]*columnData{}
		for _, col := range createTable.Cols {
			columnDataMap[col.Name.Name.O] = &columnData{filled: false, data: make([]string, count)}
		}

		var columnMap = map[string]*ast.ColumnDef{}
		for _, col := range createTable.Cols {
			columnMap[col.Name.Name.O] = col
		}

		//fill column with constraint first
		for _, constraint := range createTable.Constraints {
			switch constraint.Tp {
			case ast.ConstraintPrimaryKey:
				cname := constraint.Keys[0].Column.Name.O
				cdata := columnDataMap[cname]
				cdata.filled = true
				col := columnMap[cname]
				funcSlice := dataTypeRandGenFunMap[col.Tp.Tp]
				cornerCaseStart := 0
				cornerCaseEnd := len(funcSlice) - 2
				unsigned := mysql.HasUnsignedFlag(col.Tp.Flag)

				var idStart = -count / 2
				for i := 0; i < count; i++ {
					if cornerCaseStart <= cornerCaseEnd {
						cdata.data[i] = funcSlice[cornerCaseStart](getTypeLen(col))
						cornerCaseStart++
					} else {
						if !unsigned {
							if idStart == 0 { //avoid duplicate key '1'
								idStart++
							}
							cdata.data[i] = fmt.Sprintf("'%d'", idStart)
						} else {
							cdata.data[i] = fmt.Sprintf("'%d'", idStart+count/2+1)
						}
						idStart++
					}
				}
			case ast.ConstraintIndex:
				valueSet := getUniqueValueSet(constraint, int(float32(count)*0.4), columnMap, config)
				data := valueSet.ToSlice()
				for index, item := range data {
					valueStr := item.(string)
					value := strings.Split(valueStr, "|")
					for keyIndex, key := range constraint.Keys {
						cname := key.Column.Name.O
						cdata := columnDataMap[cname]
						cdata.filled = true
						cdata.data[index] = value[keyIndex]
					}
				}
				dataLen := len(data)
				for i := 0; i < count-dataLen; i++ {
					item := data[rand.Intn(dataLen)]
					valueStr := item.(string)
					value := strings.Split(valueStr, "|")
					for keyIndex, key := range constraint.Keys {
						cname := key.Column.Name.O
						cdata := columnDataMap[cname]
						cdata.filled = true
						cdata.data[dataLen+i] = value[keyIndex]
					}
				}

			case ast.ConstraintUniq:
				valueSet := getUniqueValueSet(constraint, count, columnMap, config)
				for index, item := range valueSet.ToSlice() {
					valueStr := item.(string)
					value := strings.Split(valueStr, "|")
					for keyIndex, key := range constraint.Keys {
						cname := key.Column.Name.O
						cdata := columnDataMap[cname]
						cdata.filled = true
						cdata.data[index] = value[keyIndex]
					}
				}
			}
		}

		//fill normal column
		for _, col := range createTable.Cols {
			cdata := columnDataMap[col.Name.Name.O]
			if !cdata.filled {
				funcSlice := dataTypeRandGenFunMap[col.Tp.Tp]
				max := len(dataTypeRandGenFunMap[col.Tp.Tp]) - 1
				funcIndex := 0
				for i := 0; i < count; i++ {
					if funcIndex >= max {
						funcIndex = max
					}
					cdata.data[i] = funcSlice[funcIndex](getTypeLen(col))
					funcIndex++
				}
			}
		}

		var insertHeader = getInsertHeader(createTable)
		for i := 0; i < count; i++ {
			insert.WriteString(insertHeader)
			for index, col := range createTable.Cols {
				cdata := columnDataMap[col.Name.Name.O]
				insert.WriteString(cdata.data[i])
				if index < len(createTable.Cols)-1 {
					insert.WriteString(",")
				}
			}
			insert.WriteString(");\n")
		}
	}
	return insert.String()
}

func getUniqueValueSet(constraint *ast.Constraint, count int, columnMap map[string]*ast.ColumnDef, config genConfig) *util.Set {
	keys := constraint.Keys
	valueSet := util.NewSet()
	for valueSet.Size() < count {
		var uniqueColumnsData bytes.Buffer
		for _, key := range keys {
			cname := key.Column.Name.O
			col := columnMap[cname]
			funcSlice := dataTypeRandGenFunMap[col.Tp.Tp]
			funcIndex := config.getColFuncIndex(cname)
			if funcIndex < config.dataFunc[cname].max {
				config.dataFunc[cname].current = funcIndex + 1
			}
			uniqueColumnsData.WriteString(funcSlice[funcIndex](getTypeLen(col)))
			uniqueColumnsData.WriteString("|")
		}
		valueSet.Put(uniqueColumnsData.String())
	}
	return valueSet
}

//return insert into table_name(col_1, col_2....) values (
func getInsertHeader(createTable *ast.CreateTableStmt) string {
	var insert bytes.Buffer
	insert.WriteString("insert into `")
	insert.WriteString(createTable.Table.Name.O)
	insert.WriteString("` (")
	for index, col := range createTable.Cols {
		insert.WriteString(col.Name.Name.O)
		if index < len(createTable.Cols)-1 {
			insert.WriteString(",")
		}
	}
	insert.WriteString(") values ( ")
	return insert.String()
}

func fillFlags(col *ast.ColumnDef) {
	for _, opt := range col.Options {
		switch opt.Tp {
		case ast.ColumnOptionNotNull:
			col.Tp.Flag |= mysql.NotNullFlag
		case ast.ColumnOptionAutoIncrement:
			col.Tp.Flag |= mysql.AutoIncrementFlag
		}
	}
}

func getTypeLen(col *ast.ColumnDef) (int, int, *ast.ColumnDef) {
	defaultFlen, defaultDecimal := mysql.GetDefaultFieldLengthAndDecimal(col.Tp.Tp)
	flen := col.Tp.Flen
	if flen <= 0 {
		flen = defaultFlen
	}
	dec := col.Tp.Decimal
	if dec <= 0 {
		dec = defaultDecimal
	}
	return flen, dec, col
}

var dataTypeRandGenFunMap = make(map[byte][]func(int, int, *ast.ColumnDef) string)

func init() {
	dataTypeRandGenFunMap[mysql.TypeDecimal] = randNewDecimalFunArray
	dataTypeRandGenFunMap[mysql.TypeTiny] = randTinyFuncArray
	dataTypeRandGenFunMap[mysql.TypeShort] = randShortFuncArray
	dataTypeRandGenFunMap[mysql.TypeLong] = randLongFuncArray
	dataTypeRandGenFunMap[mysql.TypeFloat] = randFloatFuncArray
	dataTypeRandGenFunMap[mysql.TypeDouble] = randDoubleFuncArray
	dataTypeRandGenFunMap[mysql.TypeNull] = randNullFuncArray
	dataTypeRandGenFunMap[mysql.TypeTimestamp] = randTimeStampFuncArray
	dataTypeRandGenFunMap[mysql.TypeLonglong] = randLongLongFuncArray
	dataTypeRandGenFunMap[mysql.TypeInt24] = randInt24FuncArray
	dataTypeRandGenFunMap[mysql.TypeDate] = randDateFuncArray
	dataTypeRandGenFunMap[mysql.TypeDuration] = randTimeFuncArray

	dataTypeRandGenFunMap[mysql.TypeDatetime] = randDateTimeFuncArray
	dataTypeRandGenFunMap[mysql.TypeYear] = randYearFuncArray
	dataTypeRandGenFunMap[mysql.TypeNewDate] = randNewDateFuncArray
	dataTypeRandGenFunMap[mysql.TypeVarchar] = randStringFuncArray
	dataTypeRandGenFunMap[mysql.TypeBit] = randBitFuncArray

	dataTypeRandGenFunMap[mysql.TypeNewDecimal] = randNewDecimalFunArray

	dataTypeRandGenFunMap[mysql.TypeJSON] = randJSONFuncArray
	dataTypeRandGenFunMap[mysql.TypeTinyBlob] = randString100FuncArray
	dataTypeRandGenFunMap[mysql.TypeMediumBlob] = randString200FuncArray
	dataTypeRandGenFunMap[mysql.TypeLongBlob] = randString300FuncArray
	dataTypeRandGenFunMap[mysql.TypeBlob] = randString300FuncArray
	dataTypeRandGenFunMap[mysql.TypeVarString] = randStringFuncArray
	dataTypeRandGenFunMap[mysql.TypeString] = randStringFuncArray
}

func randTimestamp() string {
	return gofakeit.DateRange(time.Unix(0, 0), time.Now().Add(time.Hour*1440)).Format("2006-01-02 15:04:05.000")
}

func randString(minLen, maxLen int) string {
	var buf bytes.Buffer
	var size = 0
	if maxLen <= minLen {
		size = maxLen
	} else {
		size = minLen + rand.Intn(maxLen-minLen)
	}
	for i := 0; i < size; i++ {
		buf.WriteString(gofakeit.Letter())
	}
	return buf.String()
}

func randTime() string {
	return fmt.Sprintf("%d:%d:%d", rand.Intn(838)*2-838, rand.Intn(60), rand.Intn(60))
}

var digits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func randBit(len int) string {
	var buf bytes.Buffer
	for i := 0; i < len; i++ {
		buf.WriteByte(digits[rand.Intn(2)])
	}
	return buf.String()
}

var randTinyFuncArray = generateRangeFuncArr("0", strconv.FormatUint(math.MaxUint8, 10), strconv.FormatInt(math.MinInt8, 10), strconv.FormatInt(math.MaxInt8, 10))
var randShortFuncArray = generateRangeFuncArr("0", strconv.FormatUint(math.MaxUint16, 10), strconv.FormatInt(math.MinInt16, 10), strconv.FormatInt(math.MaxInt16, 10))
var randInt24FuncArray = generateRangeFuncArr("0", strconv.FormatUint(1<<21-1, 10), strconv.FormatInt(-1<<23, 10), strconv.FormatInt(1<<23-1, 10))
var randLongFuncArray = generateRangeFuncArr("0", strconv.FormatUint(math.MaxUint32, 10), strconv.FormatInt(math.MinInt32, 10), strconv.FormatInt(math.MaxInt32, 10))
var randLongLongFuncArray = generateRangeFuncArr("0", strconv.FormatUint(math.MaxUint64, 10), strconv.FormatInt(math.MinInt64, 10), strconv.FormatInt(math.MaxInt64, 10))

var randFloatFuncArray = generateRangeFuncArr("1.401298464324817070923729583289916131280e-45", "3.40282346638528859811704183484516925440e+38", "1.401298464324817070923729583289916131280e-45", "3.40282346638528859811704183484516925440e+38")
var randDoubleFuncArray = generateRangeFuncArr("4.940656458412465441765687928682213723651e-324", "1.797693134862315708145274237317043567981e+308", "4.940656458412465441765687928682213723651e-324", "1.797693134862315708145274237317043567981e+308")

func generateRangeFuncArr(unsignedMin, unsignedMax, min, max string) []func(int, int, *ast.ColumnDef) string {
	return []func(int, int, *ast.ColumnDef) string{
		func(flen int, dec int, col *ast.ColumnDef) string {
			if mysql.HasUnsignedFlag(col.Tp.Flag) {
				return fmt.Sprintf("'%s'", unsignedMax)
			} else {
				return fmt.Sprintf("'%s'", max)
			}
		},

		func(flen int, dec int, col *ast.ColumnDef) string {
			if mysql.HasUnsignedFlag(col.Tp.Flag) {
				return fmt.Sprintf("'%s'", unsignedMin)
			} else {
				return fmt.Sprintf("'%s'", min)
			}
		},

		func(flen int, dec int, col *ast.ColumnDef) string {
			switch col.Tp.Tp {
			case mysql.TypeDouble:
				return fmt.Sprintf("'%f'", gofakeit.Float32Range(float32(-500.9999), float32(500.9999)))
			case mysql.TypeFloat:
				return fmt.Sprintf("'%f'", gofakeit.Float32Range(float32(-500.9999), float32(500.9999)))
			//case mysql.TypeLonglong:
			//	max := math.Min(math.Pow10(flen)-1, 0x7ffffffffffff800) // avoid dealing with float -> int rounding
			//	return fmt.Sprintf("'%.0f'", gofakeit.Number(1, 2000))
			default:
				if mysql.HasUnsignedFlag(col.Tp.Flag) {
					//max, _ := strconv.ParseUint(unsignedMax, 10, 64)
					//maxf := math.Min(math.Pow10(flen)-1, float64(max))
					return fmt.Sprintf("'%d'", gofakeit.Number(1, 4000))
				} else {
					//max, _ := strconv.ParseInt(max, 10, 64)
					//maxf := math.Min(math.Pow10(flen)-1, float64(max))
					return fmt.Sprintf("'%d'", gofakeit.Number(-1000, 1000))
				}
			}
		},
	}
}

var randNullFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "null"
	},
}

var randTimeStampFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'1970-01-01 00:00:01.0000'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'2038-01-19 03:14:07.0000'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'000-00-00 00:00:00.0000'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%s'", randTimestamp())
	},
}

var randDateFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'1970-01-01'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'000-00-00'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%v'", gofakeit.Date().Format("2006-01-02"))
	},
}

var randTimeFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'-838:59:59.000000'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'838:59:59.000000'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return "'00:00:00.000000'"
	},
	func(flen int, dec int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%s'", randTime())
	},
}

var randDateTimeFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%v'", gofakeit.Date().Format("2006-01-02 15:04:05"))
	},
}

var randYearFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%d'", gofakeit.Number(1970, 2030))
	},
}

var randNewDateFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%v'", gofakeit.Date().Format("2006-01-02"))
	},
}

var randBitFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen int, dec int, col *ast.ColumnDef) string {
		return fmt.Sprintf("B'%s'", randBit(flen))
	},
}

var randString100FuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen, d int, col *ast.ColumnDef) string {
		if mysql.HasNotNullFlag(col.Tp.Flag) {
			return fmt.Sprintf("'%s'", randString(4, 100))
		}
		return "null"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return "''"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%s'", randString(4, 100))
	},
}
var randString200FuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen, d int, col *ast.ColumnDef) string {
		if mysql.HasNotNullFlag(col.Tp.Flag) {
			return fmt.Sprintf("'%s'", randString(4, 200))
		}
		return "null"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return "''"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%s'", randString(4, 200))
	},
}
var randString300FuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen, d int, col *ast.ColumnDef) string {
		if mysql.HasNotNullFlag(col.Tp.Flag) {
			return fmt.Sprintf("'%s'", randString(4, 300))
		}
		return "null"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return "''"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%s'", randString(4, 300))
	},
}
var randStringFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen, d int, col *ast.ColumnDef) string {
		if mysql.HasNotNullFlag(col.Tp.Flag) {
			return fmt.Sprintf("'%s'", randString(4, flen))
		}
		return "null"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return "''"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'%s'", randString(4, flen))
	},
}

var randJSONFuncArray = []func(int, int, *ast.ColumnDef) string{
	func(flen, d int, col *ast.ColumnDef) string {
		return "null"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return "'{}'"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return "'{\"arr\": []}'"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return "'{\"arr\": null}'"
	},
	func(flen, d int, col *ast.ColumnDef) string {
		return fmt.Sprintf("'{\"%s\":\"%s\", \"number\": %d, \"array\":[%d,%d], \"obj\":\"%s\"}'",
			randString(4, 20), randString(4, 20),
			rand.Intn(1234), rand.Intn(1234), rand.Intn(1234), randString(4, 40),
		)
	},
}

var randNewDecimalFunArray = []func(int, int, *ast.ColumnDef) string{
	minMaxNewDecimal(true),
	minMaxNewDecimal(false),
	normalRandNewDecimal,
}

func normalRandNewDecimal(m, d int, col *ast.ColumnDef) string {
	var buf bytes.Buffer
	if i := rand.Intn(2); i == 0 {
		buf.WriteByte('-')
	}
	for i := 0; i < m-d; i++ {
		buf.WriteByte(digits[rand.Intn(len(digits))])
	}

	if d > 0 {
		buf.WriteByte('.')
		for i := 0; i < d; i++ {
			buf.WriteByte(digits[rand.Intn(len(digits))])
		}
	}
	return buf.String()
}

func minMaxNewDecimal(negative bool) func(m, d int, col *ast.ColumnDef) string {
	return func(m, d int, col *ast.ColumnDef) string {
		var buf bytes.Buffer
		if negative {
			buf.WriteByte('-')
		}
		for i := 0; i < m-d; i++ {
			buf.WriteByte('9')
		}

		if d > 0 {
			buf.WriteByte('.')
			for i := 0; i < d; i++ {
				buf.WriteByte('9')
			}
		}
		return buf.String()
	}
}

type genConfig struct {
	dataFunc map[string]*indexTuple //max, current
}

func (config genConfig) getColFuncIndex(colName string) int {
	current := config.dataFunc[colName].current
	if current >= config.dataFunc[colName].max {
		return config.dataFunc[colName].max - 1
	} else {
		return current
	}
}

type indexTuple struct {
	max, current int
}

type columnData struct {
	filled bool
	data   []string
}
