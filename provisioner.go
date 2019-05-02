package main

import (
	"errors"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type msaProvisioner struct{}

// NewMSAProvisioner : Creates the provisionner instance that implements
// the controller.Provisioner interface
func NewMSAProvisioner() controller.Provisioner {
	return &msaProvisioner{}
}

func (p *msaProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	return nil, errors.New("unimplemented")
}

func (p *msaProvisioner) Delete(*v1.PersistentVolume) error {
	return errors.New("unimplemented")
}
