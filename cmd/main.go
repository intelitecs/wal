package main

import (
	"log"
	"os"

	"github.com/intelitecs/wal/internal/adapters/app/api"
	"github.com/intelitecs/wal/internal/adapters/core/arithmetic"
	gRPC "github.com/intelitecs/wal/internal/adapters/framework/left/grpc"
	"github.com/intelitecs/wal/internal/adapters/framework/right/db"
	"github.com/intelitecs/wal/internal/ports"
)

func main() {

	var err error
	dbDriver := os.Getenv("DB_DRIVER")
	dsourceName := os.Getenv("DS_NAME")

	// ports

	var dbAdapter ports.DBPort
	var core ports.ArithmeticPort
	var appAdapter ports.APIPort
	var gRPCAdapter ports.GRPCPort

	dbAdapter, err = db.NewAdapter(dbDriver, dsourceName)
	if err != nil {
		log.Fatalf("failed to initiate database connection: %v", err)
	}
	defer dbAdapter.CloseDBConnection()

	core = arithmetic.NewAdapter()
	appAdapter = api.NewApplication(dbAdapter, core)
	gRPCAdapter = gRPC.NewAdapter(appAdapter)
	gRPCAdapter.Run()

}
