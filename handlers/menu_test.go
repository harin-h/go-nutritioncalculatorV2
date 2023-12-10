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

func TestCreateMenu(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("CreateMenu", service.NewMenuRequest{
			Name:      "Ramyeon",
			Protein:   5,
			Fat:       15,
			Carb:      65,
			CreatorId: "gooddy20",
		}).Return(nil)
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.CreateMenu).Methods("POST")
		preReqBody := map[string]interface{}{
			"name":       "Ramyeon",
			"protein":    5,
			"fat":        15,
			"carb":       65,
			"creator_id": "gooddy20",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/menu/", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.CreateMenu).Methods("POST")
		preReqBody := map[string]interface{}{
			"name":       "Ramyeon",
			"protein":    5,
			"fat":        15,
			"carb":       65,
			"creator_id": "gooddy20",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/menu/", bytes.NewReader(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateMenu")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.CreateMenu).Methods("POST")
		reqBody := []byte("")
		req := httptest.NewRequest("POST", "/menu/", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateMenu")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("CreateMenu", service.NewMenuRequest{
			Name:      "Ramyeon",
			Protein:   5,
			Fat:       15,
			Carb:      65,
			CreatorId: "gooddy20",
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.CreateMenu).Methods("POST")
		preReqBody := map[string]interface{}{
			"name":       "Ramyeon",
			"protein":    5,
			"fat":        15,
			"carb":       65,
			"creator_id": "gooddy20",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/menu/", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestDeleteMenu(t *testing.T) {
	t.Run("Complete", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("DeleteMenu", 1).Return(nil)
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/{menu_id}", hdlr.DeleteMenu).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/menu/1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Parse Int Error", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("DeleteMenu", 1).Return(nil)
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/{menu_id}", hdlr.DeleteMenu).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/menu/1.1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Parse data type error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("DeleteMenu", 1).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/{menu_id}", hdlr.DeleteMenu).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/menu/1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestUpdateMenu(t *testing.T) {
	t.Run("Complete", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("UpdateMenu", service.UpdateMenuRequest{
			Id:      1,
			Name:    "Ramyeon v2",
			Protein: 8,
			Fat:     20,
			Carb:    70,
		}).Return(nil)
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.UpdateMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":      1,
			"name":    "Ramyeon v2",
			"protein": 8,
			"fat":     20,
			"carb":    70,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/menu/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("content-type is not correct", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("UpdateMenu", service.UpdateMenuRequest{
			Id:      1,
			Name:    "Ramyeon v2",
			Protein: 8,
			Fat:     20,
			Carb:    70,
		}).Return(nil)
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.UpdateMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":      1,
			"name":    "Ramyeon v2",
			"protein": 8,
			"fat":     20,
			"carb":    70,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/menu/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
	})
	t.Run("Decode Request Body Error", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("UpdateMenu", service.UpdateMenuRequest{
			Id:      1,
			Name:    "Ramyeon v2",
			Protein: 8,
			Fat:     20,
			Carb:    70,
		}).Return(nil)
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.UpdateMenu).Methods("PUT")
		reqBody := []byte("")
		req := httptest.NewRequest("PUT", "/menu/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("UpdateMenu", service.UpdateMenuRequest{
			Id:      1,
			Name:    "Ramyeon v2",
			Protein: 8,
			Fat:     20,
			Carb:    70,
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.UpdateMenu).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":      1,
			"name":    "Ramyeon v2",
			"protein": 8,
			"fat":     20,
			"carb":    70,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/menu/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestGetAllMenues(t *testing.T) {
	t.Run("Complete", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("GetAllMenues").Return([]service.MenuResponse{
			{Id: 1,
				Name:        "Ramyeon v2",
				Protein:     8,
				Fat:         20,
				Carb:        70,
				CreatorId:   "gooddy20",
				CreatorName: "GoodDy",
				Like:        5,
				Status:      1},
			{Id: 2,
				Name:        "Fried Chicken",
				Protein:     10,
				Fat:         10,
				Carb:        5,
				CreatorId:   "gooddy20",
				CreatorName: "GoodDy",
				Like:        3,
				Status:      1},
		}, nil)
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.GetAllMenues).Methods("GET")
		req := httptest.NewRequest("GET", "/menu/", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		resultBody := []service.MenuResponse{}
		_ = json.Unmarshal(res.Body.Bytes(), &resultBody)
		expectedBody := []service.MenuResponse{
			{Id: 1,
				Name:        "Ramyeon v2",
				Protein:     8,
				Fat:         20,
				Carb:        70,
				CreatorId:   "gooddy20",
				CreatorName: "GoodDy",
				Like:        5,
				Status:      1},
			{Id: 2,
				Name:        "Fried Chicken",
				Protein:     10,
				Fat:         10,
				Carb:        5,
				CreatorId:   "gooddy20",
				CreatorName: "GoodDy",
				Like:        3,
				Status:      1},
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, resultBody)
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewMenuServiceMock()
		srv.On("GetAllMenues").Return([]service.MenuResponse{
			{Id: 1,
				Name:        "Ramyeon v2",
				Protein:     8,
				Fat:         20,
				Carb:        70,
				CreatorId:   "gooddy20",
				CreatorName: "GoodDy",
				Like:        5,
				Status:      1},
			{Id: 2,
				Name:        "Fried Chicken",
				Protein:     10,
				Fat:         10,
				Carb:        5,
				CreatorId:   "gooddy20",
				CreatorName: "GoodDy",
				Like:        3,
				Status:      1},
		}, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewMenuHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/menu/", hdlr.GetAllMenues).Methods("GET")
		req := httptest.NewRequest("GET", "/menu/", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}
