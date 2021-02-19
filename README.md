# Dothill-csi dynamic provisioner for Kubernetes
A dynamic persistent volume (PV) provisioner for Dothill AssuredSAN based storage systems.

[![Build status](https://gitlab.com/enix.io/dothill-csi/badges/main/pipeline.svg)](https://gitlab.com/enix.io/dothill-csi/-/pipelines)
[![Go Report Card](https://goreportcard.com/badge/github.com/enix/dothill-csi)](https://goreportcard.com/report/github.com/enix/dothill-csi)


## Introduction
Dealing with persistent storage on kubernetes can be particularly cumbersome, especially when dealing with on-premises installations, or when the cloud-provider persistent storage solutions are not applicable.

Entry-level SAN appliances usually propose a low-cost, still powerful, solution to enable redundant persistent storage, with the flexibility of attaching it to any host on your network.

Dothill systems was acquired by Seagate in 2015 for its AssuredSAN family of hybrid storage.

Seagate continues to maintain the line-up with subsequent series :
- [Seagate AssuredSAN](https://www.seagate.com/fr/fr/support/dothill-san/assuredsan-pro-5000-series/) 3000/4000/5000/6000 series

It is also privately labeled by some of the world's most prominent storage brands :
- [Hewlett Packard Enterprise MSA](https://www.hpe.com/fr/fr/storage/msa-shared-storage.html) 1050, 1060, 1062 and 2050, 2060, 2062 models.
- [Dell EMC PowerVault ME4](https://www.dell.com/en-us/work/shop/productdetailstxn/powervault-me4-series) series.
- Quantum StorNex series.
- Lenovo DS series.
- ...

## This project
`Dothill-CSI` implements the **Container Storage Interface** in order to facilitate dynamic provisioning of persistent volumes on kubernetes cluster.

All dothill AssuredSAN based equipements share a common API which **may or may not be advertised** by the final integrator.
Although this project is developped and tested on HPE MSA 2050 & 2060 equipments, it should work with a lot of other references from various brands.
We are therefore looking for tests and feedbacks while using other references.

Considering that this project reached a certain level of maturity, and as of version `3.0.0`, this csi driver is proposed as an open-source project under the MIT [license](./LICENSE).

## Roadmap

This project has reached a `beta` stage, and we hope to promote it to `general availability` with the help of external users and contributors. Feel free to help !

The following features are considered for the near future :
- PV snapshotting (supported by AssuredSAN appliances)
- additional prometheus metrics

To a lesser extent, the following features are considered for a longer term future :
- Raw blocks support
- FiberChannel (supported by AssuredSAN appliances)
- Authentication proxy, as appliances lack correct right management

## Features

| Features / Availability |  roadmap  | alpha | beta  | general availability |
|-------------------------|-----------|-------|-------|----------------------|
| dynamic provisioning    |           |       | 2.3.x |                      |
| resize                  |           | 2.4.x |       |                      |
| snapshot                | 3.1.x     |       |       |                      |
| prometheus metrics      | 3.2.x     |       |       |                      |
| raw blocks              | long term |       |       |                      |
| fiber channel           | long term |       |       |                      |
| authentication proxy    | long term |       |       |                      |

## Installation

### Uninstall ISCSI tools on your node(s)

`iscsid` and `multipathd` are now shipped as sidecars on each nodes, it is therefore strongly suggested to uninstall any `open-iscsi` package.

### Deploy the provisioner to your kubernetes cluster

The preferred approach to install this project is to use the provided [Helm Chart](https://artifacthub.io/packages/helm/enix/dothill-csi).

#### Configure your release

Create a `values.yaml` file. It should contain configuration for your release.

Please read the dothill-csi helm-chart [README.md](https://github.com/enix/helm-charts/blob/master/charts/dothill-csi/README.md#values) for more details about this file.

#### Install the helm chart

You should first add our charts repository, and then install the chart as follows.

```sh
helm repo add enix https://charts.enix.io/
helm install my-release enix/dothill-csi -f ./example/values.yaml
```

### Create a storage class

In order to dynamically provision persistants volumes, you first need to create a storage class as well as his associated secret. To do so, please refer to this [example](./example/storage-class.yaml).

### Run a test pod

To make sure everything went well, there's a example pod you can deploy in the `example/` directory. If the pod reaches the `Running` status, you're good to go!

```sh
kubectl apply -f example/pod.yaml
```

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