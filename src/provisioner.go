package main

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/spf13/viper"

	dothill "enix.io/dothill-api-go"
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

	dothillClient := &dothill.Client{
		Addr:     apiAddress,
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
	}

	err := dothillClient.Login()
	if err != nil {
		log.Fatal(err)
	}

	return &dothillProvisioner{
		baseInitiatorIQN: baseInitiatorIQN,
		targetIQN:        targetIQN,
		portals:          viper.GetStringSlice("portals"),
		fsType:           viper.GetString("fsType"),
		client:           dothillClient,
	}
}

func (p *dothillProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	initiatorName := fmt.Sprintf("%s:%s", p.baseInitiatorIQN, options.SelectedNode.ObjectMeta.Name)
	lun, err := p.chooseLUN(initiatorName)
	if err != nil {
		return nil, err
	}
	volumeName := fmt.Sprintf("%s.%s.%d", p.baseInitiatorIQN, options.SelectedNode.ObjectMeta.Name, lun)

	size := options.PVC.Spec.Resources.Requests["storage"]
	err = p.prepareVolume(volumeName, initiatorName, fmt.Sprintf("%sB", size.String()), lun)
	if err != nil {
		return nil, err
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: volumeName,
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
					Lun:           int32(lun),
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
	fmt.Println("ok")
	return nil
}

func (p *dothillProvisioner) chooseLUN(initiatorName string) (int, error) {
	volumes, _, err := p.client.ShowHostMaps(initiatorName)
	if err != nil {
		return 0, err
	}

	sort.Sort(Volumes(volumes))
	index := 0
	for ; index < len(volumes); index++ {
		if volumes[index].LUN != index+1 {
			return index + 1, nil
		}
	}

	if index+1 < 255 {
		return index + 1, nil
	}

	return -1, errors.New("no more available LUNs")
}

func (p *dothillProvisioner) prepareVolume(volumeName, initiatorName, size string, lun int) error {
	_, _, err := p.client.CreateVolume(volumeName, size, "A")
	if err != nil {
		return err
	}

	_, _, err = p.client.MapVolume(volumeName, initiatorName, lun)
	return err
}
