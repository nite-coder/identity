package main

import (
	"fmt"
	"identity/internal/config"
	identityProto "identity/pkg/proto"
	"net"

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
				err = fmt.Errorf("unknown error: %v", err)
			}
			log.StackTrace().Error(err)
		}
	}()

	config := config.New("app.yml")
	err := initialize(config)
	if err != nil {
		log.Panicf("main: initialize failed: %v", err)
		return
	}

	// start grpc server
	lis, err := net.Listen("tcp", config.Identity.GRPCBind)
	if err != nil {
		log.Fatalf("main: bind identity grpc failed: %v", err)
	}
	s := grpc.NewServer()

	identityProto.RegisterIdentityServiceServer(s, _identityServer)
	log.Info("main: grpc service started")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("main: failed to start grpc server: %v", err)
	}
}
