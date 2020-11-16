package main

import (
	"flag"
	"fmt"

	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/enix/dothill-storage-controller/pkg/node"
	"k8s.io/klog"
)

var kubeletPath = flag.String("kubeletpath", "/var/lib/kubelet", "Kubelet path")
var bind = flag.String("bind", fmt.Sprintf("unix:///var/lib/kubelet/plugins/%s/csi.sock", common.PluginName), "RPC bind URI (can be a UNIX socket path or any URI)")

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()
	klog.Infof("starting dothill storage node plugin %s", common.Version)
	common.NewDriver(node.NewDriver(*kubeletPath)).Start(*bind)
}
