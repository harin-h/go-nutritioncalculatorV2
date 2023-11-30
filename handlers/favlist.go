package handler

import (
	"encoding/json"
	"go-nutritioncalculator2/errs"
	service "go-nutritioncalculator2/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type favListHandler struct {
	favListSrv service.FavListService
}

func NewFavListHandler(favListSrv service.FavListService) favListHandler {
	return favListHandler{favListSrv: favListSrv}
}

// CreateFavList ... Create a "Favorite List"
// @Summary Create a "Favorite List"
// @Description Create a `Favorite List` for recording the daily meal easily
// @Tags Favorite List
// @Accept json
// @Param request body service.NewFavListRequest true "`Favorite List`'s data detail"
// @Response 200
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /favlist/ [post]
func (h favListHandler) CreateFavList(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Request body incorrect format"})
		return
	}
	var request service.NewFavListRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Paste request body error"})
		return
	}
	err = h.favListSrv.CreateFavList(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}

// DeleteFavList ... Delete a "Favorite List"
// @Summary Delete a "Favorite List"
// @Description Delete a `Favorite List`
// @Tags Favorite List
// @Param favlist_id path int true "`Favorite List`'s id that you want to delete"
// @Response 200
// @Response 406 "Request Parameter Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /favlist/{favlist_id} [delete]
func (h favListHandler) DeleteFavList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	favListId, err := strconv.ParseInt(vars["favlist_id"], 0, 0)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Parse data type error"})
		return
	}
	err = h.favListSrv.DeleteFavList(int(favListId))
	if err != nil {
		handlerError(w, err)
		return
	}
}

// UpdateFavList ... Update a "Favorite List"
// @Summary Update a "Favorite List"
// @Description Update a `Favorite List`
// @Tags Favorite List
// @Accept json
// @Param request body service.UpdateFavListRequest true "`Favorite List`'s data detail that you want to update and can ignore the unchanged parameters"
// @Response 200
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /favlist/ [put]
func (h favListHandler) UpdateFavList(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Request body incorrect format"})
		return
	}
	var request service.UpdateFavListRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Paste request body error"})
		return
	}
	err = h.favListSrv.UpdateFavList(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}

// GetFavListsByUserId ... Get all "Favorite List" of the "User Id"
// @Summary Get all "Favorite List" of the "User Id"
// @Description Get all `Favorite List` of the `User Id`
// @Tags Favorite List
// @Produce json
// @Param user_id path string true "User Id"
// @Response 200 {object} []service.FavListResponse
// @Response 500 "Internal Server Error"
// @Router /favlist/{user_id} [get]
func (h favListHandler) GetFavListsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := h.favListSrv.GetFavListsByUserId(vars["user_id"])
	if err != nil {
		handlerError(w, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
