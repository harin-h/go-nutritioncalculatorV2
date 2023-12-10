package handler

import (
	"encoding/json"
	"go-nutritioncalculator2/errs"
	service "go-nutritioncalculator2/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type menuHandler struct {
	menuSrv service.MenuService
}

func NewMenuHandler(menuSrv service.MenuService) menuHandler {
	return menuHandler{menuSrv: menuSrv}
}

// CreateMenu ... Create a "Menu"
// @Summary Create a "Menu"
// @Description Create a 'Menu'
// @Tags Menu
// @Accept json
// @Param request body service.NewMenuRequest true "`Menu`'s data detail"
// @Response 200
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /menu/ [post]
func (h menuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request service.NewMenuRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	err = h.menuSrv.CreateMenu(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}

// DeleteMenu ... Delete a "Menu"
// @Summary Delete a "Menu"
// @Description Delete a 'Menu'
// @Tags Menu
// @Param menu_id path int true "`Menu`'s id that you want to delete"
// @Response 200
// @Response 406 "Request Parameter Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /menu/{menu_id} [delete]
func (h menuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menu_id, err := strconv.ParseInt(vars["menu_id"], 0, 0)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Parse data type error"})
		return
	}
	err = h.menuSrv.DeleteMenu(int(menu_id))
	if err != nil {
		handlerError(w, err)
		return
	}
}

// UpdateMenu ... Update a "Menu"
// @Summary Update a "Menu"
// @Description Update a `Menu`
// @Tags Menu
// @Accept json
// @Param request body service.UpdateMenuRequest true "`Menu`'s data detail that you want to update and the unchanged parameters need to be input the old value"
// @Response 200
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /menu/ [put]
func (h menuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request service.UpdateMenuRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	err = h.menuSrv.UpdateMenu(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}

// GetAllMenues ... Get all "Menu"
// @Summary Get all "Menu"
// @Description Get all 'Menu'
// @Tags Menu
// @Produce json
// @Response 200 {object} []service.MenuResponse
// @Response 500 "Internal Server Error"
// @Router /menu/ [get]
func (h menuHandler) GetAllMenues(w http.ResponseWriter, r *http.Request) {
	response, err := h.menuSrv.GetAllMenues()
	if err != nil {
		handlerError(w, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
