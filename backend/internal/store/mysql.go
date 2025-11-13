package store

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DB struct {
	Conn sqlx.SqlConn
}

func NewDB(dsn string) (*DB, error) {
	conn := sqlx.NewMysql(dsn)
	// connectivity check via a lightweight query
	if _, err := conn.Exec("SELECT 1"); err != nil {
		return nil, fmt.Errorf("mysql connectivity check failed: %v", err)
	}
	return &DB{Conn: conn}, nil
}

func (d *DB) Exec(stmt string) error {
	_, err := d.Conn.Exec(stmt)
	return err
}
