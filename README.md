# Dothill provisioner for Kubernetes

A dynamic volume provisioner for Dothill storage systems.

### Quickstart

#### Install ISCSI and multipath on your node(s)

```sh
apt update
apt install open-iscsi multipath-tools
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