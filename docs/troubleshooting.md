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
