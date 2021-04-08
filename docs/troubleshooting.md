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

## AttachVolume.Attach failed for volume "xxx" : CSINode xxx does not contain driver dothill.csi.enix.io

Typically, this error happens when you use RancherOS. Since the kubelet path is `/opt/rke/var/lib/kubelet` instead of `/var/lib/kubelet`, the plugin cannot be registered using the default path.

In order to fix this issue, paste the following line in your `value.yaml` and upgrade your helm release.

```yaml
kubeletPath: /opt/rke/var/lib/kubelet
```

## Multipathd segfault or a volume got corrupted

It's a known fact that when `multipathd` segfaults, it can produce wrong mappings of device paths. When such a multipathed device is mounted, it can result in a corruption of the filesystem. Some checks were added to ensure that the different paths are consistent and lead to the same volume in the appliance.

If you still get this issue, please check that the candidate for the package `multipath-tools` on your host is on the same version as in the container. You can do so by running `apt-cache policy multipath-tools` on your host as well as in the container `multipathd` from one of the pod `dothill-node-server-xxxxx`.

## When expanding a volume, I get the error "missing API credentials"

It's because your storage class miss parameters `csi.storage.k8s.io/controller-expand-secret-name` and `csi.storage.k8s.io/controller-expand-secret-namespace`. The same can happen with volume's creation and publication. The solution is to add those parameters to your storage class. Since a storage class is immutable, you will have to delete it first and then recreate it. The CSI plugin may not take account of this change if an expansion is already in progress. A solution could be to [clone](./volume-snapshots.md#clone-a-volume) the volume you wanted to expand using your new storage class and replace it by its clone.
