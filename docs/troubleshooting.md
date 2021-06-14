# Troubleshooting

## Fixing iSCSI/multipathd state

It might happen that your iSCSI devices/sessions/whatever are in a bad state, for instance the multipath device `/dev/dm-x` might be missing.

In such case, running the following commands should fix the state by removing and recreating devices.

*Please use those commands with **EXTREME CAUTION** and **NEVER IN PRODUCTION** since it can result in data loss.*

```sh
iscsiadm -m node --logout all
iscsiadm -m discovery -t st -p 10.14.84.215
iscsiadm -m node -L all
```

## AttachVolume.Attach failed for volume "xxx" : CSINode xxx does not contain driver san-iscsi.csi.enix.io

Typically, this error happens when you use RancherOS. Since the kubelet path is `/opt/rke/var/lib/kubelet` instead of `/var/lib/kubelet`, the plugin cannot be registered using the default path.

In order to fix this issue, paste the following line in your `value.yaml` and upgrade your helm release.

```yaml
kubeletPath: /opt/rke/var/lib/kubelet
```

## Multipathd segfault or a volume got corrupted

It's a known fact that when `multipathd` segfaults, it can produce wrong mappings of device paths. When such a multipathed device is mounted, it can result in a corruption of the filesystem. Some checks were added to ensure that the different paths are consistent and lead to the same volume in the appliance.

Those segfaults being caused by a mismatch between the candidate version for the package `multipath-tools` on the host and the actually installed one and since `multipathd` now runs on the host instead of as a sidecar of the nodes, this issue shouldn't appears again. If you still get it, please open an issue.

## When expanding a volume, I get the error "missing API credentials"

It's because your storage class miss parameters `csi.storage.k8s.io/controller-expand-secret-name` and `csi.storage.k8s.io/controller-expand-secret-namespace`. The same can happen with volume's creation and publication. The solution is to add those parameters to your storage class. Since a storage class is immutable, you will have to delete it first and then recreate it. The CSI plugin may not take account of this change if an expansion is already in progress. A solution could be to [clone](./volume-snapshots.md#clone-a-volume) the volume you wanted to expand using your new storage class and replace it by its clone.

## multipath is inconsistent: devices WWIDs differ

This issue can be caused by various reasons. In the case the full message looks like the following, it's most probably because you forgot to install the configuration file for multipathd, or forgot to reload it. To find the configuration to apply, refers to the section [multipathd additionnal configuration](https://github.com/enix/san-iscsi-csi/blob/main/README.md#multipathd-additionnal-configuration) in `README.md`.

```
rpc error: code = Unavailable desc = multipath is inconsistent: devices WWIDs differ: mpathb (wwid:3600c0ff00052098834d1c16001000000) != mpathb (wwid:mpathb)
```

Otherwise it may be multipathd just being silly and mapping the wrong devices together, in this case, it may work after a few retries, or may not. If it doesn't, try to manually eject corresponding devices and try again.

## iscsiadm: can not connect to iSCSI daemon (111)!

If you get this error message, try to check that iscsid is running.

```
{output: sh: 0: getcwd() failed: No such file or directory\nFailed to connect to bus: No data available\niscsiadm: can not connect to iSCSI daemon (111)!
```

## device is not mapped to exactly one multipath device: []

If you get this error message, there is a good chance that multipathd is not running. If it is, it may work after a few retry. If it still doesn't work, try to manually eject corresponding devices and try again.

If the device is mapped to more than one multipath device, manually eject corresponding devices and try again.
