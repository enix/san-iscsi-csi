package main

import (
	"fmt"
	"io/ioutil"
	"sync"

	dothill "github.com/enix/dothill-api-go"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

// DothillController : Contains context to interact with k8s and storage
type DothillController struct {
	pluginName    string
	namespace     string
	lock          sync.Mutex
	dothillClient *dothill.Client
	kubeClient    *kubernetes.Clientset
}

// NewDothillController : Creates the controller instance that implements
// the controller.Provisioner and controller.Resizer interface
func NewDothillController(pluginName string, kubeClient *kubernetes.Clientset) *DothillController {
	namespace := "kube-system"
	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		klog.Info(errors.Wrap(err, "failed to get current namespace, using 'kube-system' as a fallback"))
	} else {
		namespace = string(namespaceBytes)
		klog.V(1).Infof("current namespace: %s", namespace)
	}

	return &DothillController{
		pluginName:    pluginName,
		namespace:     namespace,
		dothillClient: dothill.NewClient(),
		kubeClient:    kubeClient,
	}
}

func (p *DothillController) configureClient(parameters map[string]string) error {
	klog.V(1).Infof("fetching dothill credentials from secret %s in namespace %s", parameters[credentialsSecretNameConfigKey], p.namespace)
	credentials, err := p.kubeClient.CoreV1().Secrets(p.namespace).Get(parameters[credentialsSecretNameConfigKey], metav1.GetOptions{})
	if err != nil {
		return err
	}

	username := string(credentials.Data[usernameSecretKey])
	password := string(credentials.Data[passwordSecretKey])
	klog.V(1).Infof("using dothill API at address %s", parameters[apiAddressConfigKey])
	if p.dothillClient.Addr == parameters[apiAddressConfigKey] && p.dothillClient.Username == username {
		klog.V(1).Info("dothill client is already configured for this API, skipping login")
		return nil
	}

	p.dothillClient.Username = username
	p.dothillClient.Password = password
	p.dothillClient.Addr = parameters[apiAddressConfigKey]

	klog.V(1).Infof("login into %s as user %s", p.dothillClient.Addr, p.dothillClient.Username)
	err = p.dothillClient.Login()
	if err != nil {
		return err
	}

	klog.V(1).Info("login was successful")
	return nil
}

func runPreflightChecks(parameters map[string]string, accessModes []v1.PersistentVolumeAccessMode) error {
	checkIfKeyExistsInConfig := func(key string) error {
		klog.V(2).Infof("checking for %s in storage class parameters", key)
		_, ok := parameters[key]
		if !ok {
			return fmt.Errorf("'%s' is missing from configuration", key)
		}
		return nil
	}

	if err := checkIfKeyExistsInConfig(fsTypeConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(poolConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(targetIQNConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(portalsConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(initiatorNameConfigKey); err != nil {
		if err2 := checkIfKeyExistsInConfig(uniqueInitiatorNameByPvcConfigKey); err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
	}
	if err := checkIfKeyExistsInConfig(apiAddressConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(credentialsSecretNameConfigKey); err != nil {
		return err
	}

	for _, mode := range accessModes {
		if mode != v1.ReadWriteOnce {
			return errors.New("dothill storage only supports ReadWriteOnce access mode")
		}
	}

	return nil
}
