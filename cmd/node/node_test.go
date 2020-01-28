package main

import (
	"testing"

	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/node"
)

func Test(t *testing.T) {
	common.NewDriver(node.NewDriver()).Test(t)
}
