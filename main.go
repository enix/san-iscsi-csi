package main

import (
	"log"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

func start() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "unable to get kubernetes client config")
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "unable to get k8s client")
	}

	serverVersion, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		return errors.Wrap(err, "failed to get Kubernetes API server version")
	}

	pc := controller.NewProvisionController(
		kubeClient,
		"msa-provisioner",
		NewMSAProvisioner(),
		serverVersion.GitVersion,
	)

	pc.Run(nil)
	return nil
}

func main() {
	err := start()
	if err != nil {
		log.Fatal(err)
	}
}
