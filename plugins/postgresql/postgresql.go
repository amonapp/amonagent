package postgresql

import "encoding/json"

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
