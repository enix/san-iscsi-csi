package main

import "github.com/enix/dothill-storage-controller/pkg/common"

func main() {
	common.NewDriver(node.NewDriver()).Start()
}
