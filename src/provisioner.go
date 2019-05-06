package main

import (
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type dothillProvisioner struct {
	BaseIQN    string
	PortalAddr string
	FSType     string
}

// NewDothillProvisioner : Creates the provisionner instance that implements
// the controller.Provisioner interface
func NewDothillProvisioner(args *args) controller.Provisioner {
	return &dothillProvisioner{
		PortalAddr: args.PortalAddr,
		BaseIQN:    args.BaseIQN,
		FSType:     args.FSType,
	}
}

func (p *dothillProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	lun := 0
	iqn := fmt.Sprintf("%s:storage-lun%d", p.BaseIQN, lun)
	mode := v1.PersistentVolumeFilesystem

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
					TargetPortal: p.PortalAddr,
					Portals:      []string{p.PortalAddr},
					IQN:          iqn,
					Lun:          int32(lun),
					FSType:       p.FSType,
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
