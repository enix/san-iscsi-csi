package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

func loadConfiguration() {
	viper.SetConfigName("dothill")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")

	viper.SetDefault("name", "dothill-provisioner")
	viper.SetDefault("fsType", "ext4")
	viper.SetDefault("username", "manage")
	viper.SetDefault("password", "!manage")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func start() error {
	config := &rest.Config{
		Host:            "https://10.14.99.121:6443",
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImFkbWluLXRva2VuLXI4eHp3Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImFkbWluIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiZmQxODgxOTQtNzBkMy0xMWU5LTk4YTctYTI3Yjk0Nzc3ZjBiIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6YWRtaW4ifQ.wxejx7Q05jnx6Ru18sbLIxDxJ8qPljt9mG4Cy46o5cnuohCo6X9FrQ4XTi57E1M9MikhuVZIE-4v2EnioGjlnAZaCYlEKrw5XIYBrtj4MweKZ3ruZCRH46woMK5U__2jTNZ0XfsRJ3TQBgiqwDoi_GjLz1IkgJ04UZWbLZ5bgRxVEZFFAPpNuj84WD1fNxFyfczWnDaIBMyhnlHD0R3F3wV0ZjQ4SR6QGLNXWCjoDN7PIpSgaTcyPw9SUMK1eIqyIhwXavOLABOURFESZfg2DsCajVnaL7qEQhxA9ZF8tErANWTjtlDMP-zAw_kZGpkPB60TaQq4FoNqh5ICQ81Pgg",
	}

	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	return errors.Wrap(err, "unable to get kubernetes client config")
	// }

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
		viper.GetString("name"),
		NewDothillProvisioner(),
		serverVersion.GitVersion,
	)

	log.Println("provision controller listening...")
	pc.Run(wait.NeverStop)
	return nil
}

func main() {
	loadConfiguration()
	err := start()
	if err != nil {
		log.Fatal(err)
	}
}
