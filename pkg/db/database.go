package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (Db *sql.DB) {

	dbDriver := "mysql"
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbPort := "3306"
	dbName := "shillings"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(db:"+dbPort+")/"+dbName)

	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Connected to SQL DB.\n" + dbUser + ":" + dbPass + "@tcp(db:3306)/" + dbName)
	}

	db.SetConnMaxLifetime(3 * time.Hour)
	db.SetMaxOpenConns(3000)
	db.SetMaxIdleConns(3000)

	return db
}
