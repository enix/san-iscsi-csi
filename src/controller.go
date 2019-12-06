package main

import (
	"io/ioutil"
	"sync"

	dothill "github.com/enix/dothill-api-go"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

// DothillController : Contains context to interact with k8s and storage
type DothillController struct {
	namespace     string
	lock          sync.Mutex
	dothillClient *dothill.Client
	kubeClient    *kubernetes.Clientset
}

// NewDothillController : Creates the controller instance that implements
// the controller.Provisioner and controller.Resizer interface
func NewDothillController(kubeClient *kubernetes.Clientset) *DothillController {
	namespace := "kube-system"
	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		klog.Info(errors.Wrap(err, "failed to get current namespace, using 'kube-system' as a fallback"))
	} else {
		namespace = string(namespaceBytes)
		klog.V(1).Infof("current namespace: %s", namespace)
	}

	return &DothillController{
		namespace:     namespace,
		dothillClient: dothill.NewClient(),
		kubeClient:    kubeClient,
	}
}
