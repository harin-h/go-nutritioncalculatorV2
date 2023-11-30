package repository

import "time"

type Menu struct {
	Id               int       `db:"id"`
	Name             string    `db:"name"`
	Protein          float64   `db:"protein"`
	Fat              float64   `db:"fat"`
	Carb             float64   `db:"carb"`
	CreatorId        string    `db:"creator_id"`
	CreatorName      string    `db:"creator_name"`
	Like             int       `db:"count_like"`
	Status           int       `db:"status"`
	CreatedTimestamp time.Time `db:"created_timestamp"`
}

type MenuRepository interface {
	CreateMenu(Menu) (*Menu, error)
	GetAllMenues() ([]Menu, error)
	GetMenuById(int) (*Menu, error)
	UpdateMenu(Menu) error
}
