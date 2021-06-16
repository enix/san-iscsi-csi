# Migrate from v3.x to v4.x (on Kubernetes)

While this CSI plugin bumped his major version from 3 to 4, it has changed his name from `dothill-csi` to `san-iscsi-csi`. This results in the need to upgrade StorageClasses using the plugin and replace `provisioner: dothill.csi.enix.io` by `provisioner: san-iscsi.csi.enix.io`. Since StorageClasses are immutables, you need to delete the old StorageClasses and create a new one, but this leads to an issue: your already existing PVCs will not be attached to any StorageClasse anymore. This guide explains one (recommended) way to handle this issue.

## Summary

We will recreate PersistantVolumes with the new version of the provisonner and rename old volumes on the appliance in order to make Kubernetes think the old volumes were created with the new version of the plugin and the new StorageClasses.

## Steps to follow

1. Scale your deployments/statefulsets/etc... using volumes provisonned by the plugin to 0. No pod using those volumes should be running.
1. Ensure that your PersistantVolumes ReclaimPolicy are set to `retain`.
1. Delete PersistantVolumeClaims provisionned by the plugin. Because of the ReclaimPolicy, PersistantVolumes will not be deleted, keep it that way.
1. Uninstall dothill-csi v3 and delete associated StorageClasses
1. Install san-iscsi-csi v4 and recreate StorageClasses with the new `provisioner` value configured to `san-iscsi.csi.enix.io`.
1. Recreate deleted PersistantVolumeClaims.
1. Manually delete on the appliance new volumes provisonned by the plugin.
1. Manually rename on the appliance old volumes with the name of new volumes that you just deleted.
1. Scale back to whatever values your deployments/statefulsets/etc... were when you started the miration.
1. If everything works correctly, delete old PersistantVolumes. You also will have to remove `finalizers` on them to force the deletion since the provisionner handling them (dothill-csi) doesn't exists anymore.
