package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"

	// Postgres Driver
	_ "github.com/lib/pq"

	"github.com/amonapp/amonagent/plugins"
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

// Config - XXX
type Config struct {
	Host string
	DB   string
}

// Start - XXX
func (p *PostgreSQL) Start() error {
	return nil
}

// Stop - XXX
func (p *PostgreSQL) Stop() {
}

var sampleConfig = `
#   Available config options:
#
#    {"host": "postgres://user:password@localhost:port/dbname"}
#
# Config location: /etc/opt/amonagent/plugins-enabled/postgresql.conf
`

// SampleConfig - XXX
func (p *PostgreSQL) SampleConfig() string {
	return sampleConfig
}

// SetConfigDefaults - XXX
func (p *PostgreSQL) SetConfigDefaults() error {
	configFile, err := plugins.UmarshalPluginConfig("postgresql")
	if err != nil {
		log.WithFields(log.Fields{"plugin": "postgresql", "error": err.Error()}).Error("Can't read config file")
	}
	var config Config
	decodeError := mapstructure.Decode(configFile, &config)
	if decodeError != nil {
		log.WithFields(log.Fields{"plugin": "postgresql", "error": decodeError.Error()}).Error("Can't decode config file")
	}

	u, _ := url.Parse(config.Host)
	config.DB = strings.Trim(u.Path, "/")

	p.Config = config

	return nil
}

// PostgreSQL - XXX
type PostgreSQL struct {
	Config Config
}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	TableSizeData    `json:"tables_size"`
	IndexHitRateData `json:"index_hit_rate"`
	SlowQueriesData  `json:"slow_queries"`
	Gauges           map[string]interface{} `json:"gauges"`
	Counters         map[string]interface{} `json:"counters"`
}

// Description - XXX
func (p *PostgreSQL) Description() string {
	return "Read metrics from a PostgreSQL server"
}

// Collect - XXX
func (p *PostgreSQL) Collect() (interface{}, error) {
	PerformanceStruct := PerformanceStruct{}
	p.SetConfigDefaults()

	db, err := sql.Open("postgres", p.Config.Host)
	if err != nil {
		log.Errorf("Can't connect to database': %v", err)
		return PerformanceStruct, err
	}

	defer db.Close()

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
	CountersQuery := fmt.Sprintf(DatabaseStatsSQL, strings.Join(counterColumns, ", "), p.Config.DB)
	CounterRows, errCounterRows := db.Query(CountersQuery)
	if errCounterRows == nil {
		defer CounterRows.Close()

		for CounterRows.Next() {
			err = CounterRows.Scan(dest...)
			if err != nil {
				log.Errorf("Can't get Counter stats': %v", err)

			}

			for i, val := range rawResult {
				key := counterColumns[i]
				counters[key] = string(val)
			}
		}
	}

	gauges := make(map[string]interface{})
	GaugesQuery := fmt.Sprintf(DatabaseStatsSQL, "numbackends", "amon")
	GaugeRows, errGaugeRows := db.Query(GaugesQuery)

	if errGaugeRows == nil {
		defer GaugeRows.Close()

		for GaugeRows.Next() {
			var connections int
			err = GaugeRows.Scan(&connections)
			if err != nil {
				log.Errorf("Can't get Gauges': %v", err)

			}
			gauges["connections"] = connections

		}

	}

	TableSizeRows, errTableSizeRows := db.Query(TableSizeSQL)
	TableSizeHeaders := []string{"name", "type", "size"}
	TableSizeData := TableSizeData{Headers: TableSizeHeaders}

	if errTableSizeRows == nil {
		defer TableSizeRows.Close()

		for TableSizeRows.Next() {
			// TABLES_SIZE_ROWS = ['name','type','size']
			var name string
			var Type string
			var size int64

			err = TableSizeRows.Scan(&name, &Type, &size)
			if err != nil {
				log.Errorf("Can't get Table size rows': %v", err)
			}
			fields := []interface{}{}
			fields = append(fields, name)
			fields = append(fields, Type)
			fields = append(fields, size)

			TableSizeData.Data = append(TableSizeData.Data, fields)

		}

		PerformanceStruct.TableSizeData = TableSizeData

	}

	IndexHitRows, errIndexHitRows := db.Query(IndexCacheHitRateSQL)
	IndexHitHeaders := []string{"table_name", "size", "reads",
		"cumulative_pct_reads", "index_hit_rate", "cache_hit_rate"}
	IndexHitData := IndexHitRateData{Headers: IndexHitHeaders}

	if errIndexHitRows == nil {
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
				log.Errorf("Can't get Index Hit Rate tables': %v", err)

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
	}

	SlowQueriesRows, errSlowQueriesRows := db.Query(SlowQueriesSQL)
	SlowQueriesHeaders := []string{"calls", "total", "per_call", "query"}
	SlowQueriesData := SlowQueriesData{Headers: SlowQueriesHeaders}

	if errSlowQueriesRows == nil {
		defer SlowQueriesRows.Close()
		for SlowQueriesRows.Next() {

			var Calls int64
			var Total float64
			var PerCall float64
			var Query string

			err = SlowQueriesRows.Scan(&Calls, &Total, &PerCall, &Query)
			if err != nil {
				log.Errorf("Can't get Slow Queries': %v", err)
			}
			fields := []interface{}{}
			fields = append(fields, Calls)
			fields = append(fields, Total)
			fields = append(fields, PerCall)
			fields = append(fields, Query)

			SlowQueriesData.Data = append(SlowQueriesData.Data, fields)

		}

		PerformanceStruct.SlowQueriesData = SlowQueriesData
	}

	PerformanceStruct.Counters = counters
	PerformanceStruct.Gauges = gauges

	return PerformanceStruct, nil
}

func init() {
	plugins.Add("postgresql", func() plugins.Plugin {
		return &PostgreSQL{}
	})
}
