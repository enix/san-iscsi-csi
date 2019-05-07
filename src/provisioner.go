package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"

	dothill "github.com/enix/dothill-api-go"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type dothillProvisioner struct {
	baseInitiatorIQN string
	targetIQN        string
	portals          []string
	fsType           string
	client           *dothill.Client
}

// NewDothillProvisioner : Creates the provisionner instance that implements
// the controller.Provisioner interface
func NewDothillProvisioner() controller.Provisioner {
	baseInitiatorIQN := viper.GetString("baseInitiatorIQN")
	if len(baseInitiatorIQN) < 1 {
		log.Fatal("'baseInitiatorIQN' missing from configuration")
	}

	portals := viper.GetStringSlice("portals")
	if len(portals) < 1 {
		log.Fatal("'portals' missing from configuration")
	}

	apiAddress := viper.GetString("apiAddress")
	if len(apiAddress) < 1 {
		log.Fatal("'apiAddress' missing from configuration")
	}

	targetIQN := viper.GetString("targetIQN")
	if len(targetIQN) < 1 {
		log.Fatal("'targetIQN' missing from configuration")
	}

	return &dothillProvisioner{
		baseInitiatorIQN: baseInitiatorIQN,
		targetIQN:        targetIQN,
		portals:          viper.GetStringSlice("portals"),
		fsType:           viper.GetString("fsType"),
		client: &dothill.Client{
			Addr:     apiAddress,
			Username: viper.GetString("username"),
			Password: viper.GetString("password"),
		},
	}
}

func (p *dothillProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	lun, err := p.createVolume()
	if err != nil {
		return nil, err
	}

	initiatorName := fmt.Sprintf("%s:%s", p.baseInitiatorIQN, options.SelectedNode.ObjectMeta.Name)
	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s.%d", p.baseInitiatorIQN, options.SelectedNode.ObjectMeta.Name, lun),
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			VolumeMode:                    options.PVC.Spec.VolumeMode,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				ISCSI: &v1.ISCSIPersistentVolumeSource{
					InitiatorName: &initiatorName,
					TargetPortal:  p.portals[0],
					Portals:       p.portals,
					IQN:           p.targetIQN,
					Lun:           lun,
					FSType:        p.fsType,
					ReadOnly:      false,
				},
			},
		},
	}

	log.Println(pv)
	return pv, nil
}

func (p *dothillProvisioner) Delete(*v1.PersistentVolume) error {
	return errors.New("unimplemented")
}

func (p *dothillProvisioner) createVolume() (int32, error) {
	return 1, nil
}
