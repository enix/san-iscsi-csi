package common

import (
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

// PluginName is the public name to be used in storage class etc.
const PluginName = "dothill.csi.enix.io"

// Configuration constants
const (
	FsTypeConfigKey                   = "fsType"
	PoolConfigKey                     = "pool"
	TargetIQNConfigKey                = "iqn"
	PortalsConfigKey                  = "portals"
	InitiatorNameConfigKey            = "initiatorName"
	APIAddressConfigKey               = "apiAddress"
	UniqueInitiatorNameByPvcConfigKey = "uniqueInitiatorNameByPvc"
	UsernameSecretKey                 = "username"
	PasswordSecretKey                 = "password"
	StorageClassAnnotationKey         = "storageClass"

	MaximumLUN          = 255
	VolumeNameMaxLength = 32
)

// Driver contains main resources needed by the driver
// and references the underlying specific driver
type Driver struct {
	impl   csi.IdentityServer
	socket net.Listener
	server *grpc.Server
}

// DriverCtx contains data common to most calls
type DriverCtx struct {
	Credentials map[string]string
	Parameters  map[string]string
	VolumeCaps  []*csi.VolumeCapability
	Req         interface{}
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
func (driver *Driver) Start(bind string) {
	parts := strings.Split(bind, "://")
	if len(parts) < 2 {
		klog.Fatal("please specify a protocol in your bind URI (e.g. \"tcp://\")")
	}

	socket, err := net.Listen(parts[0], parts[1])
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

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		_ = <-sigc
		driver.Stop()
	}()

	klog.Infof("driver listening on %s\n\n", bind)
	driver.server.Serve(socket)
}

// Stop shuts down the driver
func (driver *Driver) Stop() {
	klog.Info("gracefully stopping...")
	driver.server.GracefulStop()
	driver.socket.Close()
}

// BeginRoutine logs every RPC
func (ctx *DriverCtx) BeginRoutine() {
	pc, _, _, _ := runtime.Caller(2)
	caller := runtime.FuncForPC(pc)
	callerNameParts := strings.Split(caller.Name(), ".")
	klog.Infof("=== [ROUTINE START] %s ===", callerNameParts[len(callerNameParts)-1])

	// TODO: find a way to hide credentials
	// klog.V(8).Infof("ARGUMENTS: %+v", ctx.Req)
}
