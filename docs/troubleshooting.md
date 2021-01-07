# Troubleshooting

## Fixing iSCSI/multipathd state

It might happen that your iSCSI devices/sessions/whatever are in a bad state, for instance the multipath device `/dev/dm-x` might be missing.

In such case, running the following commands should fix the state by removing and recreating devices.

*Please use those commands with **EXTREM CAUTION** and **NEVER IN PRODUCTION** since it can result in data loss.*

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
