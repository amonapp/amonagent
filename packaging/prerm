case "$*" in
  0)
    # We're uninstalling.

    # kill all of the user processes
    ps -u amonagent | grep -v PID | awk '{print $1}' | xargs -i kill {}
    sleep 2
    ps -u amonagent | grep -v PID | awk '{print $1}' | xargs -i kill {}

    # Systemd
    if which systemctl > /dev/null 2>&1 ; then
        systemctl stop amonagent
    # Sysv
    else
        service amonagent stop
    fi

    ;;
  1)
    # We're upgrading. Do nothing.
    ;;
  *)
    ;;
esac

exit 0
