package service

import (
	"github.com/intelitecs/wal/internal/arithmetics/domain"
	"github.com/intelitecs/wal/internal/ports"
)

type ArithmeticApplication struct {
	db    ports.ArithmeticDB
	arith domain.Arithmetics
}

func NewArithmeticApplication(db ports.ArithmeticDB, arith domain.Arithmetics) *ArithmeticApplication {
	return &ArithmeticApplication{arith: arith, db: db}
}

func (a *ArithmeticApplication) GetAddition(x, y int32) (int32, error) {
	answer, err := a.arith.Addition(x, y)
	if err != nil {
		return 0, err
	}
	err = a.db.AddArithmeticToHistory(answer, "addition")
	if err != nil {
		return 0, err
	}
	return answer, nil
}

func (a *ArithmeticApplication) GetSubtraction(x, y int32) (int32, error) {
	answer, err := a.arith.Subtraction(x, y)
	if err != nil {
		return 0, err
	}
	err = a.db.AddArithmeticToHistory(answer, "subtraction")
	if err != nil {
		return 0, err
	}
	return answer, nil
}

func (a *ArithmeticApplication) GetMultiplication(x, y int32) (int32, error) {
	answer, err := a.arith.Multiplication(x, y)
	if err != nil {
		return 0, err
	}
	err = a.db.AddArithmeticToHistory(answer, "multiplication")
	if err != nil {
		return 0, err
	}
	return answer, nil
}

func (a *ArithmeticApplication) GetDivision(x, y int32) (int32, error) {
	answer, err := a.arith.Division(x, y)
	if err != nil {
		return 0, err
	}

	err = a.db.AddArithmeticToHistory(answer, "division")
	if err != nil {
		return 0, err
	}
	return answer, nil
}
