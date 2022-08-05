package pkg

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func DbConn() (Db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := "shillings"

	var err error
	Db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(db:3306)/"+dbName)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Connected to SQL DB.\n" + dbUser + ":" + dbPass + "@tcp(db:3306)/" + dbName)
	}

	return Db
}
