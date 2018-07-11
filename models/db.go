package models

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"

	// postgres driver
	_ "github.com/lib/pq"
)

// DB contains connection to the database.
var DB *sqlx.DB

// InitDatabase enstablishes connection to the database.
func InitDatabase(driverName string, dataSourceName string) {
	var err error
	DB, err = sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}

// CreateSchema creates tables in database from schema file.
func CreateSchema(file string) {
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	DB.MustExec(string(schema))
}

// GetByID retrieves specified by id model from database.
func GetByID(table string, id string, dest interface{}) error {
	q := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", table)

	if err := DB.Get(dest, q, id); err != nil {
		return err
	}

	return nil
}
