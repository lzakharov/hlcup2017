package models

import "log"

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
func GetUser(id uint32) (User, error) {
	user := User{}
	if err := DB.Get(&user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return user, err
	}
	return user, nil
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
