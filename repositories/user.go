package repository

import "time"

type User struct {
	UserId           string    `db:"user_id"`
	Password         string    `db:"password"`
	Username         string    `db:"username"`
	Weight           float64   `db:"weight"`
	Protein          float64   `db:"protein"`
	Fat              float64   `db:"fat"`
	Carb             float64   `db:"carb"`
	FavoriteMenues   string    `db:"favorite_menues"`
	CreatedTimestamp time.Time `db:"created_timestamp"`
}

type UserRepository interface {
	GetUserById(string) (*User, error)
	GetUserByUsername(string) (*User, error)
	CreateUser(User) error
	UpdateUser(User) error
}
