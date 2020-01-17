package main

import (
	"testing"

	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/controller"
)

func Test(t *testing.T) {
	driver := common.Driver{
		Impl: &controller.Driver{},
	}

	driver.Test(t)
}
