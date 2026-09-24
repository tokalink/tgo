package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLAdapter struct {
	db *sql.DB
}

func NewMySQL(cfg Config) (*MySQLAdapter, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLAdapter{db: db}, nil
}

func (a *MySQLAdapter) Conn() *sql.DB {
	return a.db
}

func (a *MySQLAdapter) Close() error {
	return a.db.Close()
}

func (a *MySQLAdapter) Ping(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

func (a *MySQLAdapter) FromContext(ctx context.Context) (DBEngine, error) {
	return a, nil
}
