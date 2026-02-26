package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(conn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Fatal("Db connection error: ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db
}
