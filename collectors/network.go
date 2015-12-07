package collectors

import (
	"encoding/json"

	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/util"
	"github.com/shirou/gopsutil/net"
)

// {'docker0': {'inbound': '0.00', 'outbound': '0.00'}, 'eth0': {'inbound': '0.12', 'outbound': '0.00'}}
var networkLogger = logging.GetLogger("amonagent.net")

// NetworkUsageList struct
type NetworkUsageList []NetworkStruct

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

// NetworkStruct - net interfaces data
type NetworkStruct struct {
	Inbound  float64 `json:"inbound"`
	Outbound float64 `json:"outbound"`
	Name     string  `json:"name"`
}

// NetworkUsage - list
func NetworkUsage() (NetworkUsageList, error) {

	netio, _ := net.NetIOCounters(true)
	ifaces, _ := net.NetInterfaces()
	var usage NetworkUsageList

	var validInterfaces []string
	for _, iface := range ifaces {
		if !stringInSlice("loopback", iface.Flags) {
			validInterfaces = append(validInterfaces, iface.Name)
		}

	}

	for _, io := range netio {
		if stringInSlice(io.Name, validInterfaces) {

			InboundKB, _ := util.ConvertBytesTo(io.BytesRecv, "kb", 0)
			OutboundKB, _ := util.ConvertBytesTo(io.BytesSent, "kb", 0)
			n := NetworkStruct{
				Name:     io.Name,
				Inbound:  InboundKB,
				Outbound: OutboundKB,
			}

			usage = append(usage, n)
		}

	}

	return usage, nil

}
