package api

import "github.com/intelitecs/wal/internal/ports"

type Application struct {
	db    ports.ArithmeticDB
	arith ports.Arithmetics
}

func NewApplication(db ports.ArithmeticDB, arith ports.Arithmetics) *Application {
	return &Application{arith: arith, db: db}
}

func (a *Application) GetAddition(x, y int32) (int32, error) {
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

func (a *Application) GetSubtraction(x, y int32) (int32, error) {
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

func (a *Application) GetMultiplication(x, y int32) (int32, error) {
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

func (a *Application) GetDivision(x, y int32) (int32, error) {
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
