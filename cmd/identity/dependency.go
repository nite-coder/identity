package main

import (
	identityGRPC "identity/pkg/identity/delivery/grpc"
)

var (
	_identityServer *identityGRPC.IdentityServer
)

func initialize() error {
	// ctx := context.Background()

	// err := startup.InitConfig()
	// if err != nil {
	// 	return err
	// }

	// err = startup.InitLogger()
	// if err != nil {
	// 	return err
	// }

	// db, err := setupDatabase(cfg)
	// if err != nil {
	// 	return err
	// }

	// accountRepo := identityDatabase.NewAccountRepo(db)
	// accountSvc := usecase.NewAccountUsecase(accountRepo)

	// _identityServer = identityGRPC.NewIdentityServer(accountSvc)

	return nil

}
