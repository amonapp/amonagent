package collectors

import (
	"encoding/json"
	"time"

	"github.com/amonapp/amonagent/internal/logging"
	"github.com/amonapp/amonagent/internal/util"
	"github.com/shirou/gopsutil/net"
)

// {'docker0': {'inbound': '0.00', 'outbound': '0.00'}, 'eth0': {'inbound': '0.12', 'outbound': '0.00'}}
var networkLogger = logging.GetLogger("amonagent.net")

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
func (p NetworkStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// NetworkUsageList struct
type NetworkUsageList []NetworkStruct

// NetworkStruct - net interfaces data
type NetworkStruct struct {
	Inbound  float64 `json:"inbound"`
	Outbound float64 `json:"outbound"`
	Name     string  `json:"name"`
}

// NetworkUsage - list
func NetworkUsage() (NetworkUsageList, error) {

	netio, _ := net.IOCounters(true)
	time.Sleep(1000 * time.Millisecond) // Sleep 1 second to get kb/s
	netioSecondRun, _ := net.IOCounters(true)
	ifaces, _ := net.Interfaces()
	var usage NetworkUsageList

	var validInterfaces []string
	for _, iface := range ifaces {
		if !stringInSlice("loopback", iface.Flags) {
			validInterfaces = append(validInterfaces, iface.Name)
		}

	}

	for _, io := range netio {
		if stringInSlice(io.Name, validInterfaces) {

			for _, lastio := range netioSecondRun {
				if lastio.Name == io.Name {
					Inbound := lastio.BytesRecv - io.BytesRecv
					InboundKB, _ := util.ConvertBytesTo(Inbound, "kb", 0)

					Outbound := lastio.BytesSent - io.BytesSent
					OutboundKB, _ := util.ConvertBytesTo(Outbound, "kb", 0)

					n := NetworkStruct{
						Name:     io.Name,
						Inbound:  InboundKB,
						Outbound: OutboundKB,
					}

					usage = append(usage, n)
				}
			}

		}

	}

	return usage, nil

}
