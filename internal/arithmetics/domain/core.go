package domain

type Arithmetics interface {
	Addition(int32, int32) (int32, error)
	Subtraction(int32, int32) (int32, error)
	Multiplication(int32, int32) (int32, error)
	Division(int32, int32) (int32, error)
}

type Adapter struct {
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (adapter *Adapter) Addition(a, b int32) (int32, error) {
	return a + b, nil
}

func (adapter *Adapter) Subtraction(a, b int32) (int32, error) {
	return a - b, nil
}

func (adapter *Adapter) Multiplication(a, b int32) (int32, error) {
	return a * b, nil
}

func (adapter *Adapter) Division(a, b int32) (int32, error) {
	return a / b, nil
}
