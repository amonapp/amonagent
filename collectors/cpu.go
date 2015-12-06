package collectors

import "github.com/martinrusev/amonagent/logging"

var cpuLogger = logging.GetLogger("amonagent.cpu")

//{'iowait': '0.00', 'system': '0.50', 'idle': '98.50', 'user': '1.00', 'steal': '0.00', 'nice': '0.00'}
