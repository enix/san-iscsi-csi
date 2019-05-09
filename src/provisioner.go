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
	size := options.PVC.Spec.Resources.Requests["storage"]
	sizeStr := fmt.Sprintf("%sB", size.String())
	initiatorName := fmt.Sprintf("%s:%s", p.baseInitiatorIQN, options.SelectedNode.ObjectMeta.Name)
	log.Printf("creating %s volume for host %s\n", sizeStr, initiatorName)

	if err := checkAccessMode(options); err != nil {
		return nil, err
	}

	lun, err := p.chooseLUN(initiatorName)
	if err != nil {
		return nil, err
	}

	volumeName := fmt.Sprintf("%s.%s.%d", p.baseInitiatorIQN, options.SelectedNode.ObjectMeta.Name, lun)
	err = p.prepareVolume(volumeName, initiatorName, sizeStr, lun)
	if err != nil {
		return nil, err
	}

	log.Printf("created volume %s (%s) for host %s\n", volumeName, sizeStr, initiatorName)
	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: volumeName,
			Annotations: map[string]string{
				"initiatorName": initiatorName,
			},
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
	}, nil
}

func (p *dothillProvisioner) Delete(volume *v1.PersistentVolume) error {
	name := volume.ObjectMeta.Name
	initiatorName := volume.ObjectMeta.Annotations["initiatorName"]

	log.Printf("deleting volume %s\n", name)
	p.client.UnmapVolume(name, initiatorName)
	p.client.DeleteVolume(name)
	log.Printf("deleted volume %s\n", name)
	return nil
}

func checkAccessMode(options controller.VolumeOptions) error {
	for _, mode := range options.PVC.Spec.AccessModes {
		if mode != v1.ReadWriteOnce {
			return errors.New("dothill storage only supports ReadWriteOnce access mode")
		}
	}

	return nil
}

func (p *dothillProvisioner) chooseLUN(initiatorName string) (int, error) {
	volumes, status, err := p.client.ShowHostMaps(initiatorName)
	if status.ReturnCode == -10074 {
		return 1, nil
	}
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
	_, _, err := p.client.CreateVolume(volumeName, size, viper.GetString("pool"))
	if err != nil {
		return err
	}

	_, _, err = p.client.MapVolume(volumeName, initiatorName, "rw", lun)
	return err
}
