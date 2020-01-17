package main

import (
	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/controller"
)

func main() {
	driver := common.Driver{
		Impl: &controller.Driver{},
	}

	driver.Start()
}
