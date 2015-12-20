package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	// Mysql Driver
	_ "github.com/go-sql-driver/mysql"
)

// TablesSizeSQL - XXX
var TablesSizeSQL = `
SELECT table_name as 'table',
	 table_schema as 'database',
	 table_rows as rows,
	 CONCAT(table_schema, '.', table_name) as full_name,
	 data_length as size,
	 index_length as indexes,
	(data_length + index_length) as total,
ROUND(index_length / data_length, 2) as index_fraction
FROM   information_schema.TABLES
WHERE table_schema NOT IN ('information_schema', 'performance_schema', 'mysql')
ORDER  BY data_length + index_length DESC;
`

// Collect - XXX
func Collect() error {
	serv := "root:123456@tcp/employees"

	// If user forgot the '/', add it
	if strings.HasSuffix(serv, ")") {
		serv = serv + "/"
	} else if serv == "localhost" {
		serv = ""
	}

	db, err := sql.Open("mysql", serv)
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query(`SHOW /*!50002 GLOBAL */ STATUS`)
	defer rows.Close()
	if err != nil {
		return err
	}

	fields := make(map[string]interface{})
	for rows.Next() {
		var name string
		var val interface{}

		err = rows.Scan(&name, &val)
		if err != nil {
			return err
		}

		var found bool

		// for _, mapped := range mappings {
		// 	if strings.HasPrefix(name, mapped.onServer) {
		// 		i, _ := strconv.Atoi(string(val.([]byte)))
		// 		fields[mapped.inExport+name[len(mapped.onServer):]] = i
		// 		found = true
		// 	}
		// }

		if found {
			continue
		}

		switch name {
		case "Queries":
			i, err := strconv.ParseInt(string(val.([]byte)), 10, 64)
			if err != nil {
				return err
			}

			fields["queries"] = i
		case "Slow_queries":
			i, err := strconv.ParseInt(string(val.([]byte)), 10, 64)
			if err != nil {
				return err
			}

			fields["slow_queries"] = i
		}
	}

	ConnRows, err := db.Query("SELECT user, sum(1) FROM INFORMATION_SCHEMA.PROCESSLIST GROUP BY user")
	defer ConnRows.Close()
	for ConnRows.Next() {
		var user string
		var connections int64

		err = ConnRows.Scan(&user, &connections)
		if err != nil {
			return err
		}

		fields := make(map[string]interface{})

		if err != nil {
			return err
		}
		fields["connections"] = connections

	}

	TableSizeRows, err := db.Query(TablesSizeSQL)
	defer TableSizeRows.Close()
	for TableSizeRows.Next() {
		// TABLES_SIZE_ROWS = ['table', 'database', 'rows',
		// 'full_name', 'size', 'indexes', 'total', 'index_fraction']

		var table string
		var database string
		var rows int64
		var fullName string
		var size int64
		var indexes int64
		var total int64
		var indexFraction float64

		// var val interface{}
		err = TableSizeRows.Scan(&table, &database, &rows, &fullName, &size, &indexes, &total, &indexFraction)
		if err != nil {
			return err
		}

	}

	fmt.Print(fields)

	return nil
}
