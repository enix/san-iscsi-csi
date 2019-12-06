package main

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
)

// Name : Get the name of the resizer plugin
func (r *DothillController) Name() string {
	klog.Info("name asked")
	return r.pluginName
}

// CanSupport : Ensure the plugin will be able to process the resize
func (r *DothillController) CanSupport(pv *v1.PersistentVolume, pvc *v1.PersistentVolumeClaim) bool {
	klog.Info(pv)
	klog.Info(pvc)
	klog.Info("can support asked")
	return false
}

// Resize : Called when a PVC is resized
func (r *DothillController) Resize(
	pv *v1.PersistentVolume, requestSize resource.Quantity,
) (
	newSize resource.Quantity, fsResizeRequired bool, err error,
) {
	return *resource.NewQuantity(0, resource.BinarySI), false, errors.New("unsupported yet")
}
