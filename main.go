package amonagent

// package main
//
// import (
// 	"log"
// 	"time"
//
// 	"github.com/martinrusev/amonagent/collectors"
// 	"github.com/martinrusev/amonagent/logging"
// 	"github.com/martinrusev/amonagent/remote"
// )
//
// // AmonAgentLogger for the main file
// var AmonAgentLogger = logging.GetLogger("amonagent")
//
// func GatherAndSend() error {
// 	allMetrics := collectors.CollectSystem()
// 	remote.SendData(allMetrics)
//
// 	return nil
// }
//
// // Just for testing
// func main() {
//
// 	ticker := time.NewTicker(10 * time.Second)
//
// 	for {
// 		if err := GatherAndSend; err != nil {
// 			log.Printf(err.Error())
// 		}
//
// 		select {
// 		case <-shutdown:
// 			return nil
// 		case <-ticker.C:
// 			continue
// 		}
// 	}
//
// }
