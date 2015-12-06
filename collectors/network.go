package collectors

import "github.com/martinrusev/amonagent/logging"

var networkLogger = logging.GetLogger("amonagent.net")

// {'docker0': {'inbound': '0.00', 'outbound': '0.00'}, 'eth0': {'inbound': '0.12', 'outbound': '0.00'}}
