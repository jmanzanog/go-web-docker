package src

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	host     = "ruby.db.elephantsql.com"
	port     = 5432
	user     = "xqirljbw"
	password = "Mo17RtzuQSGGge9VuUDN1EHnYquIcbK9"
	dbname   = "xqirljbw"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

type Database struct {
	Client *sql.DB
}

var single *Database

//singleton de cliente base de datos
func (d Database) GetClient() *Database {
	if single == nil {
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		db.SetMaxOpenConns(30)
		db.SetMaxIdleConns(30)
		db.SetConnMaxLifetime(3600 * time.Second)
		single = &Database{Client: db}
	}
	return single
}
