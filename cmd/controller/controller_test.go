package main

import (
	"testing"

	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/controller"
)

func Test(t *testing.T) {
	common.NewDriver(controller.NewDriver()).Test(t)
}
