package models

import (
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
	err := GetByID(usersTableName, id, &user)
	return user, err
}

// GetUserVisits returns user's visits from database specified by user's id.
func GetUserVisits(id string, params map[string][]string) (Places, error) {
	params["id"] = []string{id}

	conditions := map[string]string{
		"id":       `"user"=:id`,
		"fromDate": "visited_at>:fromDate",
		"toDate":   "visited_at<:toDate",
		"country":  "country=:country",
		"distance": "distance<:distance"}

	where := []string{}
	args := map[string]interface{}{}

	for param, condition := range conditions {
		if _, ok := params[param]; ok {
			where = append(where, condition)
			args[param] = params[param][0]
		}
	}

	query := `SELECT mark, visited_at, place
	FROM visits
	INNER JOIN locations ON visits.location = locations.id
	WHERE ` + strings.Join(where, " AND ")

	places := Places{[]*Place{}}
	nstmt, err := DB.PrepareNamed(query)
	if err != nil {
		return places, err
	}
	err = nstmt.Select(&places.Rows, args)
	return places, err
}

// InsertUser inserts specified user into database.
func InsertUser(user *User) error {
	_, err := DB.NamedExec(
		`INSERT INTO users (id, email, first_name, last_name, gender, birth_date) 
		 VALUES (:id, :email, :first_name, :last_name, :gender, :birth_date)`, user)
	return err
}

// PopulateUsers inserts specified list of users into database.
func PopulateUsers(users Users) error {
	for _, user := range users.Rows {
		if err := InsertUser(user); err != nil {
			return err
		}
	}
	return nil
}
