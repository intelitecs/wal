package db

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

type Adapter struct {
	db *sql.DB
}

func NewAdapter(driverName, dataSourceName string) (*Adapter, error) {
	// connect
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("db connection failure: %v", err)
	}
	// test db connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("db ping failure: %v", err)
	}
	return &Adapter{db: db}, nil

}

func (a *Adapter) CloseDBConnection() {
	err := a.db.Close()
	if err != nil {
		log.Fatalf("close connection failure: %v", err)
	}

}

func (a *Adapter) AddArithmeticToHistory(answer int32, operation string) error {
	queryString, args, err := sq.Insert("arith_history").Columns("date", "answer", "operation").Values(time.Now(), answer, operation).ToSql()
	if err != nil {
		return err
	}
	_, err = a.db.Exec(queryString, args...)
	if err != nil {
		return err
	}
	return nil
}
