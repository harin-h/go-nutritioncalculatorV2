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

func TestCreateFavList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("CreateFavList", service.NewFavListRequest{
			UserId: "gooddy20",
			Name:   "Breakfast",
			List:   "9,10",
		}).Return(nil)
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.CreateFavList).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id": "gooddy20",
			"name":    "Breakfast",
			"list":    "9,10",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/favlist/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("CreateFavList", service.NewFavListRequest{
			UserId: "gooddy20",
			Name:   "Breakfast",
			List:   "9,10",
		}).Return(nil)
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.CreateFavList).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id": "gooddy20",
			"name":    "Breakfast",
			"list":    "9,10",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/favlist/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateFavList")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.CreateFavList).Methods("POST")
		reqBody := []byte("")
		req := httptest.NewRequest("POST", "/favlist/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateFavList")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("CreateFavList", service.NewFavListRequest{
			UserId: "gooddy20",
			Name:   "Breakfast",
			List:   "9,10",
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.CreateFavList).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id": "gooddy20",
			"name":    "Breakfast",
			"list":    "9,10",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/favlist/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestDeleteFavList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("DeleteFavList", 1).Return(nil)
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/{favlist_id}", hdlr.DeleteFavList).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/favlist/1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Parse Id (String to Int) Error", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("DeleteFavList", 1).Return(nil)
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/{favlist_id}", hdlr.DeleteFavList).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/favlist/1.1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Parse data type error", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "DeleteFavList")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("DeleteFavList", 1).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/{favlist_id}", hdlr.DeleteFavList).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/favlist/1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestUpdateFavList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("UpdateFavList", service.UpdateFavListRequest{
			Id:   1,
			Name: "Extra Breakfast",
			List: "9,9,10",
		}).Return(nil)
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.UpdateFavList).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":   1,
			"name": "Extra Breakfast",
			"list": "9,9,10",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/favlist/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.UpdateFavList).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":   1,
			"name": "Extra Breakfast",
			"list": "9,9,10",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/favlist/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.UpdateFavList).Methods("PUT")
		reqBody := []byte("")
		req := httptest.NewRequest("PUT", "/favlist/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("UpdateFavList", service.UpdateFavListRequest{
			Id:   1,
			Name: "Extra Breakfast",
			List: "9,9,10",
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/", hdlr.UpdateFavList).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":   1,
			"name": "Extra Breakfast",
			"list": "9,9,10",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/favlist/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestGetFavListsByUserId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("GetFavListsByUserId", "gooddy20").Return([]service.FavListResponse{
			{Id: 1, Name: "Breakfast", Menues: "Moo Yang-2 ,Sticky Rice-1 ", List: "9,9,10", Protein: 40, Fat: 10, Carb: 20, IsUpdated: 1},
			{Id: 2, Name: "Lunch", Menues: "Ramyeon-1 ,Moo Yang-1 ", List: "1,9", Protein: 30, Fat: 20, Carb: 63, IsUpdated: 1},
		}, nil)
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/{user_id}", hdlr.GetFavListsByUserId).Methods("GET")
		req := httptest.NewRequest("GET", "/favlist/gooddy20", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		resultBody := []service.FavListResponse{}
		_ = json.Unmarshal(res.Body.Bytes(), &resultBody)
		expectedBody := []service.FavListResponse{
			{Id: 1, Name: "Breakfast", Menues: "Moo Yang-2 ,Sticky Rice-1 ", List: "9,9,10", Protein: 40, Fat: 10, Carb: 20, IsUpdated: 1},
			{Id: 2, Name: "Lunch", Menues: "Ramyeon-1 ,Moo Yang-1 ", List: "1,9", Protein: 30, Fat: 20, Carb: 63, IsUpdated: 1},
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, resultBody)
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewFavListServiceMock()
		srv.On("GetFavListsByUserId", "gooddy20").Return([]service.FavListResponse{}, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewFavListHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/favlist/{user_id}", hdlr.GetFavListsByUserId).Methods("GET")
		req := httptest.NewRequest("GET", "/favlist/gooddy20", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}
