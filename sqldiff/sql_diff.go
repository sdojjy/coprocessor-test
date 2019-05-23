package sqldiff

import (
	"database/sql"
)

type SqlResultComparer interface {
	// execute query on two databases and compare the result
	CompareQuery(db1, db2 *sql.DB, query string) (bool, error, error)
}

type StandardComparer struct {
	Strict bool
}

func (c *StandardComparer) CompareQuery(db1, db2 *sql.DB, query string) (bool, error, error) {
	tidbResult, err1 := GetQueryResult(db1, query)
	mysqlResult, err2 := GetQueryResult(db2, query)
	if err1 != nil || err2 != nil {
		return false, err1, err2
	} else {
		// now compare the results
		equals := false
		if c.Strict {
			equals = mysqlResult.strictCompare(tidbResult)
		} else {
			equals = mysqlResult.nonOrderCompare(tidbResult)
		}
		return equals, nil, nil
	}
}
