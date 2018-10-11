package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	// postgres driver
	_ "github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	Socket           *sqlx.DB
	StatementBuilder sq.StatementBuilderType
}

const (
	connectionTimeout  = 30
	reconnectionTime   = 3
	usersTableName     = "users"
	locationsTableName = "locations"
	visitsTableName    = "visits"
)

var (
	usersTableColumns     = []string{"id", "email", "first_name", "last_name", "gender", "birth_date"}
	locationsTableColumns = []string{"id", "place", "country", "city", "distance"}
	visitsTableColumns    = []string{"id", "location", `"user"`, "visited_at", "mark"}
)

func (d *Database) Initialize(c *DBConfig) error {
	var err error
	timer := time.NewTimer(time.Duration(connectionTimeout) * time.Second)
	connected := make(chan struct{})

	go func() {
		for {
			if d.Socket, err = sqlx.Connect(c.Driver, c.GetDataSourceName()); err == nil {
				connected <- struct{}{}
			}
			log.Println("Database connection failed")
			time.Sleep(reconnectionTime * time.Second)
		}
	}()
	select {
	case <-timer.C:
		return errors.New("database connection timeout")
	case <-connected:
		timer.Stop()
		log.Println("Database connected")
	}
	d.StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return d.createSchema(c.Schema)
}

func (d *Database) createSchema(file string) error {
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	d.Socket.MustExec(string(schema))
	return nil
}

func (d *Database) getByID(table string, id string, dest interface{}) error {
	sql, args, err := d.StatementBuilder.Select("*").From(table).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		log.Println(err)
		return err
	}

	if err := d.Socket.Get(dest, sql, args...); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (d *Database) GetUser(id string) (*User, error) {
	user := new(User)
	err := d.getByID(usersTableName, id, user)
	return user, err
}

func (d *Database) GetUserVisits(id string, filter *PlaceFilter) (*Places, error) {
	places := d.StatementBuilder.
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
	if err := d.Socket.Select(&result.Rows, sql, args...); err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (d *Database) InsertUser(user *User) error {
	sql, args, err := d.StatementBuilder.
		Insert(usersTableName).
		Columns(usersTableColumns...).
		Values(user.ID, user.Email, user.FirstName, user.LastName, user.Gender, user.BirthDate).
		ToSql()
	if err != nil {
		return err
	}

	_, err = d.Socket.Exec(sql, args...)
	return err
}

func (d *Database) PopulateUsers(users *Users) error {
	for _, user := range users.Rows {
		if err := d.InsertUser(user); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) UpdateUser(id string, user *User) error {
	update := d.StatementBuilder.Update(usersTableName)

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
	if err != nil {
		return err
	}

	_, err = d.Socket.Exec(sql, args...)
	return err
}

func (d *Database) GetLocation(id string) (*Location, error) {
	location := new(Location)
	err := d.getByID(locationsTableName, id, location)
	return location, err
}

func (d *Database) GetLocationAverageMark(id string, filter *LocationFilter) (*LocationAvgMark, error) {
	const age = "date_part('year', age(to_timestamp(users.birth_date)))"

	locations := d.StatementBuilder.
		Select(`COALESCE("round"("avg"(visits.mark), 2), 0) AS "avg"`).
		From(locationsTableName).
		Join(fmt.Sprintf("%s ON %s.id = %s.location", visitsTableName, locationsTableName, visitsTableName)).
		Join(fmt.Sprintf(`%s ON %s."user" = %s.id`, usersTableName, visitsTableName, usersTableName)).
		Where(sq.Eq{locationsTableName + ".id": id})

	if filter.FromDate != nil {
		locations = locations.Where(sq.Gt{"visits.visited_at": filter.FromDate})
	}
	if filter.ToDate != nil {
		locations = locations.Where(sq.Lt{"visits.visited_at": filter.ToDate})
	}

	if filter.FromAge != nil {
		locations = locations.Where(sq.Gt{age: filter.FromAge})
	}
	if filter.ToAge != nil {
		locations = locations.Where(sq.Lt{age: filter.ToAge})
	}

	if filter.Gender != nil {
		locations = locations.Where(sq.Eq{"users.gender": filter.Gender})
	}

	sql, args, err := locations.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	average := &LocationAvgMark{}
	if err = d.Socket.Get(average, sql, args...); err != nil {
		log.Println(err)
		return nil, err
	}

	return average, err
}

func (d *Database) InsertLocation(location *Location) error {
	sql, args, err := d.StatementBuilder.
		Insert(locationsTableName).
		Columns(locationsTableColumns...).
		Values(location.ID, location.Place, location.Country, location.City, location.Distance).
		ToSql()
	if err != nil {
		return err
	}

	_, err = d.Socket.Exec(sql, args...)
	return err
}

func (d *Database) PopulateLocations(locations *Locations) error {
	for _, location := range locations.Rows {
		if err := d.InsertLocation(location); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) UpdateLocation(id string, location *Location) error {
	update := d.StatementBuilder.Update(locationsTableName)

	if location.Place != nil {
		update = update.Set("place", location.Place)
	}
	if location.Country != nil {
		update = update.Set("country", location.Country)
	}
	if location.City != nil {
		update = update.Set("city", location.City)
	}
	if location.Distance != nil {
		update = update.Set("distance", location.Distance)
	}

	update = update.Where(sq.Eq{"id": id})

	sql, args, err := update.ToSql()
	if err != nil {
		return err
	}

	_, err = d.Socket.Exec(sql, args...)
	return err
}

func (d *Database) GetVisit(id string) (*Visit, error) {
	visit := new(Visit)
	err := d.getByID(visitsTableName, id, visit)
	return visit, err
}

func (d *Database) InsertVisit(visit *Visit) error {
	sql, args, err := d.StatementBuilder.
		Insert(visitsTableName).
		Columns(visitsTableColumns...).
		Values(visit.ID, visit.Location, visit.User, visit.VisitedAt, visit.Mark).
		ToSql()
	if err != nil {
		return err
	}

	_, err = d.Socket.Exec(sql, args...)
	return err
}

func (d *Database) PopulateVisits(visits *Visits) error {
	for _, visit := range visits.Rows {
		if err := d.InsertVisit(visit); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) UpdateVisit(id string, visit *Visit) error {
	update := d.StatementBuilder.Update(visitsTableName)

	if visit.Location != nil {
		update = update.Set("location", visit.Location)
	}
	if visit.User != nil {
		update = update.Set(`"user"`, visit.User)
	}
	if visit.VisitedAt != nil {
		update = update.Set("visited_at", visit.VisitedAt)
	}
	if visit.Mark != nil {
		update = update.Set("mark", visit.Mark)
	}

	update = update.Where(sq.Eq{"id": id})

	sql, args, err := update.ToSql()
	if err != nil {
		return err
	}

	_, err = d.Socket.Exec(sql, args...)
	return err
}
