package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/neutrixs/lexiquiz-server/pkg/env"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	
	user, _ := env.Get("MYSQL_USER")
	pass, _ := env.Get("MYSQL_PASS")
	pass = url.QueryEscape(pass)
	host, _ := env.Get("MYSQL_HOST")
	port, err := env.Get("MYSQL_PORT")
	if err != nil {
		port = "3306"
	}
	name, _ := env.Get("MYSQL_DB")

	connectionFormat := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)

	db, err = sql.Open("mysql", connectionFormat)
	if err != nil {
		log.Println(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		db, _ = sql.Open("mysql", connectionFormat)
	}
	dbCheck()
	go func() {
		for {
			time.Sleep(10 * time.Second)
			err = db.Ping()
			if err != nil {
				log.Println(err)
				db, _ = sql.Open("mysql", connectionFormat)
			}
		}
	}()
}