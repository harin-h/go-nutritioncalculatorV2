package repository

import "time"

type Record struct {
	Id               int       `db:"id"`
	UserId           string    `db:"user_id"`
	List             string    `db:"list"`
	Menues           string    `db:"menues"`
	Note             string    `db:"note"`
	Weight           float64   `db:"weight"`
	Protein          float64   `db:"protein"`
	Fat              float64   `db:"fat"`
	Carb             float64   `db:"carb"`
	EventTimestamp   time.Time `db:"event_timestamp"`
	Status           int       `db:"status"`
	IsUpdated        int       `db:"is_updated"`
	CreatedTimestamp time.Time `db:"created_timestamp"`
}

type RecordRepository interface {
	GetRecordsByUserId(string) ([]Record, error)
	GetRecordById(int) (*Record, error)
	CreateRecord(Record) (*Record, error)
	UpdateRecord(Record) error
}
