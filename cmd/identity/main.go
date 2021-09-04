package main

import (
	"fmt"
	identityProto "identity/pkg/identity/proto"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nite-coder/blackbear/pkg/log"
	"google.golang.org/grpc"
)

func main() {
	defer log.Flush()
	defer func() {
		if r := recover(); r != nil {
			// unknown error
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("unknown error: %v", r)
			}
			log.Err(err).Panic("unknown error")
		}
	}()

	err := initialize()
	if err != nil {
		log.Panicf("main: initialize failed: %v", err)
		return
	}

	// start grpc server
	lis, err := net.Listen("tcp", "cfg.Identity.GRPCBind")
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
