package grpc

import (
	"context"

	arithmetics "github.com/intelitecs/wal/api/v1/arithmetics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAddition gets the result of adding parameters a and b
func (grpca *Adapter) GetAddition(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error) {
	var err error
	ans := &arithmetics.Answer{}

	if req.GetA() == 0 || req.GetB() == 0 {
		return ans, status.Error(codes.InvalidArgument, "missing required")
	}

	answer, err := grpca.api.GetAddition(req.A, req.B)
	if err != nil {
		return ans, status.Error(codes.Internal, "unexpected error")
	}

	ans = &arithmetics.Answer{
		Value: answer,
	}

	return ans, nil
}

// GetSubtraction gets the result of subtracting parameters a and b
func (grpca *Adapter) GetSubtraction(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error) {
	var err error
	ans := &arithmetics.Answer{}

	if req.GetA() == 0 || req.GetB() == 0 {
		return ans, status.Error(codes.InvalidArgument, "missing required")
	}

	answer, err := grpca.api.GetSubtraction(req.A, req.B)
	if err != nil {
		return ans, status.Error(codes.Internal, "unexpected error")
	}

	ans = &arithmetics.Answer{
		Value: answer,
	}

	return ans, nil
}

// GetMultiplication gets the result of multiplying parameters a and b
func (grpca *Adapter) GetMultiplication(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error) {
	var err error
	ans := &arithmetics.Answer{}

	if req.GetA() == 0 || req.GetB() == 0 {
		return ans, status.Error(codes.InvalidArgument, "missing required")
	}

	answer, err := grpca.api.GetMultiplication(req.A, req.B)
	if err != nil {
		return ans, status.Error(codes.Internal, "unexpected error")
	}

	ans = &arithmetics.Answer{
		Value: answer,
	}

	return ans, nil
}

// GetDivision gets the result of dividing parameters a and b
func (grpca *Adapter) GetDivision(ctx context.Context, req *arithmetics.OperationParameters) (*arithmetics.Answer, error) {
	var err error
	ans := &arithmetics.Answer{}

	if req.GetA() == 0 || req.GetB() == 0 {
		return ans, status.Error(codes.InvalidArgument, "missing required")
	}

	answer, err := grpca.api.GetDivision(req.A, req.B)
	if err != nil {
		return ans, status.Error(codes.Internal, "unexpected error")
	}

	ans = &arithmetics.Answer{
		Value: answer,
	}

	return ans, nil
}
