package main

import (
	"testing"

	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/node"
)

func Test(t *testing.T) {
	driver := common.Driver{
		Impl: &node.Driver{},
	}

	driver.Test(t)
}
