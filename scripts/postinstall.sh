#!/bin/sh
set -e

# Create dripfile group (if it doesn't exist)
if ! getent group dripfile >/dev/null; then
    groupadd --system dripfile
fi

# Create dripfile user (if it doesn't exist)
if ! getent passwd dripfile >/dev/null; then
    useradd                        \
        --system                   \
        --gid dripfile             \
        --shell /usr/sbin/nologin  \
        dripfile
fi

# Update config file permissions (idempotent)
chown root:dripfile /etc/dripfile.conf
chmod 0640 /etc/dripfile.conf

# Reload systemd to pickup dripfile service(s)
systemctl daemon-reload

# Start or restart dripfile components
components="dripfile-migrate dripfile-scheduler dripfile-worker dripfile-web"
for component in $components; do
    if ! systemctl is-enabled $component >/dev/null
    then
        systemctl enable $component >/dev/null
        systemctl start $component >/dev/null
    else
        systemctl restart $component >/dev/null
    fi
done
