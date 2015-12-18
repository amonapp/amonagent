#!/usr/bin/env bash

# chkconfig: 2345 95 05
# description: Amon agent - collects system and process information.
# processname: amonagent
# pidfile: /var/run/amonagent/amonagent.pid


### BEGIN INIT INFO
# Provides:          amonagent
# Required-Start:    $remote_fs $syslog
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Starts the Amon agent
# Description:       Amon agent - collects system and process information.
### END INIT INFO

DAEMON='/opt/amonagent/amonagent'
USER="amonagent"
NAME="amonagent"
GROUP="amonagent"
PIDFILE="/var/run/amonagent/amonagent.pid"
PIDPATH=`dirname $PIDFILE`
CONFING="/etc/opt/amonagent/amonagent.conf"

[ -f $AGENTPATH ] || echo "$AGENTPATH not found"


# pid file for the daemon
if [ ! -d "$PIDPATH" ]; then
    mkdir -p $PIDPATH
    chown -R $USER:$GROUP $PIDPATH
fi


DEFAULT=/etc/default/amonagent

if [ -r $DEFAULT ]; then
    source $DEFAULT
fi

# Max open files
OPEN_FILE_LIMIT=65536

if [ -r /lib/lsb/init-functions ]; then
    source /lib/lsb/init-functions
fi

# Logging
if [ -z "$STDOUT" ]; then
    STDOUT=/dev/null
fi

if [ ! -f "$STDOUT" ]; then
    mkdir -p $(dirname $STDOUT)
fi

if [ -z "$STDERR" ]; then
    STDERR=/var/log/amonagent/amonagent.log
fi

if [ ! -f "$STDERR" ]; then
    mkdir -p $(dirname $STDERR)
fi

function pidofproc() {
    if [ $# -ne 3 ]; then
        echo "Expected three arguments, e.g. $0 -p pidfile daemon-name"
    fi

    PID=`pgrep -f $3`
    local PIDFILE=`cat $2`

    if [ "x$PIDFILE" == "x" ]; then
        return 1
    fi

    if [ "x$PID" != "x" -a "$PIDFILE" == "$PID" ]; then
        return 0
    fi

    return 1
}

function killproc() {
    if [ $# -ne 3 ]; then
        echo "Expected three arguments, e.g. $0 -p pidfile signal"
    fi

    PID=`cat $2`

    /bin/kill -s $3 $PID
    while true; do
        pidof `basename $DAEMON` >/dev/null
        if [ $? -ne 0 ]; then
            return 0
        fi

        sleep 1
        n=$(expr $n + 1)
        if [ $n -eq 30 ]; then
            /bin/kill -s SIGKILL $PID
            return 0
        fi
    done
}

function log_failure_msg() {
    echo "$@" "[ FAILED ]"
}

function log_success_msg() {
    echo "$@" "[ OK ]"
}

case $1 in
    start)
        # Check if config file exist
        if [ ! -r $CONFIG ]; then
            log_failure_msg "config file doesn't exists"
            exit 4
        fi

        # Checked the PID file exists and check the actual status of process
        if [ -e $PIDFILE ]; then
            pidofproc -p $PIDFILE $DAEMON > /dev/null 2>&1 && STATUS="0" || STATUS="$?"
            # If the status is SUCCESS then don't need to start again.
            if [ "x$STATUS" = "x0" ]; then
                log_failure_msg "$NAME process is running"
                exit 0 # Exit
            fi
        # if PID file does not exist, check if writable
        else
            su -s /bin/sh -c "touch $PIDFILE" $USER > /dev/null 2>&1
            if [ $? -ne 0 ]; then
                log_failure_msg "$PIDFILE not writable, check permissions"
                exit 5
            fi
        fi

        # Bump the file limits, before launching the daemon. These will carry over to
        # launched processes.
        ulimit -n $OPEN_FILE_LIMIT
        if [ $? -ne 0 ]; then
            log_failure_msg "set open file limit to $OPEN_FILE_LIMIT"
            exit 1
        fi

        log_success_msg "Starting the process" "$NAME"
        if which start-stop-daemon > /dev/null 2>&1; then
            start-stop-daemon --chuid $GROUP:$USER --start --quiet --pidfile $PIDFILE --exec $DAEMON -- -pidfile $PIDFILE >>$STDOUT 2>>$STDERR &
        else
            su -s /bin/sh -c "nohup $DAEMON -pidfile $PIDFILE >>$STDOUT 2>>$STDERR &" $USER
        fi
        log_success_msg "$NAME process was started"
        ;;

    stop)
        # Stop the daemon.
        if [ -e $PIDFILE ]; then
            pidofproc -p $PIDFILE $DAEMON > /dev/null 2>&1 && STATUS="0" || STATUS="$?"
            if [ "$STATUS" = 0 ]; then
                if killproc -p $PIDFILE SIGTERM && /bin/rm -rf $PIDFILE; then
                    log_success_msg "$NAME process was stopped"
                else
                    log_failure_msg "$NAME failed to stop service"
                fi
            fi
        else
            log_failure_msg "$NAME process is not running"
        fi
        ;;

    restart)
        # Restart the daemon.
        $0 stop && sleep 2 && $0 start
        ;;

    status)
        # Check the status of the process.
        if [ -e $PIDFILE ]; then
            if pidofproc -p $PIDFILE $DAEMON > /dev/null; then
                log_success_msg "$NAME Process is running"
                exit 0
            else
                log_failure_msg "$NAME Process is not running"
                exit 1
            fi
        else
            log_failure_msg "$NAME Process is not running"
            exit 3
        fi
        ;;

    version)
        $DAEMON version
        ;;

    *)
        # For invalid arguments, print the usage message.
        echo "Usage: $0 {start|stop|restart|status|version}"
        exit 2
        ;;
esac
