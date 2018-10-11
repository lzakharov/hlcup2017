package main

type User struct {
	ID        *uint32 `json:"id" db:"id"`
	Email     *string `json:"email" db:"email"`
	FirstName *string `json:"first_name" db:"first_name"`
	LastName  *string `json:"last_name" db:"last_name"`
	Gender    *string `json:"gender" db:"gender"`
	BirthDate *int32  `json:"birth_date" db:"birth_date"`
}

type Users struct {
	Rows []*User `json:"users"`
}
