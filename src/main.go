package main

import (
	"log"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

const (
	pluginName = "dothill"

	fsTypeConfigKey                = "fsType"
	poolConfigKey                  = "pool"
	targetIQNConfigKey             = "iqn"
	portalsConfigKey               = "portals"
	initiatorNameConfigKey         = "initiatorName"
	apiAddressConfigKey            = "apiAddress"
	credentialsSecretNameConfigKey = "credentialsSecretName"
	usernameSecretKey              = "username"
	passwordSecretKey              = "password"
	storageClassAnnotationKey      = "storageClass"

	maximumLUN                    = 255
	hostDoesNotExistsErrorCode    = -10386
	hostMapDoesNotExistsErrorCode = -10074
)

func start(config *rest.Config) error {
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
		pluginName,
		NewDothillProvisioner(kubeClient),
		serverVersion.GitVersion,
	)

	log.Println("provision controller listening...")
	pc.Run(wait.NeverStop)
	return nil
}

func loadConfiguration(kubeconfigPath string) (*rest.Config, error) {
	var config *rest.Config
	var err error

	if len(kubeconfigPath) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, errors.Wrap(err, "unable to get kubernetes client config")
	}

	return config, nil
}

func main() {
	kubeconfigPath := flag.String("kubeconfig", "", "path to the kubeconfig file to use instead of in-cluster configuration")
	err := start()
	flag.Parse()

	config, err := loadConfiguration(*kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}

	err = start(config)
	if err != nil {
		log.Fatal(err)
	}
}
