package ports

import (
	"context"

	"github.com/intelitecs/wal/internal/adapters/framework/left/grpc/pb"
)

type Arithmetics interface {
	Addition(int32, int32) (int32, error)
	Subtraction(int32, int32) (int32, error)
	Multiplication(int32, int32) (int32, error)
	Division(int32, int32) (int32, error)
}

type APIPort interface {
	GetAddition(a, b int32) (int32, error)
	GetSubtraction(a, b int32) (int32, error)
	GetMultiplication(a, b int32) (int32, error)
	GetDivision(a, b int32) (int32, error)
}

type ArithmeticDB interface {
	CloseDBConnection()
	AddArithmeticToHistory(answer int32, operation string) error
}

type GRPCPort interface {
	Run()
	GetAddition(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error)
	GetSubtraction(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error)
	GetMultiplication(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error)
	GetDivision(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error)
}
