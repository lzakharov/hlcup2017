package models

import "log"

// User contains information about user.
type User struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate int32  `json:"birth_date"`
}

// Users contains list of users.
type Users struct {
	Rows []*User `json:"users"`
}

// InsertUser inserts specified user into database.
func InsertUser(user *User) {
	_, err := DB.NamedExec(
		`INSERT INTO users (id, email, first_name, last_name, gender, birth_date) 
		 VALUES (:id, :email, :firstname, :lastname, :gender, :birthdate)`, user)
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
