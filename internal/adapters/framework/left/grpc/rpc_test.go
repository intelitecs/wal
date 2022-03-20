package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/intelitecs/wal/internal/adapters/app/api"
	"github.com/intelitecs/wal/internal/adapters/core/arithmetic"
	"github.com/intelitecs/wal/internal/adapters/framework/left/grpc/pb"
	"github.com/intelitecs/wal/internal/adapters/framework/right/db"
	"github.com/intelitecs/wal/internal/ports"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	var err error
	dbDriver := os.Getenv("DB_DRIVER")
	dsourceName := os.Getenv("DS_NAME")
	lis = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

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
	gRPCAdapter = NewAdapter(appAdapter)
	pb.RegisterArithmeticServiceServer(grpcServer, gRPCAdapter)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("test server start failed: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func getGRPCConnection(ctx context.Context, t *testing.T) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	return conn
}

func TestGetAddition(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()
	client := pb.NewArithmeticServiceClient(conn)
	params := &pb.OperationParameters{
		A: 1,
		B: 2,
	}
	answer, err := client.GetAddition(ctx, params)
	if err != nil {
		t.Fatalf("expect :%v, got %v", nil, err)
	}
	require.Equal(t, answer.Value, int32(3))

}

func TestGeSubtration(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()
	client := pb.NewArithmeticServiceClient(conn)
	params := &pb.OperationParameters{
		A: 10,
		B: 2,
	}
	answer, err := client.GetSubtraction(ctx, params)
	if err != nil {
		t.Fatalf("expect :%v, got %v", nil, err)
	}
	require.Equal(t, answer.Value, int32(8))

}

func TestGetMultiplication(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()
	client := pb.NewArithmeticServiceClient(conn)
	params := &pb.OperationParameters{
		A: 8,
		B: 2,
	}
	answer, err := client.GetMultiplication(ctx, params)
	if err != nil {
		t.Fatalf("expect :%v, got %v", nil, err)
	}
	require.Equal(t, answer.Value, int32(16))

}

func TestGetDivision(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()
	client := pb.NewArithmeticServiceClient(conn)
	params := &pb.OperationParameters{
		A: 10,
		B: 2,
	}
	answer, err := client.GetDivision(ctx, params)
	if err != nil {
		t.Fatalf("expect :%v, got %v", nil, err)
	}
	require.Equal(t, answer.Value, int32(5))

}
