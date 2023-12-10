package handler_test

import (
	"bytes"
	"encoding/json"
	"go-nutritioncalculator2/errs"
	handler "go-nutritioncalculator2/handlers"
	service "go-nutritioncalculator2/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRecoverDeletedMenu(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		menuSrv := service.NewMenuServiceMock()
		menuSrv.On("RecoverMenu", 1, "ramyeon v2").Return(&service.MenuResponse{
			Id:          15,
			Name:        "ramyeon v2",
			Protein:     10,
			Fat:         10,
			Carb:        50,
			CreatorId:   "gooddy20",
			CreatorName: "GoodDy",
			Like:        0,
			Status:      1,
		}, nil)
		userSrv := service.NewUserServiceMock()
		userSrv.On("RecoverFavoriteMenues", "gooddy20", 1).Return(nil)
		favListSrv := service.NewFavListServiceMock()
		favListSrv.On("GetFavListsByUserId", "gooddy20").Return([]service.FavListResponse{
			{Id: 1, Name: "Breakfast", Menues: "Moo Yang-2 ,Sticky Rice-1 ", List: "9,9,10", Protein: 40, Fat: 10, Carb: 20, IsUpdated: 1},
			{Id: 2, Name: "Lunch", Menues: "Ramyeon-1 ,Moo Yang-1 ", List: "1,9", Protein: 30, Fat: 20, Carb: 63, IsUpdated: 1},
		}, nil)
		favListSrv.On("RecoverFavList", 1, 1, 15).Return(nil)
		favListSrv.On("RecoverFavList", 2, 1, 15).Return(nil)
		hdlr := handler.NewMultiHandler(menuSrv, userSrv, favListSrv)
		r := mux.NewRouter()
		r.HandleFunc("/recover/", hdlr.RecoverDeletedMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"deleted_menu_id": 1,
			"new_menu_name":   "ramyeon v2",
			"is_create":       1,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/recover/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		menuSrv := service.NewMenuServiceMock()
		userSrv := service.NewUserServiceMock()
		favListSrv := service.NewFavListServiceMock()
		hdlr := handler.NewMultiHandler(menuSrv, userSrv, favListSrv)
		r := mux.NewRouter()
		r.HandleFunc("/recover/", hdlr.RecoverDeletedMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"deleted_menu_id": 1,
			"new_menu_name":   "ramyeon v2",
			"is_create":       1,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/recover/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		menuSrv.AssertNotCalled(t, "RecoverMenu")
		userSrv.AssertNotCalled(t, "RecoverFavoriteMenues")
		favListSrv.AssertNotCalled(t, "GetFavListsByUserId")
		favListSrv.AssertNotCalled(t, "RecoverFavList")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		menuSrv := service.NewMenuServiceMock()
		userSrv := service.NewUserServiceMock()
		favListSrv := service.NewFavListServiceMock()
		hdlr := handler.NewMultiHandler(menuSrv, userSrv, favListSrv)
		r := mux.NewRouter()
		r.HandleFunc("/recover/", hdlr.RecoverDeletedMenu).Methods("PUT")
		reqBody := []byte("")
		req := httptest.NewRequest("PUT", "/recover/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		menuSrv.AssertNotCalled(t, "RecoverMenu")
		userSrv.AssertNotCalled(t, "RecoverFavoriteMenues")
		favListSrv.AssertNotCalled(t, "GetFavListsByUserId")
		favListSrv.AssertNotCalled(t, "RecoverFavList")
	})
	t.Run("Menu Service Error", func(t *testing.T) {
		menuSrv := service.NewMenuServiceMock()
		menuSrv.On("RecoverMenu", 1, "ramyeon v2").Return(&service.MenuResponse{}, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		userSrv := service.NewUserServiceMock()
		favListSrv := service.NewFavListServiceMock()
		hdlr := handler.NewMultiHandler(menuSrv, userSrv, favListSrv)
		r := mux.NewRouter()
		r.HandleFunc("/recover/", hdlr.RecoverDeletedMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"deleted_menu_id": 1,
			"new_menu_name":   "ramyeon v2",
			"is_create":       1,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/recover/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
		userSrv.AssertNotCalled(t, "RecoverFavoriteMenues")
		favListSrv.AssertNotCalled(t, "GetFavListsByUserId")
		favListSrv.AssertNotCalled(t, "RecoverFavList")
	})
	t.Run("User Service Error", func(t *testing.T) {
		menuSrv := service.NewMenuServiceMock()
		menuSrv.On("RecoverMenu", 1, "ramyeon v2").Return(&service.MenuResponse{
			Id:          15,
			Name:        "ramyeon v2",
			Protein:     10,
			Fat:         10,
			Carb:        50,
			CreatorId:   "gooddy20",
			CreatorName: "GoodDy",
			Like:        0,
			Status:      1,
		}, nil)
		userSrv := service.NewUserServiceMock()
		userSrv.On("RecoverFavoriteMenues", "gooddy20", 1).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		favListSrv := service.NewFavListServiceMock()
		hdlr := handler.NewMultiHandler(menuSrv, userSrv, favListSrv)
		r := mux.NewRouter()
		r.HandleFunc("/recover/", hdlr.RecoverDeletedMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"deleted_menu_id": 1,
			"new_menu_name":   "ramyeon v2",
			"is_create":       1,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/recover/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
		favListSrv.AssertNotCalled(t, "GetFavListsByUserId")
		favListSrv.AssertNotCalled(t, "RecoverFavList")
	})
	t.Run("Get Favorite List Service Error", func(t *testing.T) {
		menuSrv := service.NewMenuServiceMock()
		menuSrv.On("RecoverMenu", 1, "ramyeon v2").Return(&service.MenuResponse{
			Id:          15,
			Name:        "ramyeon v2",
			Protein:     10,
			Fat:         10,
			Carb:        50,
			CreatorId:   "gooddy20",
			CreatorName: "GoodDy",
			Like:        0,
			Status:      1,
		}, nil)
		userSrv := service.NewUserServiceMock()
		userSrv.On("RecoverFavoriteMenues", "gooddy20", 1).Return(nil)
		favListSrv := service.NewFavListServiceMock()
		favListSrv.On("GetFavListsByUserId", "gooddy20").Return([]service.FavListResponse{}, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewMultiHandler(menuSrv, userSrv, favListSrv)
		r := mux.NewRouter()
		r.HandleFunc("/recover/", hdlr.RecoverDeletedMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"deleted_menu_id": 1,
			"new_menu_name":   "ramyeon v2",
			"is_create":       1,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/recover/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
		favListSrv.AssertNotCalled(t, "RecoverFavList")
	})
	t.Run("Recover Favorite List Service Error", func(t *testing.T) {
		menuSrv := service.NewMenuServiceMock()
		menuSrv.On("RecoverMenu", 1, "ramyeon v2").Return(&service.MenuResponse{
			Id:          15,
			Name:        "ramyeon v2",
			Protein:     10,
			Fat:         10,
			Carb:        50,
			CreatorId:   "gooddy20",
			CreatorName: "GoodDy",
			Like:        0,
			Status:      1,
		}, nil)
		userSrv := service.NewUserServiceMock()
		userSrv.On("RecoverFavoriteMenues", "gooddy20", 1).Return(nil)
		favListSrv := service.NewFavListServiceMock()
		favListSrv.On("GetFavListsByUserId", "gooddy20").Return([]service.FavListResponse{
			{Id: 1, Name: "Breakfast", Menues: "Moo Yang-2 ,Sticky Rice-1 ", List: "9,9,10", Protein: 40, Fat: 10, Carb: 20, IsUpdated: 1},
			{Id: 2, Name: "Lunch", Menues: "Ramyeon-1 ,Moo Yang-1 ", List: "1,9", Protein: 30, Fat: 20, Carb: 63, IsUpdated: 1},
		}, nil)
		favListSrv.On("RecoverFavList", 1, 1, 15).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		favListSrv.On("RecoverFavList", 2, 1, 15).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewMultiHandler(menuSrv, userSrv, favListSrv)
		r := mux.NewRouter()
		r.HandleFunc("/recover/", hdlr.RecoverDeletedMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"deleted_menu_id": 1,
			"new_menu_name":   "ramyeon v2",
			"is_create":       1,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/recover/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}
