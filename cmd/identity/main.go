package main

import (
	"fmt"
	"identity/internal/pkg/config"
	identityProto "identity/pkg/identity/proto"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonsoft/log"
	"google.golang.org/grpc"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// unknown error
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("unknown error: %v", r)
			}
			trace := make([]byte, 4096)
			runtime.Stack(trace, true)
			log.WithField("stack_trace", string(trace)).WithError(err).Panic("unknown error")

		}
	}()

	config.EnvPrefix = "IDENTITY"
	cfg := config.New("app.yml")
	err := initialize(cfg)
	if err != nil {
		log.Panicf("main: initialize failed: %v", err)
		return
	}

	// start grpc server
	lis, err := net.Listen("tcp", cfg.Identity.GRPCBind)
	if err != nil {
		log.Fatalf("main: bind identity grpc failed: %v", err)
	}
	grpcServer := grpc.NewServer()

	identityProto.RegisterIdentityServiceServer(grpcServer, _identityServer)
	log.Info("main: grpc service started")

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("main: failed to start grpc server: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("main: shutting down server...")

	grpcServer.GracefulStop()
}
