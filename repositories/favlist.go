package repository

import "time"

type FavList struct {
	Id               int       `db:"id"`
	UserId           string    `db:"user_id"`
	Name             string    `db:"name"`
	Menues           string    `db:"menues"`
	List             string    `db:"list"`
	Protein          float64   `db:"protein"`
	Fat              float64   `db:"fat"`
	Carb             float64   `db:"carb"`
	Status           int       `db:"status"`
	IsUpdated        int       `db:"is_updated"`
	CreatedTimestamp time.Time `db:"created_timestamp"`
}

type FavListRepository interface {
	GetFavListsByUserId(string) ([]FavList, error)
	GetFavListById(int) (*FavList, error)
	CreateFavList(FavList) (*FavList, error)
	UpdateFavList(FavList) error
}
