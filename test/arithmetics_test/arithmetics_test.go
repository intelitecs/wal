package arithmetics_test

import (
	"log"
	"testing"

	"github.com/intelitecs/wal/internal/arithmetics/domain"
	"github.com/stretchr/testify/require"
)

func TestAdditon(t *testing.T) {
	arith := domain.NewAdapter()

	answer, err := arith.Addition(2, 3)
	if err != nil {
		log.Fatalf("addition failure: %v", err)
	}
	require.Equal(t, answer, int32(5))
}

func TestSubtraction(t *testing.T) {
	arith := domain.NewAdapter()

	answer, err := arith.Subtraction(10, 3)
	if err != nil {
		log.Fatalf("subtraction failure: %v", err)
	}
	require.Equal(t, answer, int32(7))
}

func TestMultication(t *testing.T) {
	arith := domain.NewAdapter()

	answer, err := arith.Multiplication(2, 3)
	if err != nil {
		log.Fatalf("multiplication failure: %v", err)
	}
	require.Equal(t, answer, int32(6))
}

func TestDivision(t *testing.T) {
	arith := domain.NewAdapter()

	answer, err := arith.Division(12, 3)
	if err != nil {
		log.Fatalf("division failure: %v", err)
	}
	require.Equal(t, answer, int32(4))
}
