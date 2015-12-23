package mongodb

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/influxdb/telegraf/plugins/mongodb"

	// MongoDB Driver

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Server - XXX
type Server struct {
	URL        *url.URL
	Session    *mgo.Session
	lastResult *mongodb.ServerStatus
}

var localhost = &url.URL{Host: "127.0.0.1:27017"}

// // Connect - XXX
// func Connect(server *Server) error {
// 	if server.Session == nil {
// 		var dialAddrs []string
// 		if server.Url.User != nil {
// 			dialAddrs = []string{server.URL.String()}
// 		} else {
// 			dialAddrs = []string{server.URL.Host}
// 		}
// 		dialInfo, err := mgo.ParseURL(dialAddrs[0])
// 		if err != nil {
// 			return fmt.Errorf("Unable to parse URL (%s), %s\n", dialAddrs[0], err.Error())
// 		}
// 		dialInfo.Direct = true
// 		dialInfo.Timeout = time.Duration(10) * time.Second
//
// 		if m.Ssl.Enabled {
// 			tlsConfig := &tls.Config{}
// 			if len(m.Ssl.CaCerts) > 0 {
// 				roots := x509.NewCertPool()
// 				for _, caCert := range m.Ssl.CaCerts {
// 					ok := roots.AppendCertsFromPEM([]byte(caCert))
// 					if !ok {
// 						return fmt.Errorf("failed to parse root certificate")
// 					}
// 				}
// 				tlsConfig.RootCAs = roots
// 			} else {
// 				tlsConfig.InsecureSkipVerify = true
// 			}
// 			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
// 				conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
// 				if err != nil {
// 					fmt.Printf("error in Dial, %s\n", err.Error())
// 				}
// 				return conn, err
// 			}
// 		}
//
// 		sess, err := mgo.DialWithInfo(dialInfo)
// 		if err != nil {
// 			fmt.Printf("error dialing over ssl, %s\n", err.Error())
// 			return fmt.Errorf("Unable to connect to MongoDB, %s\n", err.Error())
// 		}
// 		server.Session = sess
// 	}
//
// }

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
	// "member_status":            "NodeType",
}

// MmapStats - XXX
var MmapStats = map[string]string{
	"mapped_megabytes":     "Mapped",
	"non-mapped_megabytes": "NonMapped",
	"page_faults_per_sec":  "Faults",
}

// WiredTigerStats - XXX
var WiredTigerStats = map[string]string{
	"percent_cache_dirty": "CacheDirtyPercent",
	"percent_cache_used":  "CacheUsedPercent",
}

// Collect - XXX
func Collect(server *Server) error {

	if server.Session == nil {
		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:    []string{server.URL.Host},
			Timeout:  10 * time.Second,
			Database: "amon",
		}

		session, connectionError := mgo.DialWithInfo(mongoDBDialInfo)
		if connectionError != nil {
			return fmt.Errorf("Unable to connect to URL (%s), %s\n", server.URL.Host, connectionError.Error())
		}
		server.Session = session
		server.lastResult = nil

		server.Session.SetMode(mgo.Eventual, true)
		server.Session.SetSocketTimeout(0)
	}

	result := &mongodb.ServerStatus{}
	err := server.Session.DB("amon").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, result)
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

		data := mongodb.NewStatLine(*server.lastResult, *result, server.URL.Host, true, durationInSeconds)
		fmt.Print(data.NodeType)

		statLine := reflect.ValueOf(data).Elem()
		storageEngine := statLine.FieldByName("StorageEngine").Interface()

		for key, value := range DefaultStats {
			val := statLine.FieldByName(value).Interface()
			fmt.Print(key + ":")
			fmt.Print(val)
			fmt.Println("\n-----")
		}

		if storageEngine == "mmapv1" {
			for key, value := range MmapStats {
				val := statLine.FieldByName(value).Interface()
				fmt.Print(key + ":")
				fmt.Print(val)
				fmt.Println("\n-----")
			}
		} else if storageEngine == "wiredTiger" {
			for key, value := range WiredTigerStats {
				val := statLine.FieldByName(value).Interface()
				percentVal := fmt.Sprintf("%.1f", val.(float64)*100)
				floatVal, _ := strconv.ParseFloat(percentVal, 64)
				fmt.Print(key + ":")
				fmt.Print(floatVal)
				fmt.Println("\n-----")
			}
		}
		// for key, value := range data {
		// 	val := statLine.FieldByName(value).Interface()
		// 	fmt.Print(key)
		// }

	}

	// Optional. Switch the session to a monotonic behavior.
	// session.SetMode(mgo.Monotonic, true)
	// result := &ServerStatus{}
	// if err := session.DB(mongoDBDialInfo.Database).Run(bson.D{{"serverStatus", 1}}, &result); err != nil {
	// 	return fmt.Errorf("Unable to collect Mongo Stats from URL (%s), %s\n", localhost.Host, err.Error())
	// }
	//
	// fmt.Print(result)

	// for k, v := range result {
	// 	if k == "connections" {
	// 		var conn ConnectionStats
	// 		err := mapstructure.Decode(v, &conn)
	// 		if err != nil {
	// 			return fmt.Errorf("Unable to collect connection stats %s\n", err.Error())
	// 		}
	//
	// 	}
	// 	if k == "cursors" {
	// 		fmt.Print(v)
	// 	}
	//
	// 	if k == "mem" {
	// 		var mem MemStats
	// 		err := mapstructure.Decode(v, &mem)
	// 		if err != nil {
	// 			return fmt.Errorf("Unable to collect mem stats %s\n", err.Error())
	// 		}
	// 	}
	//
	// 	// fmt.Print(k)
	// 	// fmt.Print(v)
	// 	// fmt.Println("--")
	//
	// }
	// fmt.Print(result)

	return nil
}

//
// func main() {
// 	server := Server{URL: localhost}
// 	f := Collect(&server)
// 	time.Sleep(time.Duration(1) * time.Second)
// 	f = Collect(&server)
// 	fmt.Print(f)
// 	defer server.Session.Close()
// }
