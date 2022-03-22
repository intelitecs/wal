package ports

import (
	"context"

	"github.com/intelitecs/wal/api/v1/arithmetics"
)

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
	GetAddition(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error)
	GetSubtraction(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error)
	GetMultiplication(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error)
	GetDivision(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error)
}
