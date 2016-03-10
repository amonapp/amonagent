#! /bin/sh

case "$*" in
  0)
    # We're uninstalling.
    # Systemd
    if which systemctl > /dev/null 2>&1 ; then
        systemctl disable amonagent
    # Sysv
    else

        # update-rc.d sysv service:
        if which update-rc.d > /dev/null 2>&1 ; then
            update-rc.d -f amonagent remove
            update-rc.d amonagent stop
        # CentOS-style sysv:
        else
            chkconfig --del amonagent
        fi
    fi

   	  getent passwd amonagent >/dev/null && userdel  amonagent
	  getent group  amonagent >/dev/null && groupdel amonagent

      # Remove config
      rm -rf /etc/opt/amonagent
      # Remove binary
      rm -rf /opt/amonagent
    ;;
  1)
    # We're upgrading.
    ;;
  *)
    ;;
esac

exit 0
