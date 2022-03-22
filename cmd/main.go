package main

import (
	"log"
	"os"

	"github.com/intelitecs/wal/internal/arithmetics/domain"
	service "github.com/intelitecs/wal/internal/arithmetics/service"
	"github.com/intelitecs/wal/internal/ports"
	"github.com/intelitecs/wal/internal/server/db"
	gRPC "github.com/intelitecs/wal/internal/server/grpc"
)

func main() {
	var err error
	dbDriver := os.Getenv("DB_DRIVER")
	dsourceName := os.Getenv("DS_NAME")

	// ports

	var dbAdapter ports.ArithmeticDB
	var core domain.Arithmetics
	var appAdapter ports.APIPort
	var gRPCAdapter ports.GRPCPort

	dbAdapter, err = db.NewAdapter(dbDriver, dsourceName)
	if err != nil {
		log.Fatalf("failed to initiate database connection: %v", err)
	}
	defer dbAdapter.CloseDBConnection()

	// core

	core = domain.NewAdapter()

	appAdapter = service.NewArithmeticApplication(dbAdapter, core)
	gRPCAdapter = gRPC.NewAdapter(appAdapter)
	gRPCAdapter.Run()

}
