0.7.1 - 24.08.2016
==============

* 32bit binary
* Compile with go 1.7


0.7 - 04.08.2016
==============

* Statsd plugin
* Refactoring, speed and stability improvements

0.6.5 - 10.07.2016
==============

* Systemd, RPM distros - properly creates the /var/run/amonagent directory on reboot

0.6.4 - 08.07.2016
==============

* Postinst script permissions changes for properly running Sensu plugins on CentOS

0.6.3 - 08.07.2016
==============

* `amonagent -test` displays each plugin output individualy + show execution time
* New CLI param `amonagent -debug` - starts the agent and display the metrics sent in the terminal

0.6.2 - 03.07.2016
==============

* Automatically start on Ubuntu / Debian
* Generate machine id on first install / run. Fixes duplicate server creation in the Amon interface

0.6 - 27.06.2016
==============

* Stability fixes - cleanup Golang data races in plugin collectors

0.5.5 - 10.05.2016
==============

* Cleanup and update test suite
* Remove hard coded config directories(prepare for future cross platform releases)


0.5 - 30.03.2016
==============

* Remove Metadata URL print statement
* Systemd fixes
* Recompile with Go 1.6


0.4.9 - 03.02.2016
==============

* Telegraf Plugin

0.4.8 - 22.01.2016
==============

* Fixes an issue where the agent will hang and not send data on specific servers

0.4.7 - 21.01.2016
==============

* CloudID fixes


0.4.6 - 20.01.2016
==============

* Postinst script fixes
* Improved test command
* PostgreSQL plugin - slow queries parser fix

0.4.5 - 19.01.2016
==============

* Missing Full Name for database/index size in the MySQL Plugin(Triggers an error in Amon)


0.4.4 - 15.01.2016
==============

* Recompile with Go 1.5.3(Fixes security-related issue in Go) - https://groups.google.com/forum/#!topic/golang-dev/MEATuOi_ei4

0.4.3 - 13.01.2016
==============

* Properly format process memory metric
* Properly format disc metrics for values bigger than Terrabyte

0.4.2 - 09.01.2016
==============

* Generate machine id on first install

0.4.1 - 08.01.2016
==============

* Fix init script on systemd distros

0.4 - 07.01.2016
==============

* Collects all metrics in parallel
* MySQL Plugin
* PostgreSQL Plugin
* MongoDB Plugin
* Redis Plugin
* HAProxy Plugin
* Sensu Plugin
* Nginx Plugin
* Apache Plugin
* Health Checks Plugin
* Can run health checks locally and send the results to Amon
* Custom Plugin - you can write custom plugins in any language with just a couple lines of code.
* New command line options - `list-plugins`, `test-plugin`, `plugin-config`
* Gets Amazon, Google and DigitalOcean instance ids
* Works with self-signed certificates(skips TSL verification)

0.3 - 20.12.2015
==============

* More detailed error messages
* Improve testing command

0.2.5 - 17.12.2015
==============

* Fix permissions issues in the systemd service file

0.2.1 - 15.12.2015
==============

* Machine id parameter
* Fix CPU collector, format data to float

0.2 - 14.12.2015
==============

* Initial release
