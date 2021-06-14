# SAN iSCSI CSI driver

A dynamic persistent volume (PV) provisioner for iSCSI-compatible SAN based storage systems.

[![Build status](https://gitlab.com/enix.io/san-iscsi-csi/badges/main/pipeline.svg)](https://gitlab.com/enix.io/san-iscsi-csi/-/pipelines)
[![Go Report Card](https://goreportcard.com/badge/github.com/enix/san-iscsi-csi)](https://goreportcard.com/report/github.com/enix/san-iscsi-csi)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

## Introduction

Dealing with persistent storage on kubernetes can be particularly cumbersome, especially when dealing with on-premises installations, or when the cloud-provider persistent storage solutions are not applicable.

Entry-level SAN appliances usually propose a low-cost, still powerful, solution to enable redundant persistent storage, with the flexibility of attaching it to any host on your network.

## This project

`san-iscsi-csi` implements the **Container Storage Interface** in order to facilitate dynamic provisioning of persistent volumes, on an iSCSI-compatible SAN based storage, from a kubernetes cluster.

Considering that this project reached a certain level of maturity, and as of version `3.0.0`, this csi driver is proposed as an open-source project.
As of version `4.0.0` it is distributed under the Apache 2.0 [license](./LICENSE).

### Equipement compatibility

Although this project is developped and tested on HPE MSA 2050 & 2060 equipments, it should work with a lot of other references from various brands.
All dothill AssuredSAN based equipements share a common API which **may or may not be advertised** by the final integrator.
We are therefore looking for tests and feedbacks while using other references.

Dothill systems was acquired by Seagate in 2015 for its AssuredSAN family of hybrid storage.

Seagate continues to maintain the line-up with subsequent series :
- [Seagate AssuredSAN](https://www.seagate.com/fr/fr/support/dothill-san/assuredsan-pro-5000-series/) 3000/4000/5000/6000 series

It is also privately labeled by some of the world's most prominent storage brands :
- [Hewlett Packard Enterprise MSA](https://www.hpe.com/fr/fr/storage/msa-shared-storage.html) 1050, 1060, 1062 and 2050, 2060, 2062 models.
- [Dell EMC PowerVault ME4](https://www.dell.com/en-us/work/shop/productdetailstxn/powervault-me4-series) series.
- Quantum StorNex series.
- Lenovo DS series.
- ...

## Roadmap

This project has reached a `beta` stage, and we hope to promote it to `general availability` with the help of external users and contributors. Feel free to help !

The following features are considered for the near future :
- modular equipment / API support
- additional prometheus metrics

To a lesser extent, the following features are considered for a longer term future :
- Raw blocks support
- Authentication proxy, as appliances lack correct right management

## Features

| Features / Availability   |  roadmap  | alpha | beta  | general availability |
|---------------------------|-----------|-------|-------|----------------------|
| dynamic provisioning      |           |       | 2.3.x |                      |
| resize                    |           | 2.4.x | 3.0.0 |                      |
| snapshot                  |           | 3.1.x |       |                      |
| prometheus metrics        |           | 3.1.x |       |                      |
| modular API support       | mid term  |       |       |                      |
| raw blocks                | long term |       |       |                      |
| iscsi chap authentication | long term |       |       |                      |
| authentication proxy      | long term |       |       |                      |
| overview web ui           | long term |       |       |                      |
| fiber channel             | maybe     |       |       |                      |

## Installation

### Install dependencies on your node(s)

After having shipped `iscsid` and `multipathd` as sidecars on each nodes, we finally decided to extract them and let them run on the host, especially for portability purposes. Indeed, it may happens that the version of `multipathd` running in sidecars don't match the version which would run on the host if it was installed and may produce incompatibilities with the kernel and other tools like `udev`. More about this in the [FAQ](./docs/troubleshooting.md#multipathd-segfault-or-a-volume-got-corrupted).

To run `iscsid` and `multipathd` on your host, first install `open-iscsi` and `multipath-tools` packages, then start the corresponding services. Here is an example on how to do it, it may vary depending on your OS. This was tested against ubuntu 20 and debian buster.

```bash
apt update
apt install open-iscsi multipath-tools -y
service iscsid start
service multipathd start
```

### Multipathd additionnal configuration

For the plugin to work with multipathd, you have to install the following configuration on your hosts. We advise to put it in `/etc/multipath/conf.d/san-iscsi-csi.conf`.

```conf
defaults {
  polling_interval 2
  find_multipaths "yes"
  retain_attached_hw_handler "no"
  disable_changed_wwids "yes"
  user_friendly_names "no"
}
```

After the configuration has been created, restart multipathd to reload it (`service multipathd restart`).

### Deploy the provisioner to your kubernetes cluster

The preferred approach to install this project is to use the provided [Helm Chart](https://artifacthub.io/packages/helm/enix/san-iscsi-csi).

#### Configure your release

Create a `values.yaml` file. It should contain configuration for your release.

Please read the san-iscsi-csi helm-chart [README.md](https://github.com/enix/helm-charts/blob/master/charts/san-iscsi-csi/README.md#values) for more details about this file.

#### Install the helm chart

You should first add our charts repository, and then install the chart as follows.

```sh
helm repo add enix https://charts.enix.io/
helm install my-release enix/san-iscsi-csi -f ./example/values.yaml
```

### Create a storage class

In order to dynamically provision persistants volumes, you first need to create a storage class as well as his associated secret. To do so, please refer to this [example](./example/storage-class.yaml).

### Run a test pod

To make sure everything went well, there's a example pod you can deploy in the `example/` directory. If the pod reaches the `Running` status, you're good to go!

```sh
kubectl apply -f example/pod.yaml
```

## Documentation

You can find more documentation in the [docs](./docs) directory.

## Command-line arguments

You can have a list of all available command line flags using the `-help` switch.

### Logging

Logging can be modified using the `-v` flag :

- `-v 0` : Standard logs to follow what's going on (default if not specified)
- `-v 9` : Debug logs (quite awful to see)

For advanced logging configuration, see [klog](https://github.com/kubernetes/klog).

### Development

You can start the drivers over TCP so your remote dev cluster can connect to them.

```
go run ./cmd/<driver> -bind=tcp://0.0.0.0:10000
```

## Testing

You can run sanity checks by using the `sanity` helper script in the `test/` directory:

```
./test/sanity
```
