package redis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/amonapp/amonagent/plugins"

	"gopkg.in/redis.v3"
)

func (p PerformanceStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	Gauges map[string]interface{} `json:"gauges"`
}

// RedisPerformanceFields - XXX
var RedisPerformanceFields = map[string]string{
	"aof_last_rewrite_time_sec": "aof.last_rewrite_time",
	"aof_rewrite_in_progress":   "aof.rewrite",
	"aof_current_size":          "aof.size",
	"aof_buffer_length":         "aof.buffer_length",

	// Network
	"connected_clients":    "net.clients",
	"connected_slaves":     "net.slaves",
	"rejected_connections": "net.rejected",

	// clients
	"blocked_clients":            "clients.blocked",
	"client_biggest_input_buf":   "clients.biggest_input_buf",
	"client_longest_output_list": "clients.longest_output_list",

	// Keys
	"evicted_keys": "keys.evicted",
	"expired_keys": "keys.expired",

	// stats
	"keyspace_hits":    "stats.keyspace_hits",
	"keyspace_misses":  "stats.keyspace_misses",
	"latest_fork_usec": "perf.latest_fork_usec",

	// pubsub
	"pubsub_channels": "pubsub.channels",
	"pubsub_patterns": "pubsub.patterns",

	// rdb
	"rdb_bgsave_in_progress":      "rdb.bgsave",
	"rdb_changes_since_last_save": "rdb.changes_since_last",
	"rdb_last_bgsave_time_sec":    "rdb.last_bgsave_time",

	// memory
	"mem_fragmentation_ratio": "mem.fragmentation_ratio",
	"used_memory":             "mem.used",
	"used_memory_lua":         "mem.lua",
	"used_memory_peak":        "mem.peak",
	"used_memory_rss":         "mem.rss",

	// replication
	"master_last_io_seconds_ago": "replication.last_io_seconds_ago",
	"master_sync_in_progress":    "replication.sync",
	"master_sync_left_bytes":     "replication.sync_left_bytes",
}

// Description - XXX
func (r *Redis) Description() string {
	return "Read metrics from a Redis server"
}

// Redis - XXX
type Redis struct {
}

// Collect - XXX
func (r *Redis) Collect() (interface{}, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	defer client.Close()

	val, err := client.Info().Result()
	if err != nil {
		fmt.Print(err.Error())
	}
	PerformanceStruct := PerformanceStruct{}

	gauges := make(map[string]interface{})
	split := strings.Split(val, "\n")
	for _, line := range split {
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		name := string(parts[0])
		metric, ok := RedisPerformanceFields[name]
		if !ok {
			continue
		}
		val := strings.TrimSpace(parts[1])

		gauges[metric] = val

	}
	PerformanceStruct.Gauges = gauges

	return PerformanceStruct, nil
}

func init() {
	plugins.Add("redis", func() plugins.Plugin {
		return &Redis{}
	})
}
