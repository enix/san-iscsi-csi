package node

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kubernetes-csi/csi-lib-iscsi/iscsi"
)

var _getConnectorFromFile = iscsi.GetConnectorFromFile

func getConnectorFromFile(filePath string) (connector *iscsi.Connector, err error) {
	connector, err = _getConnectorFromFile(filePath)
	if err != nil {
		return
	}

	if connector.TargetIqn != "" && connector.TargetPortals != nil && len(connector.TargetPortals) != 0 { // handle multiple version of this file
		return
	}

	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	type TargetInfo struct {
		Iqn    string `json:"iqn"`
		Portal string `json:"portal"`
		Port   string `json:"port"`
	}
	type OldConnector struct {
		Targets []TargetInfo `json:"targets"`
	}

	oldConnector := OldConnector{}
	err = json.Unmarshal([]byte(f), &oldConnector)
	if err != nil {
		return
	}

	connector.TargetIqn = oldConnector.Targets[0].Iqn
	for _, target := range oldConnector.Targets {
		connector.TargetPortals = append(connector.TargetPortals, target.Portal+":"+target.Port)
	}

	return
}
