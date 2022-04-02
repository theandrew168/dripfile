#!/bin/sh
set -e

systemctl stop dripfile-worker
systemctl stop dripfile-clock
systemctl stop dripfile-web
