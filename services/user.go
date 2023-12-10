package service

type NewUserRequest struct {
	UserId   string  `json:"user_id" example:"gooddy20" binding:"required"`      // "User Id"
	Password string  `json:"password" example:"zxc123zxc123" binding:"required"` // "Password"
	Username string  `json:"username" example:"GoodDy" binding:"required"`       // "Username"
	Weight   float64 `json:"weight" example:"70"`                                // Default weight (kg.) of the "User"
	Protein  float64 `json:"protein" example:"120"`                              // Default protein (g.) of the "User"
	Fat      float64 `json:"fat" example:"60"`                                   // Default fat (g.) of the "User"
	Carb     float64 `json:"carb" example:"120"`                                 // Default carb (g.) of the "User"
}

type UpdateUserRequest struct {
	UserId         string  `json:"user_id" example:"gooddy20" binding:"required"` // "User Id"
	Password       string  `json:"password" example:"zxc123zxc456"`               // "Password" that you want to change
	Username       string  `json:"username" example:"GooDDy19"`                   // "Username" that you want to change to
	Weight         float64 `json:"weight" example:"72"`                           // Weight (kg.) that you want to change to
	Protein        float64 `json:"protein" example:"150"`                         // Protein (g.) that you want to change to
	Fat            float64 `json:"fat" example:"70"`                              // Fat (g.) that you want to change to
	Carb           float64 `json:"carb" example:"160"`                            // Carb that you want to change to
	FavoriteMenues string  `json:"favorite_menues" example:"4,7,9,10,11"`         // Favorite Menues's id that you want to change to e.g. "9,10" 9 = "Moo Yang" and 10 = "Sticky Rice" so this "User" got "Moo Yang" and "Sticky Rice" as "Favorite Menu"
}

type UserResponse struct {
	Username       string  `json:"username" example:"GoodDy"`      // "Username"
	Weight         float64 `json:"weight" example:"62"`            // Default weight (kg.) of the "User"
	Protein        float64 `json:"protein" example:"140"`          // Default protein (g.) of the "User"
	Fat            float64 `json:"fat" example:"40"`               // Default fat (g.) of the "User"
	Carb           float64 `json:"carb" example:"130"`             // Default carb (g.) of the "User"
	FavoriteMenues string  `json:"favorite_menues" example:"9,10"` // Favorite Menues's id e.g. "9,10" 9 = "Moo Yang" and 10 = "Sticky Rice" so this "User" got "Moo Yang" and "Sticky Rice" as "Favorite Menu"
}

type LogInRequest struct {
	UserId   string `json:"user_id" example:"gooddy20" binding:"required"`      // "User Id"
	Password string `json:"password" example:"zxc123zxc123" binding:"required"` // "Password"
}

type LogInResponse struct {
	IsLogIn bool `json:"IsLogIn" example:"true"` // "true" = Pass, "false" = Incorrect "User Id" or "Password"
}

type UserService interface {
	CheckLogIn(LogInRequest) (*LogInResponse, error)
	GetUserDetail(string) (*UserResponse, error)
	CreateUser(NewUserRequest) error
	UpdateUser(UpdateUserRequest) error
	RecoverFavoriteMenues(string, int) error
}
