#!/bin/bash

LOG_DIR=/var/log/amonagent
SCRIPT_DIR=/opt/amonagent/scripts/

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/amonagent
    chmod +x /etc/init.d/amonagent

    echo "### You can start amonagent by executing"
    echo ""
    echo " sudo service start amonagent"
    echo ""
    echo "###"

    invoke-rc.d amonagent restart
}

function install_systemd {
    cp -f $SCRIPT_DIR/amonagent.service /lib/systemd/system/amonagent.service
    systemctl enable amonagent || true
    systemctl daemon-reload || true
    echo "### You can start amonagent by executing"
    echo ""
    echo "sudo systemctl amonagent start"
    echo ""
    echo "###"

    systemctl amonagent restart
}

function install_update_rcd {
    update-rc.d amonagent defaults
}

function install_chkconfig {
    chkconfig --add amonagent
}

id amonagent &>/dev/null
if [[ $? -ne 0 ]]; then
    useradd --system -U -M amonagent -s /bin/false -d /etc/amonagent
fi

test -d $LOG_DIR || mkdir -p $LOG_DIR
chown -R -L amonagent:amonagent $LOG_DIR
chmod 755 $LOG_DIR

# Remove legacy symlink, if it exists
if [[ -L /etc/init.d/amonagent ]]; then
    rm -f /etc/init.d/amonagent
fi

# Add defaults file, if it doesn't exist
if [[ ! -f /etc/default/amonagent ]]; then
    touch /etc/default/amonagent
fi

# Make sure the config directory exists
if [ ! -d /etc/opt/amonagent ]; then
    mkdir -p /etc/opt/amonagent
fi

# Make sure the pid directory exists
if [ ! -d /var/run/amonagent ]; then
    mkdir -p /var/run/amonagent
fi

chown -R -L amonagent:amonagent  /var/run/amonagent
chmod 775 /var/run/amonagent


# Generate machine id if it does not exists
if [ ! -f /etc/opt/amonagent/machine-id ]; then
    dbus-uuidgen > /etc/opt/amonagent/machine-id
fi

# Distribution-specific logic
if [[ -f /etc/redhat-release ]]; then
    # RHEL-variant logic
    which systemctl &>/dev/null
    if [[ $? -eq 0 ]]; then
       install_systemd
    else
       # Assuming sysv
       install_init
       install_chkconfig
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    which systemctl &>/dev/null
    if [[ $? -eq 0 ]]; then
        install_systemd
    else
        # Assuming sysv
        install_init
        install_update_rcd
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
        # Amazon Linux logic
        install_init
        install_chkconfig
    fi
fi
