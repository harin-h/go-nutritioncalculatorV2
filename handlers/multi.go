package handler

import (
	"encoding/json"
	"go-nutritioncalculator2/errs"
	service "go-nutritioncalculator2/services"
	"net/http"
)

type multiHandler struct {
	menuSrv    service.MenuService
	userSrv    service.UserService
	favListSrv service.FavListService
}

type MultiRequest struct {
	UserId        string `json:"user_id" example:"gooddy20" binding:"required"`  // "User Id" that want to recover the deleted "Menu"
	DeletedMenuId int    `json:"deleted_menu_id" example:"9" binding:"required"` // "Menu"'s id that was deleted
	NewMenuName   string `json:"new_menu_name" example:"Moo Yang V2"`            // New name of recovered "Menu"
	IsCreate      int    `json:"is_create" example:"1" binding:"required"`       // 1 = Want to create new "Menu" for replace "Menu" in the "Favorite List", 0 = Dont want to create new "Menu" so the "Favorite List" that contain the deleted "Menu" will be updated by get the "Menu" off
}

func NewMultiHandler(menuSrv service.MenuService, userSrv service.UserService, favListSrv service.FavListService) multiHandler {
	return multiHandler{menuSrv: menuSrv, userSrv: userSrv, favListSrv: favListSrv}
}

// RecoverDeletedMenu ... Recover a deleted "Menu"
// @Summary Recover a deleted "Menu"
// @Description Get the deleted `Menu` off from `Favorite Menu` and {1. replace the deleted `Menu` in `Favorite List` with the new `Menu` that has the same detail (Can change the "Menu"'s name) / 2. get the deleted `Menu` off from `Favorite List`}
// @Tags Recover
// @Accept json
// @Param request body MultiRequest true "The data detail that you want"
// @Response 200
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /recover/ [put]
func (h multiHandler) RecoverDeletedMenu(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request MultiRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	userId := request.UserId
	deletedMenuId := request.DeletedMenuId
	newMenuName := request.NewMenuName
	isCreateMenu := request.IsCreate
	var newMenuId int
	if isCreateMenu == 1 {
		newMenu, err := h.menuSrv.RecoverMenu(deletedMenuId, newMenuName)
		if err != nil {
			handlerError(w, err)
			return
		}
		newMenuId = newMenu.Id
	}
	err = h.userSrv.RecoverFavoriteMenues(userId, deletedMenuId)
	if err != nil {
		handlerError(w, err)
		return
	}
	favLists, err := h.favListSrv.GetFavListsByUserId(userId)
	if err != nil {
		handlerError(w, err)
		return
	}
	for i := 0; i < len(favLists); i++ {
		err = h.favListSrv.RecoverFavList(favLists[i].Id, deletedMenuId, newMenuId)
		if err != nil {
			handlerError(w, err)
			return
		}
	}
}
