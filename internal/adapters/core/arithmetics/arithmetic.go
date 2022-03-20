package arithmetic

type Adapter struct {
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Addition(x, y int32) (int32, error) {
	return (x + y), nil
}

func (a *Adapter) Subtraction(x, y int32) (int32, error) {
	return (x - y), nil
}

func (a *Adapter) Multiplication(x, y int32) (int32, error) {
	return (x * y), nil
}

func (a *Adapter) Division(x, y int32) (int32, error) {
	return (x / y), nil
}
