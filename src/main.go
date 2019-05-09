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
	viper.SetDefault("pool", "A")
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
		BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImFkbWluLXRva2VuLXJkZGtoIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImFkbWluIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiY2Y3ZDEyMjgtNzI1MS0xMWU5LTliOTktYTI3Yjk0Nzc3ZjBiIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6YWRtaW4ifQ.co3p-eQSJt2kNHN-uYgNH4-4FDfA3hmQBHH25TyZ2RmAAe-tlQ9qrEPGPALeiJWiXsfq326HaWqtCMV2tMwCkkdVIwh3b0XdzTdP3DxunMXUUh-Ie6P3nFx-8lK2xwsYIVnc8_U5iHV8YxmVWQl8Ll3ypDGTR8qjJTiq_q1WCWyqaxYpNwPKmpWbqw2zDc4b2OCL72ikAdTlbcgoKJ7XaVZo6VME1iUKglioBWKM_die1SNbgK_CM9eTs5mATSy4Dsyk9OFrXW8ZVjMKM6OwLac695r_b2fN8MV7fezZVG7B4IiYbpC4wQQ_6yTpKNVBoM2qzmoV25lSGwibuvnwSg",
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
