## Vapor

Creates simple cloud-init templates (in terrible, terrible go) for running CoreOS on vSphere / VMWare Fusion. This is needed because variables like $public_ipv4 are only available from certain providers. configdrive.iso must be built with a label of `config-2` and attached to the VM. CoreOS will execute the user\_data script contained in the ISO and pull down a dynamic cloudinit file[0]. Most of this was built thanks to a hackzilla post[1].

### Build Config Drive

This references /openstack due to it originally being built for Openstack Nova[2].

```
mkdir -p /tmp/new-drive/openstack/latest
cp user_data /tmp/new-drive/openstack/latest/user_data
mkisofs -R -V config-2 -o configdrive.iso /tmp/new-drive
rm -r /tmp/new-drive
```

### user_data script

```
#!/bin/bash
 
MAC=`ifconfig eno16780032 | grep -o -E '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}'`
 
URL="https://vapor.interval.io/config/host/${MAC}"
 
coreos-cloudinit --from-url="${URL}"
```

[0] https://coreos.com/docs/cluster-management/setup/cloudinit-config-drive/

[1] http://blog.hackzilla.org/posts/2014/09/02/coreos-dynamic-cloud-init

[2] http://docs.openstack.org/user-guide/content/enable_config_drive.html#config_drive_contents
