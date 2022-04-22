#!/bin/sh
set -e

systemctl stop dripfile-scheduler
systemctl stop dripfile-worker
systemctl stop dripfile-web
