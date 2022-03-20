package arithmetic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddtion(t *testing.T) {
	arith := NewAdapter()
	answer, err := arith.Addition(1, 1)
	if err != nil {
		t.Fatalf("expected: %v, got: %v", nil, err)
	}
	require.Equal(t, answer, int32(2))
}

func TestSubtracion(t *testing.T) {
	arith := NewAdapter()
	answer, err := arith.Subtraction(15, 12)
	if err != nil {
		t.Fatalf("expected: %v, got: %v", nil, err)
	}
	require.Equal(t, answer, int32(3))
}

func TestMultiplication(t *testing.T) {
	arith := NewAdapter()
	answer, err := arith.Multiplication(8, 5)
	if err != nil {
		t.Fatalf("expected: %v, got: %v", nil, err)
	}
	require.Equal(t, answer, int32(40))
}

func TestDivision(t *testing.T) {
	arith := NewAdapter()
	answer, err := arith.Division(10, 2)
	if err != nil {
		t.Fatalf("expected: %v, got: %v", nil, err)
	}
	require.Equal(t, answer, int32(5))
}
