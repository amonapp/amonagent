package mongodb

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	// MongoDB Driver
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/amonapp/amonagent/plugins"
)

var localhost = &url.URL{Host: "127.0.0.1:27017"}

// Server - XXX
type Server struct {
	URL        *url.URL
	Session    *mgo.Session
	lastResult *ServerStatus
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
	TableSizeData   `json:"tables_size"`
	SlowQueriesData `json:"slow_queries"`
	Gauges          map[string]interface{} `json:"gauges"`
}

// DefaultStats - XXX
var DefaultStats = map[string]string{
	"operations.inserts_per_sec":  "Insert",
	"operations.queries_per_sec":  "Query",
	"operations.updates_per_sec":  "Update",
	"operations.deletes_per_sec":  "Delete",
	"operations.getmores_per_sec": "GetMore",
	"operations.commands_per_sec": "Command",
	"operations.flushes_per_sec":  "Flushes",
	"memory.vsize_megabytes":      "Virtual",
	"memory.resident_megabytes":   "Resident",
	"queued.reads":                "QueuedReaders",
	"queued.writes":               "QueuedWriters",
	"active.reads":                "ActiveReaders",
	"active.writes":               "ActiveWriters",
	"net.bytes_in":                "NetIn",
	"net.bytes_out":               "NetOut",
	"open_connections":            "NumConnections",
}

// DefaultReplStats - XXX
var DefaultReplStats = map[string]string{
	"replica.inserts_per_sec":  "InsertR",
	"replica.queries_per_sec":  "QueryR",
	"replica.updates_per_sec":  "UpdateR",
	"replica.deletes_per_sec":  "DeleteR",
	"replica.getmores_per_sec": "GetMoreR",
	"replica.commands_per_sec": "CommandR",
}

// MmapStats - XXX
var MmapStats = map[string]string{
	"mapped_megabytes":               "Mapped",
	"non-mapped_megabytes":           "NonMapped",
	"operations.page_faults_per_sec": "Faults",
}

// WiredTigerStats - XXX
var WiredTigerStats = map[string]string{
	"percent_cache_dirty": "CacheDirtyPercent",
	"percent_cache_used":  "CacheUsedPercent",
}

// CollectionStats - XXX
// COLLECTION_ROWS = ['count','ns','avgObjSize', 'totalIndexSize', 'indexSizes', 'size']
type CollectionStats struct {
	Count          int64            `json:"number_of_documents"`
	Ns             string           `json:"ns"`
	AvgObjSize     int64            `json:"avgObjSize"`
	TotalIndexSize int64            `json:"total_index_size"`
	StorageSize    int64            `json:"storage_size"`
	IndexSizes     map[string]int64 `json:"index_sizes"`
	Size           int64            `json:"size"`
}

// // CollectSlowQueries - XXX
// func CollectSlowQueries(server *Server, perf *PerformanceStruct) error {
// 	//
// 	// params = {"millis": { "$gt" : slowms }}
// 	// 	performance = db['system.profile']\
// 	// 	.find(params)\
// 	// 	.sort("ts", pymongo.DESCENDING)\
// 	// 	.limit(10)
// 	db := strings.Replace(server.URL.Path, "/", "", -1) // remove slash from Path
// 	result := []bson.M{}
//
// 	params := bson.M{"millis": bson.M{"$gt": 10}}
// 	c := server.Session.DB(db).C("system.profile")
// 	err := c.Find(params).All(&result)
// 	if err != nil {
// 		return err
// 	}
// 	for _, r := range result {
// 		pluginLogger.Errorfln(r)
// 		pluginLogger.Errorfln("-----")
// 	}
// 	// pluginLogger.Errorfln(result)
// 	return nil
// }

// CollectCollectionSize - XXX
func CollectCollectionSize(server *Server, perf *PerformanceStruct) error {
	TableSizeHeaders := []string{"count", "ns", "avgObjSize", "totalIndexSize", "storageSize", "indexSizes", "size"}
	TableSizeData := TableSizeData{Headers: TableSizeHeaders}

	db := strings.Replace(server.URL.Path, "/", "", -1) // remove slash from Path
	collections, err := server.Session.DB(db).CollectionNames()
	if err != nil {
		return err
	}
	for _, col := range collections {

		result := bson.M{}
		err := server.Session.DB(db).Run(bson.D{{"collstats", col}}, &result)

		if err != nil {
			log.WithFields(log.Fields{"plugin": "mongodb", "error": err.Error()}).Error("Can't get stats for collection")
		}
		var CollectionResult CollectionStats
		decodeError := mapstructure.Decode(result, &CollectionResult)
		if decodeError != nil {
			log.WithFields(log.Fields{"plugin": "mongodb", "error": decodeError.Error()}).Error("Can't decode collection stats")
		}

		TableSizeData.Data = append(TableSizeData.Data, CollectionResult)
	}

	perf.TableSizeData = TableSizeData

	return nil
}

