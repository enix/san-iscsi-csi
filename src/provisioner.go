package main

import (
	"errors"
	"fmt"

	dothill "enix.io/dothill-api-go"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type dothillProvisioner struct {
	baseIQN    string
	portalAddr string
	fsType     string
	client     *dothill.Client
}

// NewDothillProvisioner : Creates the provisionner instance that implements
// the controller.Provisioner interface
func NewDothillProvisioner(args *args) controller.Provisioner {
	return &dothillProvisioner{
		portalAddr: args.PortalAddr,
		baseIQN:    args.BaseIQN,
		fsType:     args.FSType,
		client: dothill.NewClient(&dothill.Options{
			Addr:     args.APIAddr,
			Username: args.Username,
			Password: args.Password,
		}),
	}
}

func (p *dothillProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	lun, err := p.createVolume()
	if err != nil {
		return nil, err
	}

	iqn := fmt.Sprintf("%s:storage-lun%d", p.baseIQN, lun)
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
					TargetPortal: p.portalAddr,
					Portals:      []string{p.portalAddr},
					IQN:          iqn,
					Lun:          int32(lun),
					FSType:       p.fsType,
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

func (p *dothillProvisioner) createVolume() (int32, error) {
	return 0, nil
}
