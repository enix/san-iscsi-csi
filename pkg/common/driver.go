package common

import (
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/klog"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

// PluginName is the public name to be used in storage class etc.
const PluginName = "dothill.csi.enix.io"

// Configuration constants
const (
	FsTypeConfigKey           = "fsType"
	PoolConfigKey             = "pool"
	TargetIQNConfigKey        = "iqn"
	PortalsConfigKey          = "portals"
	APIAddressConfigKey       = "apiAddress"
	UsernameSecretKey         = "username"
	PasswordSecretKey         = "password"
	StorageClassAnnotationKey = "storageClass"

	MaximumLUN          = 255
	VolumeNameMaxLength = 32
)

// Driver contains main resources needed by the driver
// and references the underlying specific driver
type Driver struct {
	impl   DriverImpl
	socket net.Listener
	server *grpc.Server
}

type DriverImpl interface {
	NewServerInterceptors(logRoutineServerInterceptor grpc.UnaryServerInterceptor) *[]grpc.UnaryServerInterceptor
	ShouldLogRoutine(fullMethod string) bool
}

type WithSecrets interface {
    GetSecrets() map[string]string
}

type WithParameters interface {
    GetParameters() *map[string]string
}

type WithVolumeCaps interface {
    GetVolumeCapabilities() *[]*csi.VolumeCapability
}

// NewDriver is a convenience function for creating an abstract driver
func NewDriver(impl DriverImpl) *Driver {
	return &Driver{
		impl:   impl,
		server: grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				*impl.NewServerInterceptors(newLogRoutineServerInterceptor(impl))...
			)),
		),
	}
}

func newLogRoutineServerInterceptor(impl DriverImpl) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if impl.ShouldLogRoutine(info.FullMethod) {
			klog.Infof("=== [ROUTINE START] %s ===", info.FullMethod)
			defer klog.Infof("=== [ROUTINE END] %s ===", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

// Start does the boilerplate stuff for starting the driver
// it loads its configuration from cli flags
func (driver *Driver) Start(bind string) {
	parts := strings.Split(bind, "://")
	if len(parts) < 2 {
		klog.Fatal("please specify a protocol in your bind URI (e.g. \"tcp://\")")
	}

	if parts[0][:4] == "unix" {
		syscall.Unlink(parts[1])
	}
	socket, err := net.Listen(parts[0], parts[1])
	if err != nil {
		klog.Fatal(err)
	}
	driver.socket = socket

	if identity, ok := driver.impl.(csi.IdentityServer); ok {
		csi.RegisterIdentityServer(driver.server, identity)
	}
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
