package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	DB     *sqlx.DB
	ctx    context.Context
	cancel func()
	DSN    string
	Now    func() time.Time
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() (err error) {
  if db.DSN == "" {
    return fmt.Errorf("dsn is required")
  }

  db.DB, err = sqlx.Open("postgres", db.DSN)
  if err != nil {
    return err
  }

  if err = db.DB.Ping(); err != nil {
    return err
  }

  log.Println("🔗 Connected to database successfully");

  return nil
}

func (db *DB) Close() error {
  db.cancel()
  if db.DB != nil {
    return db.DB.Close()
  }
  return nil
}
