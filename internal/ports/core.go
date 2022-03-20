package ports

import (
	"io"

	logApi "github.com/intelitecs/wal/api/v1/log"
)

type ArithmeticPort interface {
	Addition(a int32, b int32) (int32, error)
	Subtraction(a int32, b int32) (int32, error)
	Division(a int32, b int32) (int32, error)
	Multiplication(a int32, b int32) (int32, error)
}

type LogPort interface {
	Append(record *logApi.Record) (int64, error)
	Read(int64) (logApi.Record, error)
	Setup() error
	Close() error
	Remove() error
	LowestOffset() (int64, error)
	HighestOffset() (int64, error)
	Truncate(int64) error
	Reader() io.Reader
}
