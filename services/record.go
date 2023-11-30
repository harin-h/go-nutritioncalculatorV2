package service

import "time"

type NewRecordRequest struct {
	UserId         string  `json:"user_id" example:"gooddy20" binding:"required"`                    // "User Id" that create this "Record"
	List           string  `json:"list" example:"9,9,10" binding:"required"`                         // Summary meal with "Menu"'s id e.g. "9,9,10" -> 9 = "Moo Yang" and 10 = "Sticky Rice" so the "Record" contain "Moo Yang" 2 ea and "Sticky Rice" 1 ea
	Note           string  `json:"note" example:"Breakfast"`                                         // Note for this "Record"
	Weight         float64 `json:"weight" example:"63"`                                              // Weight (kg.) that you are on that day
	EventTimestamp string  `json:"event_timestamp" example:"2023-11-01 09:30:00" binding:"required"` // Timestamp that you eat *format="2023-01-01 00:00:00"
}

type UpdateRecordRequest struct {
	Id             int     `json:"id" example:"1" binding:"required"`             // "Record"'s id that you want to update
	List           string  `json:"list" example:"9,9,10"`                         // Summary meal with "Menu"'s id that you want to change to e.g. "9,9,10" -> 9 = "Moo Yang" and 10 = "Sticky Rice" so the "Record" contain "Moo Yang" 2 ea and "Sticky Rice" 1 ea
	Note           string  `json:"note" example:"Lunch"`                          // Note that you want to change to
	Weight         float64 `json:"weight" example:"63"`                           // Weight (kg.) that you want to change to
	EventTimestamp string  `json:"event_timestamp" example:"2023-11-01 12:30:00"` // Timestamp that you want to change to *format="2023-01-01 00:00:00"
}

type RecordResponse struct {
	Id             int       `db:"id"`              // "Record"'s id
	List           string    `db:"list"`            // Summary meal with "Menu"'s id e.g. "9,9,10" -> 9 = "Moo Yang" and 10 = "Sticky Rice" so the "Record" contain "Moo Yang" 2 ea and "Sticky Rice" 1 ea
	Note           string    `db:"note"`            // Note for the "Record"
	Weight         float64   `db:"weight"`          // Weight (kg.) that you are on that day
	Protein        float64   `db:"protein"`         // Total protein (g.) of the "Record"
	Fat            float64   `db:"fat"`             // Total fat (g.) of the "Record"
	Carb           float64   `db:"carb"`            // Total carb (g.) of the "Record"
	EventTimestamp time.Time `db:"event_timestamp"` // Timestamp that you eat *format="2023-01-01 00:00:00"
	IsUpdated      int       `db:"is_updated"`      // 1 = All "Menu" in the "Record" are up to date, 0 = atleast one "Menu" in the "Record" are not up to date
}

type RecordService interface {
	GetAllRecordsByUserId(string) ([]RecordResponse, error)
	CreateRecord(NewRecordRequest) error
	DeleteRecord(int) error
	UpdateRecord(UpdateRecordRequest) error
}
