package node

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-lib-iscsi/iscsi"
	"github.com/stretchr/testify/assert"
)

func Test_getConnectorFromFile(t *testing.T) {
	assert := assert.New(t)

	// overwrite _getConnectorFromFile since it use lsblk by default to add additional information
	_getConnectorFromFile = func(filePath string) (*iscsi.Connector, error) {
		f, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		c := iscsi.Connector{}
		err = json.Unmarshal([]byte(f), &c)
		if err != nil {
			return nil, err
		}
		return &c, nil
	}

	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "old format",
			data: []byte(`{
				"volume_name":"test_name",
				"targets":[
					{
						"iqn": "iqn.test",
						"portal": "10.0.0.1",
						"port": "4242"
					},
					{
						"iqn": "iqn.test",
						"portal": "10.0.0.2",
						"port": "1313"
					}
				]
			}`),
		},
		{
			name: "new format",
			data: []byte(`{
				"volume_name": "test_name",
				"target_iqn": "iqn.test",
				"target_portal": ["10.0.0.1:4242", "10.0.0.2:1313"]
			}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile("/tmp/connector.json", tt.data, 0644)
			if err != nil {
				t.Error(err)
			}

			connector, err := getConnectorFromFile("/tmp/connector.json")
			assert.Equal(nil, err)
			assert.Equal(connector.VolumeName, "test_name")
			assert.Equal(connector.TargetIqn, "iqn.test")
			assert.Equal(connector.TargetPortals, []string{"10.0.0.1:4242", "10.0.0.2:1313"})
		})
	}
}
