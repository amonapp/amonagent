#!/bin/bash

LOG_DIR=/var/log/amonagent
SCRIPT_DIR=/opt/amonagent/scripts/
HOME_DIR=/home/amonagent

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/amonagent
    chmod +x /etc/init.d/amonagent

    echo "### You can start amonagent by executing"
    echo ""
    echo " sudo service amonagent start"
    echo ""
    echo "###"
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
}

function install_update_rcd {
    update-rc.d amonagent defaults
}

function install_chkconfig {
    chkconfig --add amonagent
}

id amonagent &>/dev/null
if [[ $? -ne 0 ]]; then
    useradd --system --user-group --key USERGROUPS_ENAB=yes -M amonagent --shell /bin/false -d /etc/opt/amonagent
fi


test -d $LOG_DIR || mkdir -p $LOG_DIR
chown -R -L amonagent:amonagent $LOG_DIR
chmod 755 $LOG_DIR


# Create a dummy home directory, Sensu plugins need this for some reason
test -d $HOME_DIR || mkdir -p $HOME_DIR
chown -R -L amonagent:amonagent $HOME_DIR

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


# Make sure the plugin config directory exists
if [ ! -d /etc/opt/amonagent/plugins-enabled ]; then
    mkdir -p /etc/opt/amonagent/plugins-enabled
fi



# Make sure the pid directory exists
if [ ! -d /var/run/amonagent ]; then
    mkdir -p /var/run/amonagent
fi

chown -R -L amonagent:amonagent  /var/run/amonagent
chmod 775 /var/run/amonagent


# Make sure the binary is executable
chmod +x /usr/bin/amonagent
chmod +x /opt/amonagent/amonagent


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
        systemctl restart amonagent || echo "WARNING: systemd not running."
    else
        # Assuming sysv
        install_init
        install_update_rcd
        invoke-rc.d amonagent restart
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
        # Amazon Linux logic
        install_init
        install_chkconfig
    fi
fi



# Generate machine id, if it does not exists
if [ ! -d /etc/opt/amonagent/machine-id ]; then
    echo "### Checking machine id:"
    echo ""
    amonagent -machineid
    echo ""
    echo "###"
    
fi