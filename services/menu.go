package service

type NewMenuRequest struct {
	Name      string  `json:"name" example:"7-11 Pepper Chicken Breast" binding:"required"` // Name of this "Menu"
	Protein   float64 `json:"protein" example:"19" binding:"required"`                      // Protein (g.) of this "Menu"
	Fat       float64 `json:"fat" example:"0.5" binding:"required"`                         // Fat (g.) of this "Menu"
	Carb      float64 `json:"carb" example:"0" binding:"required"`                          // Carb (g.) of this "Menu"
	CreatorId string  `json:"creator_id" example:"gooddy20" binding:"required"`             // "User Id" that create this "Menu"
}

type UpdateMenuRequest struct {
	Id      int     `json:"id" example:"1" binding:"required"`                            // "Menu"'s id that you want to update
	Name    string  `json:"name" example:"7-11 Chilli Chicken Breast" binding:"required"` // The name that you want to change to
	Protein float64 `json:"protein" example:"20" binding:"required"`                      // The protein (g.) that you want to change to
	Fat     float64 `json:"fat" example:"0.5" binding:"required"`                         // The fat (g.) that you want to change to
	Carb    float64 `json:"carb" example:"1" binding:"required"`                          // The carb (g.) that you want to change to
}

type MenuResponse struct {
	Id          int     `json:"id" example:"9"`                // "Menu"'s id that generate by system
	Name        string  `json:"name" example:"Moo Yang"`       // Name of "Menu" that named by the user
	Protein     float64 `json:"protein" example:"20"`          // Protein of "Menu"
	Fat         float64 `json:"fat" example:"5"`               // Fat of "Menu"
	Carb        float64 `json:"carb" example:"0"`              // Carb of "Menu"
	CreatorId   string  `json:"creator_id" example:"gooddy20"` // "User Id" that create the "Menu"
	CreatorName string  `json:"creator_name" example:"GoodDy"` // "Username" that create the "Menu"
	Like        int     `json:"like" example:"1"`              // Amount of using as favorite menu by "User Id"
	Status      int     `json:"status" example:"1"`            // 1 = Active, 0 = Deleted
}

type MenuService interface {
	CreateMenu(NewMenuRequest) error
	GetAllMenues() ([]MenuResponse, error)
	UpdateMenu(UpdateMenuRequest) error
	RecoverMenu(int, string) (*MenuResponse, error)
	DeleteMenu(int) error
}
