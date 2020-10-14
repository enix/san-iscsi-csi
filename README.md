# Dothill (Seagate) AssuredSAN provisioner for Kubernetes

A dynamic persistent volume (PV) provisioner for Dothill AssuredSAN based storage systems.

Developped and tested on HPE MSA2050 series.

Should work with its rebrands like :
- Lenovo S2200/ThinkSystem DS Series
- HPE MSA2000/P2000/MSA2040/MSA2050 Series
- Dell EMC PowerVault ME4 Series
- Quantum StorNex
- ...

### Quickstart

#### Install ISCSI on your node(s)

- Debian/Ubuntu
```sh
apt install open-iscsi
```

- CentOS/RHEL
```sh
yum -y install iscsi-initiator-utils
```

> Please note that multipath is currently not working, so it **should not** be installed!

#### Deploy the provisioner to your cluster

We are using Helm to deploy the provisioner, here is [how to install Helm](https://helm.sh/docs/intro/install/).

Our charts are availables on [charts.enix.io](https://charts.enix.io/).

##### Configure your release

Create a `values.yaml` file. It should contain configuration for your release.

You can find a [basic configuration snippet](./example/values.yaml) in the example directory.

You can have as many storage classes and secrets as you want, they will be created in the same namespace as your release.

##### Install the helm chart

You should first add our charts repository, and then install the chart as follows.

```sh
helm repo add enix https://charts.enix.io/
helm install my-release enix/dothill-provisioner -f ./example/values.yaml
```

#### Run a test pod

To make sure everything went well, there's a example pod you can deploy in the `example/` directory. If the pod reaches the `Running` status, you're good to go!

```sh
kubectl apply -f example/pod.yml
```

### Command-line arguments

You can have a list of all available command line flags using the `-help` switch.

#### Logging

Logging can be modified using the `-v` flag :

- `-v 0` : Standard logs to follow what's going on (default if not specified)
- `-v 9` : Debug logs (quite awful to see)

For advanced logging configuration, see [klog](https://github.com/kubernetes/klog).

#### Development

You can start the drivers over TCP so your remote dev cluster can connect to them.

```
go run ./cmd/<driver> -bind=tcp://0.0.0.0:10000
```

### Testing

You can run sanity checks by using the `sanity` helper script in the `test/` directory:

```
./test/sanity
```