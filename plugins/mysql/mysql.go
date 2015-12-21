package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	// Mysql Driver
	_ "github.com/go-sql-driver/mysql"
)

// Counters - XXX
var Counters = map[string]string{
	"Connections":              "net.connections",
	"Innodb_data_reads":        "innodb.data_reads",
	"Innodb_data_writes":       "innodb.data_writes",
	"Innodb_os_log_fsyncs":     "innodb.os_log_fsyncs",
	"Innodb_row_lock_waits":    "innodb.row_lock_waits",
	"Innodb_row_lock_time":     "innodb.row_lock_time",
	"Innodb_mutex_spin_waits":  "innodb.mutex_spin_waits",
	"Innodb_mutex_spin_rounds": "innodb.mutex_spin_rounds",
	"Innodb_mutex_os_waits":    "innodb.mutex_os_waits",
	"Slow_queries":             "performance.slow_queries",
	"Questions":                "performance.questions",
	"Queries":                  "performance.queries",
	"Com_select":               "performance.com_select",
	"Com_insert":               "performance.com_insert",
	"Com_update":               "performance.com_update",
	"Com_delete":               "performance.com_delete",
	"Com_insert_select":        "performance.com_insert_select",
	"Com_update_multi":         "performance.com_update_multi",
	"Com_delete_multi":         "performance.com_delete_multi",
	"Com_replace_select":       "performance.com_replace_select",
	"Qcache_hits":              "performance.qcache_hits",
	"Created_tmp_tables":       "performance.created_tmp_tables",
	"Created_tmp_disk_tables":  "performance.created_tmp_disk_tables",
	"Created_tmp_files":        "performance.created_tmp_files",
}

// Gauges - XXX
var Gauges = map[string]string{
	"Max_used_connections":     "net.max_connections",
	"Open_tables":              "performance.open_tables",
	"Open_files":               "performance.open_files",
	"Table_locks_waited":       "performance.table_locks_waited",
	"Threads_connected":        "performance.threads_connected",
	"Innodb_current_row_locks": "innodb.current_row_locks",
}

// SlowQueriesSQL - XXX
var SlowQueriesSQL = `SELECT
	mysql.slow_log.query_time,
	mysql.slow_log.rows_sent,
	mysql.slow_log.rows_examined,
	mysql.slow_log.lock_time,
	mysql.slow_log.db,
	mysql.slow_log.sql_text AS query,
	mysql.slow_log.start_time
FROM
	mysql.slow_log
WHERE
	mysql.slow_log.query_time > 1
ORDER BY
	start_time DESC
LIMIT 30
`

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

// TableSizeData - XXX
type TableSizeData struct {
	Headers []string      `json:"headers"`
	Data    []interface{} `json:"data"`
}

// SlowQueriesData - XXX
type SlowQueriesData struct {
	Headers []string      `json:"headers"`
	Data    []interface{} `json:"data"`
}

func (p PerformanceStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	TableSizeData   `json:"tables_size"`
	SlowQueriesData `json:"slow_queries"`
	Gauges          map[string]interface{} `json:"gauges"`
	Counters        map[string]interface{} `json:"counters"`
}

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

	PerformanceStruct := PerformanceStruct{}

	rows, err := db.Query(`SHOW /*!50002 GLOBAL */ STATUS`)
	defer rows.Close()
	if err != nil {
		return err
	}

	counters := make(map[string]interface{})
	gauges := make(map[string]interface{})
	for rows.Next() {
		var name string
		var val interface{}

		err = rows.Scan(&name, &val)
		if err != nil {
			return err
		}

		for RawKey, FormatedKey := range Gauges {
			if name == RawKey {
				i, err := strconv.ParseInt(string(val.([]byte)), 10, 64)
				if err != nil {
					return err
				}
				gauges[FormatedKey] = i
			}

		}

		for RawKey, FormatedKey := range Counters {
			if name == RawKey {
				i, err := strconv.ParseInt(string(val.([]byte)), 10, 64)
				if err != nil {
					return err
				}
				counters[FormatedKey] = i
			}

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

		gauges["connections.active_connections"] = connections

	}
	PerformanceStruct.Gauges = gauges
	PerformanceStruct.Counters = counters

	TableSizeRows, err := db.Query(TablesSizeSQL)
	TableSizeHeaders := []string{"table", "database", "rows", "full_name", "size", "indexes", "total", "index_fraction"}
	TableSizeData := TableSizeData{Headers: TableSizeHeaders}

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

		err = TableSizeRows.Scan(&table, &database, &rows, &fullName, &size, &indexes, &total, &indexFraction)
		if err != nil {
			return err
		}
		fields := []interface{}{}
		fields = append(fields, table)
		fields = append(fields, database)
		fields = append(fields, rows)
		fields = append(fields, size)
		fields = append(fields, indexes)
		fields = append(fields, total)
		fields = append(fields, indexFraction)

		TableSizeData.Data = append(TableSizeData.Data, fields)

	}
	PerformanceStruct.TableSizeData = TableSizeData

	SlowQueriesRows, err := db.Query(SlowQueriesSQL)
	SlowQueriesHeader := []string{"query_time", "rows_sent", "rows_examined", "lock_time", "db", "query", "start_time"}
	SlowQueriesData := SlowQueriesData{Headers: SlowQueriesHeader}

	defer SlowQueriesRows.Close()
	for SlowQueriesRows.Next() {

		var queryTime string
		var rowsSent int64
		var rowsExamined int64
		var lockTime string
		var db string
		var query string
		var startTime string

		err = SlowQueriesRows.Scan(&queryTime, &rowsSent, &rowsExamined, &lockTime, &db, &query, &startTime)
		if err != nil {
			return err
		}

		fields := []interface{}{}
		fields = append(fields, queryTime)
		fields = append(fields, rowsSent)
		fields = append(fields, rowsExamined)
		fields = append(fields, lockTime)
		fields = append(fields, db)
		fields = append(fields, query)
		fields = append(fields, startTime)

		SlowQueriesData.Data = append(SlowQueriesData.Data, fields)

	}
	PerformanceStruct.SlowQueriesData = SlowQueriesData
	fmt.Print(PerformanceStruct)

	return nil
}
