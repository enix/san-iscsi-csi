package main

import (
	"errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type dothillProvisioner struct{}

// NewDothillProvisioner : Creates the provisionner instance that implements
// the controller.Provisioner interface
func NewDothillProvisioner() controller.Provisioner {
	return &dothillProvisioner{}
}

func (p *dothillProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	mode := v1.PersistentVolumeFilesystem
	iqn := "iqn.2019-05.io.enix:storage00"

	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: iqn,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			VolumeMode:                    &mode,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				ISCSI: &v1.ISCSIPersistentVolumeSource{
					TargetPortal: "1.2.3.4:3260",
					Portals:      []string{"1.2.3.4:3260"},
					IQN:          iqn,
					Lun:          0,
					FSType:       "ext4",
					ReadOnly:     false,
					// DiscoveryCHAPAuth: true,
					// SessionCHAPAuth: true,
					// SecretRef: &v1.SecretReference{
					// 	Name: "chap-secret",
					// },
				},
			},
		},
	}, nil
}

func (p *dothillProvisioner) Delete(*v1.PersistentVolume) error {
	return errors.New("unimplemented")
}
