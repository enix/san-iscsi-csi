package main

import (
	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/controller"
)

func main() {
	common.NewDriver(controller.NewDriver()).Start()
}
