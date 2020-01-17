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

#### Install ISCSI and multipath on your node(s)

- Debian/Ubuntu
```sh
apt install open-iscsi multipath-tools
```

- CentOS/RHEL
```sh
yum -y install iscsi-initiator-utils device-mapper-multipath
```

#### Deploy the provisioner to your cluster

As the image is on a private registry for now, make sure you have a secret named `regcred` accessible from the `kube-system` namespace that allows to pull from `docker-registry.enix.io`.

```sh
kubectl apply -f deploy/
```

#### Create a secret containing the Dothill API credentials

Make sure to update the username and password in the file before blindly running this!

```sh
kubectl apply -f example/secret.yml
```

#### Create the storage class

Again, open the file and update all the fields so it match your configuration.

```sh
kubectl apply -f example/storageclass.yml
```

#### Run a test pod

To make sure everything went well, there's a example pod you can deploy in the `example/` directory. If the pod reach the `Running` status, you're good to go!

```sh
kubectl apply -f example/pod.yml
```

### Command-line arguments

You can have a list of all available command line flags using the `-help` switch.

#### Logging

Logging can be modified using the `-v` flag :

- `-v 0` : Almost no logs (default if not specified)
- `-v 1` : Standard logs to follow what's going on
- `-v 2` : Debug logs (quite awful to see)

By default the `rc` image is launched with `-v 1`. For advanced logging configuration, see [klog](https://github.com/kubernetes/klog).

#### Development

You can start the drivers over TCP so your remote dev cluster can connect to them.

```
go run ./cmd/<driver> -transport=tcp -bind=0.0.0.0:10000
```

### Testing

You can run sanity checks by using:

```
go test ./cmd/<driver>
```