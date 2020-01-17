package main

import (
	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/node"
)

func main() {
	driver := common.Driver{
		Impl: &node.Driver{},
	}

	driver.Start()
}
