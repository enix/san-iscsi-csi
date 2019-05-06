package main

import (
	"log"
	"os"

	"github.com/pborman/getopt/v2"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type args struct {
	Help       bool
	Name       string
	PortalAddr string
	BaseIQN    string
	FSType     string
	Remaining  []string
}

func loadArguments() *args {
	args := args{
		Name:       "dothill-provisioner",
		BaseIQN:    "iqn.2019-05.io.enix",
		PortalAddr: "1.2.3.4:3260",
		FSType:     "ext4",
	}

	getopt.FlagLong(&args.Help, "help", 'h', "display this message")
	getopt.FlagLong(&args.Name, "name", 'n', "provisioner name", args.Name)
	getopt.FlagLong(&args.PortalAddr, "portal", 'p', "portal full address", args.PortalAddr)
	getopt.FlagLong(&args.BaseIQN, "iqn", 'i', "iqn static part", args.BaseIQN)
	getopt.FlagLong(&args.FSType, "fs", 'f', "filesytem to use when formatting the block device", args.FSType)

	opts := getopt.CommandLine
	opts.Parse(os.Args)
	for opts.NArgs() > 0 {
		args.Remaining = append(args.Remaining, opts.Arg(0))
		opts.Parse(opts.Args())
	}

	return &args
}

func start(args *args) error {
	config := &rest.Config{
		Host:            "https://185.145.251.10:6443",
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi11c2VyLXRva2VuLWJ3ejlzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImFkbWluLXVzZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI5ZWVmNTlkMy02Y2RjLTExZTktOGZkNy1mYTE2M2U2ZjhkZmEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06YWRtaW4tdXNlciJ9.HBUDMQYuhAdAYv8RHssc1-D2bVCRbp-2N2uBu-W-rcobF1HPho3MpcFZMFhij4pDyhupiDHKRHv6G2Lo1HCUUxdM7sBUFxliegjB-0JN3JhR9dwkyuW7UqMr_PvHgajDyJYm6muz5PKJlnRyKC3XfDsZrx2WTHs1SPmCxS3CQsCvPUNB871Q1zFn5acCjqbDqQYVK9uP5Hkg3-Qks34z7nglZGuaVB0F_eP2PBNIjGypJSMiNBkd0xtjtlb0dKz50Ed_DRA746CeAubZWHrQn6ySvaeuqwKjVOAVSmzN3MmdeLTgKaMxdDQEtJnDDJslTMcdhbhWdVZGGV5c4fJhEQ",
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
		args.Name,
		NewDothillProvisioner(args),
		serverVersion.GitVersion,
	)

	log.Println("provision controller listening...")
	pc.Run(wait.NeverStop)
	return nil
}

func main() {
	args := loadArguments()
	if args.Help || len(args.Remaining) > 0 {
		getopt.Usage()

		if len(args.Remaining) > 0 {
			os.Exit(1)
		}
		return
	}

	err := start(args)
	if err != nil {
		log.Fatal(err)
	}
}
