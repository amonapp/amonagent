#!/bin/bash

function disable_systemd {
    systemctl disable amonagent
    rm -f /lib/systemd/system/amonagent.service
}

function disable_update_rcd {
    update-rc.d -f amonagent remove
    rm -f /etc/init.d/amonagent
}

function disable_chkconfig {
    chkconfig --del amonagent
    rm -f /etc/init.d/amonagent
}

if [[ -f /etc/redhat-release ]]; then
    # RHEL-variant logic
    if [[ "$1" = "0" ]]; then
  		rm -f /etc/default/amonagent

	  which systemctl &>/dev/null
	  if [[ $? -eq 0 ]]; then
	      disable_systemd
	  else
	      # Assuming sysv
	      disable_chkconfig
	  fi
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    if [[ "$1" != "upgrade" ]]; then
      # Remove/purge
      rm -f /etc/default/amonagent

      which systemctl &>/dev/null
      if [[ $? -eq 0 ]]; then
          disable_systemd
          deb-systemd-invoke stop amonagent.service
      else
          # Assuming sysv
          disable_update_rcd
          invoke-rc.d amonagent stop
      fi
    fi
fi