curl -s -w 'ping.amoncx.lookup_time:%{time_namelookup}|gauge
ping.amoncx.connect_time:%{time_connect}|gauge
ping.amoncx.total:%{time_total}|gauge\n' -o /dev/null https://www.amon.cx
