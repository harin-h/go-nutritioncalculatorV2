package service

type FavListResponse struct {
	Id        int     `json:"id" example:"1"`                              // "Favorite List"'s id that generate by system
	Name      string  `json:"name" example:"Daily Breakfast"`              // Name of "Favorite List" that named by the user
	Menues    string  `json:"menues" example:"Moo Yang-2, Sticky Rice-1 "` // Summary each "Menu"'s name and amount of the "Favorite List"
	List      string  `json:"list" example:"9,9,10"`                       // Summary meal with "Menu"'s id e.g. "9,9,10" -> 9 = "Moo Yang" and 10 = "Sticky Rice" so the "Favorite List" contain "Moo Yang" 2 ea and "Sticky Rice" 1 ea
	Protein   float64 `json:"protein" example:"40"`                        // Total protein (g.) in the "Favorite List"
	Fat       float64 `json:"fat" example:"10"`                            // Total fat (g.) in the "Favorite List"
	Carb      float64 `json:"carb" example:"20"`                           // Total carb (g.) in the "Favorite List"
	IsUpdated int     `json:"is_updated" example:"1"`                      // 1 = All "Menu" in the "Favorite List" are up to date, 0 = atleast one "Menu" in the "Favorite List" are not up to date
}

type NewFavListRequest struct {
	UserId string `json:"user_id" example:"gooddy20" binding:"required"`     // The "User Id" that create this "Favorite List"
	Name   string `json:"name" example:"Daily Breakfast" binding:"required"` // The name of this "Favorite List"
	List   string `json:"list" example:"9,9,10" binding:"required"`          // Summary meal with "Menu"'s id  e.g. "9,9,10" -> 9 = "Moo Yang" and 10 = "Sticky Rice" so the "Favorite List" contain "Moo Yang" 2 ea and "Sticky Rice" 1 ea
}

type UpdateFavListRequest struct {
	Id   int    `json:"id" example:"1" binding:"required"` // The "Favorite List"'s id that is updated
	Name string `json:"name" example:"Daily Breakfast"`    // The name that you want to change to
	List string `json:"list" example:"9,10"`               // Summary meal with "Menu"'s id that you want to change e.g. "9,9,10" -> 9 = "Moo Yang" and 10 = "Sticky Rice" so the "Favorite List" contain "Moo Yang" 2 ea and "Sticky Rice" 1 ea
}

type FavListService interface {
	GetFavListsByUserId(string) ([]FavListResponse, error)
	CreateFavList(NewFavListRequest) error
	DeleteFavList(int) error
	UpdateFavList(UpdateFavListRequest) error
	RecoverFavList(int, int, int) error
}
