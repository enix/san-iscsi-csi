# Volume snapshots

## Installation

In order to enable volume snapshotting feature one your cluster, you first need to install the snapshot-controller as well as snapshot CRDs. You can do so by following those [instructions](https://github.com/kubernetes-csi/external-snapshotter#usage).

You will also need to install the snapshot validation webhook, by following those [instructions](https://github.com/kubernetes-csi/external-snapshotter/tree/master/deploy/kubernetes/webhook-example).

## Create a snapshot

To create a snapshot of a volume, you first have to create a `VolumeSnapshotClass`, which is equivalent of a `StorageClass` but for snapshots. Then you can create a `VolumeSnapshot` which use the newly created `VolumeSnapshotClass`. You can follow this [snapshot example](../example/snapshot.yaml). For more informations, please refer to the kubernetes [documentation](https://kubernetes.io/docs/concepts/storage/volume-snapshots/).

## Restore a snapshot

To restore a snapshot, you have to create a new `PersistantVolumeClaim` and specify the desired snapshot as a dataSource. You can find an example [here](https://github.com/kubernetes-csi/external-snapshotter/blob/release-4.0/examples/kubernetes/restore.yaml). You can also refer to the kubernetes [documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#volume-snapshot-and-restore-volume-from-snapshot-support).

## Clone a volume

To clone a volume, you can follow the same procedure than to restore a snapshot, but configure another volume instead of a snapshot. An example can be found [here](https://github.com/kubernetes-csi/csi-driver-host-path/blob/master/examples/csi-clone.yaml) and the kubernetes documentation [here](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#volume-cloning).

---

References:
- https://kubernetes.io/docs/concepts/storage/volume-snapshots
- https://github.com/kubernetes-csi/external-snapshotter
- https://kubernetes-csi.github.io/docs/snapshot-controller
- https://kubernetes-csi.github.io/docs/snapshot-validation-webhook
- https://kubernetes-csi.github.io/docs/snapshot-restore-feature
- https://kubernetes-csi.github.io/docs/volume-cloning
- https://github.com/kubernetes-csi/external-snapshotter/tree/release-4.0/examples/kubernetes