// GetSession - XXX
func GetSession(server *Server) error {

	if server.Session == nil {

		dialInfo := &mgo.DialInfo{
			Addrs:    []string{server.URL.Host},
			Database: server.URL.Path,
		}
		dialInfo.Timeout = 5 * time.Second
		if server.URL.User != nil {
			password, _ := server.URL.User.Password()
			dialInfo.Username = server.URL.User.Username()
			dialInfo.Password = password
		}

		session, connectionError := mgo.DialWithInfo(dialInfo)
		if connectionError != nil {
			return fmt.Errorf("Unable to connect to URL (%s), %s\n", server.URL.Host, connectionError.Error())
		}
		server.Session = session
		server.lastResult = nil

		server.Session.SetMode(mgo.Eventual, true)
		server.Session.SetSocketTimeout(0)

	}

	return nil
}

// CollectGauges - XXX
func CollectGauges(server *Server, perf *PerformanceStruct) error {
	db := strings.Replace(server.URL.Path, "/", "", -1) // remove slash from Path
	result := &ServerStatus{}
	err := server.Session.DB(db).Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, result)
	if err != nil {
		return err
	}
	defer func() {
		server.lastResult = result
	}()

	result.SampleTime = time.Now()

	if server.lastResult != nil && result != nil {
		duration := result.SampleTime.Sub(server.lastResult.SampleTime)
		durationInSeconds := int64(duration.Seconds())
		if durationInSeconds == 0 {
			durationInSeconds = 1
		}

		data := NewStatLine(*server.lastResult, *result, server.URL.Host, true, durationInSeconds)

		statLine := reflect.ValueOf(data).Elem()
		storageEngine := statLine.FieldByName("StorageEngine").Interface()
		// nodeType := statLine.FieldByName("NodeType").Interface()

		gauges := make(map[string]interface{})
		for key, value := range DefaultStats {
			val := statLine.FieldByName(value).Interface()
			gauges[key] = val
		}

		if storageEngine == "mmapv1" {
			for key, value := range MmapStats {
				val := statLine.FieldByName(value).Interface()
				gauges[key] = val
			}
		} else if storageEngine == "wiredTiger" {
			for key, value := range WiredTigerStats {
				val := statLine.FieldByName(value).Interface()
				percentVal := fmt.Sprintf("%.1f", val.(float64)*100)
				floatVal, _ := strconv.ParseFloat(percentVal, 64)
				gauges[key] = floatVal
			}
		}

		perf.Gauges = gauges

	}

	return nil
}

// MongoDB - XXX
type MongoDB struct {
	Config Config
}

// Config - XXX
type Config struct {
	URI string
}

// Start - XXX
func (m *MongoDB) Start() error {
	return nil
}

// Stop - XXX
func (m *MongoDB) Stop() {
}

var sampleConfig = `
#   Available config options:
#
#    {"uri": "mongodb://username:password@host:port/database"}
#
# Config location: /etc/opt/amonagent/plugins-enabled/mongodb.conf
`

// SampleConfig - XXX
func (m *MongoDB) SampleConfig() string {
	return sampleConfig
}

// SetConfigDefaults - XXX
func (m *MongoDB) SetConfigDefaults() error {
	configFile, err := plugins.UmarshalPluginConfig("mongodb")
	if err != nil {
		log.WithFields(log.Fields{"plugin": "mongodb", "error": err.Error()}).Error("Can't read config file")
	}
	var config Config
	decodeError := mapstructure.Decode(configFile, &config)
	if decodeError != nil {
		log.WithFields(log.Fields{"plugin": "mongodb", "error": decodeError.Error()}).Error("Can't decode config file")
	}

	m.Config = config

	return nil
}

// Description - XXX
func (m *MongoDB) Description() string {
	return "Read metrics from a MongoDB server"
}

// Collect - XXX
func (m *MongoDB) Collect() (interface{}, error) {
	m.SetConfigDefaults()
	PerformanceStruct := PerformanceStruct{}

	url, err := url.Parse(m.Config.URI)
	if err != nil {
		log.Errorf("Can't parse Mongo URI': %v", err)
		return PerformanceStruct, err
	}

	server := Server{URL: url}
	sessionError := GetSession(&server)
	if sessionError != nil {
		log.Errorf("Can't connect to server': %v", sessionError)
		return PerformanceStruct, err
	}
	CollectGauges(&server, &PerformanceStruct)
	time.Sleep(time.Duration(1) * time.Second)
	CollectGauges(&server, &PerformanceStruct)

	CollectCollectionSize(&server, &PerformanceStruct)
	// // CollectSlowQueries(&server, &PerformanceStruct)

	if server.Session != nil {
		defer server.Session.Close()
	}

	return PerformanceStruct, nil
}

func init() {
	plugins.Add("mongodb", func() plugins.Plugin {
		return &MongoDB{}
	})
}
