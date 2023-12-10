package handler

import (
	"encoding/json"
	"go-nutritioncalculator2/errs"
	service "go-nutritioncalculator2/services"
	"net/http"

	"github.com/gorilla/mux"
)

type userHandler struct {
	userSrv service.UserService
}

func NewUserHandler(userSrv service.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}

// LogIn ... Check "User Id" and "Password" are correct or not
// @Summary Check "User Id" and "Password" are correct or not
// @Description Check `User Id` and `Password` are correct or not
// @Tags User
// @Accept json
// @Produce json
// @Param request body service.LogInRequest true "`User Id` and `Password`"
// @Response 200 {object} service.LogInResponse
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /user/login [put]
func (h userHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request service.LogInRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	isLogIn, err := h.userSrv.CheckLogIn(request)
	if err != nil {
		handlerError(w, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(isLogIn)
}

// CreateUser ... Create a "User"
// @Summary Create a "User"
// @Description Create a `User`
// @Tags User
// @Accept json
// @Param request body service.NewUserRequest true "`User`'s data detail"
// @Response 200
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /user/ [post]
func (h userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request service.NewUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	err = h.userSrv.CreateUser(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}

// GetUserDetail ... Get a "User"'s detail
// @Summary Get a "User"'s detail
// @Description Get a `User`'s detail by `User Id`
// @Tags User
// @Param user_id path string true "`User Id`"
// @Response 200 {object} service.UserResponse
// @Response 406 "`User Id` is not found"
// @Response 500 "Internal Server Error"
// @Router /user/{user_id} [get]
func (h userHandler) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := h.userSrv.GetUserDetail(vars["user_id"])
	if err != nil {
		handlerError(w, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUserDetail ... Update a "User"'s detail
// @Summary Update a "User"'s detail
// @Description Update a `User`'s detail
// @Tags User
// @Param request body service.UpdateUserRequest true "`User`'s data detail that you want to update and can ignore the unchanged parameters"
// @Response 200
// @Response 406 "Request Body Not Acceptable or `User Id` is not found"
// @Response 500 "Internal Server Error"
// @Router /user/userdetail [put]
func (h userHandler) UpdateUserDetail(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request service.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	err = h.userSrv.UpdateUser(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}
