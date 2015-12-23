package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	// Postgres Driver
	_ "github.com/lib/pq"
)

// Counters - XXX
var Counters = map[string]string{
	"xact_commit":   "xact.commits",
	"xact_rollback": "xact.rollbacks",
	"blks_read":     "performance.disk_read",
	"blks_hit":      "performance.buffer_hit",
	"tup_returned":  "rows.returned",
	"tup_fetched":   "rows.fetched",
	"tup_inserted":  "rows.inserted",
	"tup_updated":   "rows.updated",
	"tup_deleted":   "rows.deleted",
}

// Gauges - XXX
var Gauges = map[string]string{
	"numbackends": "connections",
}

// DatabaseStatsSQL - XXX
var DatabaseStatsSQL = `SELECT %s FROM pg_stat_database WHERE datname IN ('%s')`

// MissingIndexesSQL - XXX
var MissingIndexesSQL = `SELECT
		  relname AS table,
		  CASE idx_scan
			WHEN 0 THEN 'Insufficient data'
			ELSE (100 * idx_scan / (seq_scan + idx_scan))::text
		  END percent_of_times_index_used,
		  n_live_tup rows_in_table
		FROM
		  pg_stat_user_tables
		WHERE
		  idx_scan > 0
		  AND (100 * idx_scan / (seq_scan + idx_scan)) < 95
		  AND n_live_tup >= 10000
		ORDER BY
		  n_live_tup DESC,
		  relname ASC
`

// SlowQueriesSQL - XXX
var SlowQueriesSQL = `SELECT * FROM
		(SELECT
				calls,
				round(total_time :: NUMERIC, 1) AS total,
				round(
					(total_time / calls) :: NUMERIC,
					3
				) AS per_call,
				regexp_replace(query, '[ \t\n]+', ' ', 'g') AS query
		FROM
				pg_stat_statements
		WHERE
		calls > 100
		ORDER BY
				total_time / calls DESC
		LIMIT 15
		) AS inner_table
		WHERE
		  per_call > 5
`

// https://gist.github.com/mattsoldo/3853455
var IndexCacheHitRateSQL = `
		-- Index hit rate
		WITH idx_hit_rate as (
		SELECT
		  relname as table_name,
		  n_live_tup,
		  round(100.0 * idx_scan / (seq_scan + idx_scan + 0.000001),2) as idx_hit_rate
		FROM pg_stat_user_tables
		ORDER BY n_live_tup DESC
		),

		-- Cache hit rate
		cache_hit_rate as (
		SELECT
		 relname as table_name,
		 heap_blks_read + heap_blks_hit as reads,
		 round(100.0 * sum (heap_blks_read + heap_blks_hit) over (ORDER BY heap_blks_read + heap_blks_hit DESC) / sum(heap_blks_read + heap_blks_hit + 0.000001) over (),4) as cumulative_pct_reads,
		 round(100.0 * heap_blks_hit / (heap_blks_hit + heap_blks_read + 0.000001),2) as cache_hit_rate
		FROM  pg_statio_user_tables
		WHERE heap_blks_hit + heap_blks_read > 0
		ORDER BY 2 DESC
		)

		SELECT
		  idx_hit_rate.table_name,
		  idx_hit_rate.n_live_tup as size,
		  cache_hit_rate.reads,
		  cache_hit_rate.cumulative_pct_reads,
		  idx_hit_rate.idx_hit_rate,
		  cache_hit_rate.cache_hit_rate
		FROM idx_hit_rate, cache_hit_rate
		WHERE idx_hit_rate.table_name = cache_hit_rate.table_name
		  AND cumulative_pct_reads < 100.0
		ORDER BY reads DESC;
`

// TableSizeSQL - XXX
var TableSizeSQL = `
		SELECT
			C .relname AS NAME,
			CASE
		WHEN C .relkind = 'r' THEN
			'table'
		ELSE
			'index'
		END AS TYPE,
		 pg_table_size(C .oid) AS SIZE
		FROM
			pg_class C
		LEFT JOIN pg_namespace n ON (n.oid = C .relnamespace)
		WHERE
			n.nspname NOT IN (
				'pg_catalog',
				'information_schema'
			)
		AND n.nspname !~ '^pg_toast'
		AND C .relkind IN ('r', 'i')
		ORDER BY
			pg_table_size (C .oid) DESC,
			NAME ASC
`

// IndexHitRateData - XXX
type IndexHitRateData struct {
	Headers []string      `json:"headers"`
	Data    []interface{} `json:"data"`
}

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
	TableSizeData    `json:"tables_size"`
	IndexHitRateData `json:"index_hit_rate"`
	SlowQueriesData  `json:"slow_queries"`
	Gauges           map[string]interface{} `json:"gauges"`
	Counters         map[string]interface{} `json:"counters"`
}

