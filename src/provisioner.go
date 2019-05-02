package main

import (
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type dothillProvisioner struct{}

// NewDothillProvisioner : Creates the provisionner instance that implements
// the controller.Provisioner interface
func NewDothillProvisioner() controller.Provisioner {
	return &dothillProvisioner{}
}

func (p *dothillProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	fmt.Println(options)
	return nil, errors.New("unimplemented")
}

func (p *dothillProvisioner) Delete(*v1.PersistentVolume) error {
	return errors.New("unimplemented")
}
