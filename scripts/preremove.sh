#!/bin/sh
set -e

systemctl stop dripfile-web
systemctl stop dripfile-worker
systemctl stop dripfile-scheduler
