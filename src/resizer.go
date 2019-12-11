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
	return true
}

// Resize : Called when a PVC is resized
func (r *DothillController) Resize(pv *v1.PersistentVolume, requestSize resource.Quantity) (resource.Quantity, bool, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.dothillClient.HTTPClient.CloseIdleConnections()
	klog.V(2).Infof("Resize() called with pv: %+v", pv)

	currentSize := pv.Spec.Capacity[v1.ResourceName(v1.ResourceStorage)]
	klog.V(1).Infof("received resize request for volume %s: %s -> %s", pv.ObjectMeta.Name, currentSize.ToDec(), requestSize.ToDec())

	storageClassName := pv.ObjectMeta.Annotations[storageClassAnnotationKey]
	klog.V(1).Infof("fetching storage class %s", storageClassName)
	storageClass, err := r.kubeClient.StorageV1().StorageClasses().Get(storageClassName, metav1.GetOptions{})
	if err != nil {
		return currentSize, false, err
	}
	klog.V(2).Info(storageClass)

	if err = runPreflightChecks(storageClass.Parameters, nil); err != nil {
		return currentSize, false, err
	}

	err = r.configureClient(storageClass.Parameters)
	if err != nil {
		return currentSize, false, err
	}

	additionalSize := requestSize.Copy()
	additionalSize.Sub(currentSize)
	_, _, err = r.dothillClient.ExpandVolume(pv.ObjectMeta.Name, additionalSize.String())
	if err != nil {
		return currentSize, false, err
	}

	return requestSize, true, nil
}
