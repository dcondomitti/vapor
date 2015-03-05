#!/bin/bash

mkdir -p /tmp/new-drive/openstack/latest
cp user_data /tmp/new-drive/openstack/latest/user_data
mkisofs -R -V config-2 -o configdrive.iso /tmp/new-drive
rm -r /tmp/new-drive