// Collect - XXX
func Collect() error {
	var serv string
	serv = "postgres://postgres:123456@localhost/amon"

	// If user forgot the '/', add it
	if strings.HasSuffix(serv, ")") {
		serv = serv + "/"
	} else if serv == "localhost" {
		serv = ""
	}

	db, err := sql.Open("postgres", serv)
	if err != nil {
		return err
	}

	defer db.Close()
	PerformanceStruct := PerformanceStruct{}

	rawResult := make([][]byte, len(Counters))
	counters := make(map[string]interface{})
	dest := make([]interface{}, len(Counters)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}
	var counterColumns []string
	for key := range Counters {
		counterColumns = append(counterColumns, key)
	}
	CountersQuery := fmt.Sprintf(DatabaseStatsSQL, strings.Join(counterColumns, ", "), "amon")
	CounterRows, err := db.Query(CountersQuery)
	defer CounterRows.Close()
	for CounterRows.Next() {
		err = CounterRows.Scan(dest...)
		if err != nil {
			return err
		}

		for i, val := range rawResult {
			key := counterColumns[i]
			counters[key] = string(val)
		}
	}

	gauges := make(map[string]interface{})
	GaugesQuery := fmt.Sprintf(DatabaseStatsSQL, "numbackends", "amon")
	GaugeRows, err := db.Query(GaugesQuery)
	defer GaugeRows.Close()
	for GaugeRows.Next() {
		var connections int
		err = GaugeRows.Scan(&connections)
		if err != nil {
			return err
		}
		gauges["connections"] = connections

	}

	TableSizeRows, err := db.Query(TableSizeSQL)
	TableSizeHeaders := []string{"name", "type", "size"}
	TableSizeData := TableSizeData{Headers: TableSizeHeaders}

	defer TableSizeRows.Close()
	for TableSizeRows.Next() {
		// TABLES_SIZE_ROWS = ['name','type','size']
		var name string
		var Type string
		var size int64

		err = TableSizeRows.Scan(&name, &Type, &size)
		if err != nil {
			return err
		}
		fields := []interface{}{}
		fields = append(fields, name)
		fields = append(fields, Type)
		fields = append(fields, size)

		TableSizeData.Data = append(TableSizeData.Data, fields)

	}

	PerformanceStruct.TableSizeData = TableSizeData

	IndexHitRows, err := db.Query(IndexCacheHitRateSQL)
	IndexHitHeaders := []string{"table_name", "size", "reads",
		"cumulative_pct_reads", "index_hit_rate", "cache_hit_rate"}
	IndexHitData := IndexHitRateData{Headers: IndexHitHeaders}
	defer IndexHitRows.Close()

	for IndexHitRows.Next() {
		// INDEX_CACHE_HIT_RATE_ROWS = ['table_name','size','reads',
		// 'cumulative_pct_reads', 'index_hit_rate', 'cache_hit_rate']

		var TableName string
		var Size string
		var Read int64
		var CumulativeReads float64
		var IndexHitRate float64
		var CacheHitRate float64

		err = IndexHitRows.Scan(&TableName, &Size, &Read, &CumulativeReads, &IndexHitRate, &CacheHitRate)
		if err != nil {
			return err
		}
		fields := []interface{}{}
		fields = append(fields, TableName)
		fields = append(fields, Size)
		fields = append(fields, Read)
		fields = append(fields, CumulativeReads)
		fields = append(fields, IndexHitRate)
		fields = append(fields, CacheHitRate)

		IndexHitData.Data = append(IndexHitData.Data, fields)

	}
	PerformanceStruct.IndexHitRateData = IndexHitData

	SlowQueriesRows, err := db.Query(SlowQueriesSQL)
	SlowQueriesHeaders := []string{"calls", "total", "per_call", "query"}
	SlowQueriesData := SlowQueriesData{Headers: SlowQueriesHeaders}

	defer SlowQueriesRows.Close()
	for SlowQueriesRows.Next() {

		var Calls int64
		var Total int64
		var PerCall float64
		var Query string

		err = SlowQueriesRows.Scan(&Calls, &Total, &PerCall, &Query)
		if err != nil {
			return err
		}
		fields := []interface{}{}
		fields = append(fields, Calls)
		fields = append(fields, Total)
		fields = append(fields, PerCall)
		fields = append(fields, Query)

		SlowQueriesData.Data = append(SlowQueriesData.Data, fields)

	}
	PerformanceStruct.SlowQueriesData = SlowQueriesData

	PerformanceStruct.Counters = counters
	PerformanceStruct.Gauges = gauges

	fmt.Print(PerformanceStruct)

	return nil
}
