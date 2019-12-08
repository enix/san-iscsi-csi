package main

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// Name : Get the name of the resizer plugin
func (r *DothillController) Name() string {
	return r.pluginName
}

// CanSupport : Ensure the plugin will be able to process the resize
func (r *DothillController) CanSupport(pv *v1.PersistentVolume, pvc *v1.PersistentVolumeClaim) bool {
	currentStorage := pv.Spec.Capacity[v1.ResourceName(v1.ResourceStorage)]
	currentBytes, conversionSucceed := currentStorage.AsInt64()
	if !conversionSucceed {
		return false
	}

	requestedStorage := pvc.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	requestedBytes, conversionSucceed := requestedStorage.AsInt64()
	if !conversionSucceed {
		return false
	}

	if requestedBytes < currentBytes {
		klog.Error("volume can only be expanded, not reduced")
		return false
	}

	return true
}

// Resize : Called when a PVC is resized
func (r *DothillController) Resize(pv *v1.PersistentVolume, requestSize resource.Quantity) (resource.Quantity, bool, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.dothillClient.HTTPClient.CloseIdleConnections()
	klog.V(2).Infof("Resize() called with pv: %+v", pv)

	currentStorage := pv.Spec.Capacity[v1.ResourceName(v1.ResourceStorage)]
	klog.V(1).Infof("received resize request for volume %s: %s -> %s", pv.ObjectMeta.Name, currentStorage.ToDec(), requestSize.ToDec())

	storageClassName := pv.ObjectMeta.Annotations[storageClassAnnotationKey]
	klog.V(1).Infof("fetching storage class %s", storageClassName)
	storageClass, err := r.kubeClient.StorageV1().StorageClasses().Get(storageClassName, metav1.GetOptions{})
	if err != nil {
		return currentStorage, false, err
	}
	klog.V(2).Info(storageClass)

	if err = runPreflightChecks(storageClass.Parameters, nil); err != nil {
		return currentStorage, false, err
	}

	err = r.configureClient(storageClass.Parameters)
	if err != nil {
		return currentStorage, false, err
	}

	_, _, err = r.dothillClient.ExpandVolume(pv.ObjectMeta.Name, requestSize.String())
	if err != nil {
		return currentStorage, false, err
	}

	return requestSize, false, nil
}
