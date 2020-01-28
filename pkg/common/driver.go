package common

import (
	"flag"
	"fmt"
	"net"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-test/pkg/sanity"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

// PluginName is the public name to be used in storage class etc.
const PluginName = "dothill"

var (
	transport = flag.String("transport", "unix", "transport protocol tu use (unix|tcp)")
	bind      = flag.String("bind", fmt.Sprintf("/var/lib/kubelet/plugins/%s/csi.sock", PluginName), "RPC bind URI (can be a UNIX socket path or any URI)")
)

// Driver contains main resources needed by the driver
// and references the underlying specific driver
type Driver struct {
	impl   csi.IdentityServer
	socket net.Listener
	server *grpc.Server
}

// NewDriver is a convenience function for creating an abstract driver
func NewDriver(impl csi.IdentityServer) *Driver {
	return &Driver{
		impl:   impl,
		server: grpc.NewServer(),
	}
}

// Start does the boilerplate stuff for starting the driver
// it loads its configuration from cli flags
func (driver *Driver) Start() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()
	klog.Infof("starting dothill storage controller v%s", Version)

	socket, err := net.Listen(*transport, *bind)
	if err != nil {
		klog.Fatal(err)
	}
	driver.socket = socket

	csi.RegisterIdentityServer(driver.server, driver.impl)
	if controller, ok := driver.impl.(csi.ControllerServer); ok {
		csi.RegisterControllerServer(driver.server, controller)
	} else if node, ok := driver.impl.(csi.NodeServer); ok {
		csi.RegisterNodeServer(driver.server, node)
	} else {
		klog.Fatalf("cannot start a driver which does not implement anything")
	}

	klog.Infof("driver listening on %s", *bind)
	driver.server.Serve(socket)
}

// Stop shuts down the driver
func (driver *Driver) Stop() {
	driver.server.GracefulStop()
	driver.socket.Close()
}

// Test starts the driver in background
// and runs k8s sanity checks
// It is implemented here in order to avoid duplicating code
func (driver *Driver) Test(t *testing.T) {
	socketPath := "/tmp/csi.sock"
	flag.Set("bind", socketPath)
	go driver.Start()
	defer driver.Stop()

	sanity.Test(t, &sanity.Config{
		TargetPath:  "/tmp/csi-mnt",
		StagingPath: "/tmp/csi-mnt-staging",
		Address:     fmt.Sprintf("unix://%s", socketPath),
	})
}
