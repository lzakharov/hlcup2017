package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

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

// Where contains data for SQL Where Clause.
type Where struct {
	Statement string
	Args      map[string]interface{}
}

func prepareWhere(conditions map[string]string, params map[string][]string) Where {
	where := []string{}
	args := map[string]interface{}{}

	for param, condition := range conditions {
		if _, ok := params[param]; ok {
			where = append(where, condition)
			args[param] = params[param][0]
		}
	}

	return Where{strings.Join(where, " AND "), args}
}

func prepareUpdate(table string, columns []string, params map[string]interface{}) string {
	statement := []string{}
	for _, col := range columns {
		if _, ok := params[col]; ok {
			statement = append(statement, col+"=:"+col)
		}
	}
	query := "UPDATE " + table + " SET " + strings.Join(statement, ", ") + " WHERE id=:id"
	return query
}
