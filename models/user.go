package models

import (
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
)

const usersTableName = "users"

var usersTableColumns = []string{"id", "email", "first_name", "last_name", "gender", "birth_date"}

// User contains information about user.
type User struct {
	ID        *uint32 `json:"id" db:"id"`
	Email     *string `json:"email" db:"email"`
	FirstName *string `json:"first_name" db:"first_name"`
	LastName  *string `json:"last_name" db:"last_name"`
	Gender    *string `json:"gender" db:"gender"`
	BirthDate *int32  `json:"birth_date" db:"birth_date"`
}

// Users contains list of users.
type Users struct {
	Rows []*User `json:"users"`
}

// Place conrains short description about visited place.
type Place struct {
	Mark      uint8  `json:"mark" db:"mark"`
	VisitedAt int32  `json:"visited_at" db:"visited_at"`
	Place     string `json:"place" db:"place"`
}

// Places contains list of visited places.
type Places struct {
	Rows []*Place `json:"visits"`
}

// PlaceFilter contains filtering parameters for visited places.
type PlaceFilter struct {
	FromDate *int32  `schema:"fromDate"`
	ToDate   *int32  `schema:"toDate"`
	Country  *string `schema:"country"`
	Distance *uint32 `schema:"distance"`
}

// GetUser returns user from database specified by id.
func GetUser(id string) (*User, error) {
	user := new(User)
	err := GetByID(usersTableName, id, user)
	return user, err
}

// GetUserVisits returns user's visits from database specified by user's id.
func GetUserVisits(id string, filter *PlaceFilter) (*Places, error) {
	places := psql.
		Select("mark", "visited_at", "place").
		From(visitsTableName).
		Join(fmt.Sprintf("%s ON %s.location = %s.id", locationsTableName, visitsTableName, locationsTableName)).
		Where(sq.Eq{`"user"`: id})

	if filter.FromDate != nil {
		places = places.Where(sq.Gt{"visited_at": filter.FromDate})
	}
	if filter.ToDate != nil {
		places = places.Where(sq.Lt{"visited_at": filter.ToDate})
	}

	if filter.Country != nil {
		places = places.Where(sq.Eq{"country": filter.Country})
	}
	if filter.Distance != nil {
		places = places.Where(sq.Lt{"distance": filter.Distance})
	}

	sql, args, err := places.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := &Places{[]*Place{}}
	if err := DB.Select(&result.Rows, sql, args...); err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

// InsertUser inserts specified user into database.
func InsertUser(user *User) error {
	sql, args, err := psql.
		Insert(usersTableName).
		Columns(usersTableColumns...).
		Values(user.ID, user.Email, user.FirstName, user.LastName, user.Gender, user.BirthDate).
		ToSql()

	_, err = DB.Exec(sql, args...)
	if err != nil {
		log.Println(err)
	}
	return err
}

// PopulateUsers inserts specified list of users into database.
func PopulateUsers(users *Users) error {
	for _, user := range users.Rows {
		if err := InsertUser(user); err != nil {
			return err
		}
	}
	return nil
}

// UpdateUser updates specified user's row in database.
func UpdateUser(id string, user *User) error {
	update := psql.Update(usersTableName)

	if user.Email != nil {
		update = update.Set("email", user.Email)
	}
	if user.FirstName != nil {
		update = update.Set("first_name", user.FirstName)
	}
	if user.LastName != nil {
		update = update.Set("last_name", user.LastName)
	}
	if user.Gender != nil {
		update = update.Set("gender", user.Gender)
	}
	if user.BirthDate != nil {
		update = update.Set("birth_date", user.BirthDate)
	}

	update = update.Where(sq.Eq{"id": id})

	sql, args, err := update.ToSql()

	_, err = DB.Exec(sql, args...)
	if err != nil {
		log.Println(err)
	}

	return err
}
