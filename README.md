# Dothill-csi dynamic provisioner for Kubernetes

[![Build status](https://gitlab.com/enix.io/dothill-csi/badges/main/pipeline.svg)](https://gitlab.com/enix.io/dothill-csi/-/pipelines)
[![Go Report Card](https://goreportcard.com/badge/github.com/enix/dothill-csi)](https://goreportcard.com/report/github.com/enix/dothill-csi)

A dynamic persistent volume (PV) provisioner for Dothill AssuredSAN based storage systems.

## Introduction

Dothill systems was acquired by Seagate in 2015 for its AssuredSAN family of hybrid storage.

Seagate continues to maintain the line-up with subsequent series :
- Seagate AssuredSAN 3000/4000/5000/6000 series

It is also privately labeled by some of the world's most prominent storage brands :
- Hewlett Packard Enterprise MSA 10xx and 20xx series
- Dell EMC PowerVault ME4 Series
- Quantum StorNex series
- Lenovo DS series
- ...

This project is mainly developped and tested on HPE MSA 2050 & 2060 models.
We value your potential tests on other models.

## Roadmap

Considering that this project reached a certain level of maturity, and as of version `3.0.0`, we propose this csi driver as an opensource project under the MIT license.

This project has reached a `beta` stage, and we hope to promote it to `general availability` with the help of external users and contributors. Feel free to help !

The following features are considered for the near future :
- PV snapshotting (supported by AssuredSAN appliances)
- additional prometheus metrics

To a lesser extent, the following features are considered for a longer term future :
- FiberChannel (supported by AssuredSAN appliances)
- Authentication proxy, as any appliance account has full access to all volumes

## Features

| Features / Availability | alpha | beta  | general availability | roadmap   |
|-------------------------|-------|-------|----------------------|-----------|
| Dynamic provisioning    |       | 2.3.x |                      |           |
| Resize                  | 2.4.x |       |                      |           |
| Snapshot                |       |       |                      | 3.1.x     |
| prometheus metrics      |       |       |                      | 3.x       |
| FiberChannel            |       |       |                      | long term |

## Usage

### Uninstall ISCSI tools on your node(s)

`iscsid` and `multipathd` are now shipped as sidecars on each nodes, it is therefore suggested to uninstall any `open-iscsi` package.

### Deploy the provisioner to your cluster

We are using Helm to deploy the provisioner, here is [how to install Helm](https://helm.sh/docs/intro/install/).

Our charts are availables on [charts.enix.io](https://charts.enix.io/).

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

In order to dynamically provision persistants volumes, you first need to create a storage class as well as his associated secret. To do so, please refer to this [example](./example/storage-class).

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