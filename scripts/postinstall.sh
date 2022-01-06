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

# Reload systemd to pickup dripfile.service
systemctl daemon-reload
