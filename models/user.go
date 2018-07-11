package models

import (
	"fmt"
	"log"
	"strings"
)

const usersTableName = "users"

// User contains information about user.
type User struct {
	ID        uint32 `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Gender    string `json:"gender" db:"gender"`
	BirthDate int32  `json:"birth_date" db:"birth_date"`
}

// Users contains list of users.
type Users struct {
	Rows []*User `json:"users"`
}

// GetUser returns user from database specified by id.
func GetUser(id string) (User, error) {
	user := User{}
	if err := GetByID(usersTableName, id, &user); err != nil {
		return user, err
	}
	return user, nil
}

// GetUserVisits returns user's visits from database specified by user's id.
func GetUserVisits(id string, predicates map[string][]string) (Places, error) {
	var (
		names      = []string{"fromDate", "toDate", "country", "toDistance"}
		statements = []string{"visited_at > ", "visited_at < ", "country = ", "distance < "}
		values     = []interface{}{id}
		where      = []string{"\"user\" = $1"}
	)

	for i, name := range names {
		if value, ok := predicates[name]; ok {
			values = append(values, value[0])
			where = append(where, fmt.Sprintf("%s$%d", statements[i], len(values)))
		}
	}

	places := Places{[]*Place{}}
	q := `SELECT mark, visited_at, place 
		  FROM visits 
		  INNER JOIN locations ON visits.location = locations.id
		  WHERE ` + strings.Join(where, " AND ")

	if err := DB.Select(&places.Rows, q, values...); err != nil {
		return places, err
	}
	return places, nil
}

// InsertUser inserts specified user into database.
func InsertUser(user *User) {
	_, err := DB.NamedExec(
		`INSERT INTO users (id, email, first_name, last_name, gender, birth_date) 
		 VALUES (:id, :email, :first_name, :last_name, :gender, :birth_date)`, user)
	if err != nil {
		log.Fatal(err)
	}
}

// PopulateUsers inserts specified list of users into database.
func PopulateUsers(users Users) {
	for _, user := range users.Rows {
		InsertUser(user)
	}
}
