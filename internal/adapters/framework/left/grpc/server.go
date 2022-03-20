package grpc

import (
	"log"
	"net"

	"github.com/intelitecs/wal/internal/adapters/framework/left/grpc/pb"
	"github.com/intelitecs/wal/internal/ports"
	"google.golang.org/grpc"
)

type Adapter struct {
	api ports.APIPort
}

func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

func (grpca *Adapter) Run() {
	var err error
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("listen on port 9000 failure: %v", err)
	}

	arithmeticServiceServer := grpca
	grpcServer := grpc.NewServer()
	pb.RegisterArithmeticServiceServer(grpcServer, arithmeticServiceServer)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC server over port 9000: %v", err)
	}

}
