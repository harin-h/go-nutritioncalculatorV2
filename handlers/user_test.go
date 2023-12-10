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

func TestLogIn(t *testing.T) {
	t.Run("Success Case: Correct Password", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("CheckLogIn", service.LogInRequest{
			UserId:   "gooddy20",
			Password: "correctPassword",
		}).Return(&service.LogInResponse{
			IsLogIn: true,
		}, nil)
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/login", hdlr.LogIn).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":  "gooddy20",
			"password": "correctPassword",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/user/login", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		type ResponseBody struct {
			IsLogIn bool `json:"IsLogIn"`
		}
		resultBody := ResponseBody{}
		_ = json.Unmarshal(res.Body.Bytes(), &resultBody)
		expectedBody := ResponseBody{IsLogIn: true}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, resultBody)
	})
	t.Run("Success Case: Incorrect User Id or Password", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("CheckLogIn", service.LogInRequest{
			UserId:   "gooddy20",
			Password: "wrongPassword",
		}).Return(&service.LogInResponse{
			IsLogIn: false,
		}, nil)
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/login", hdlr.LogIn).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":  "gooddy20",
			"password": "wrongPassword",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/user/login", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		type ResponseBody struct {
			IsLogIn bool `json:"IsLogIn"`
		}
		resultBody := ResponseBody{}
		_ = json.Unmarshal(res.Body.Bytes(), &resultBody)
		expectedBody := ResponseBody{IsLogIn: false}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, resultBody)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/login", hdlr.LogIn).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":  "gooddy20",
			"password": "correctPassword",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/user/login", bytes.NewReader(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CheckLogIn")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/login", hdlr.LogIn).Methods("PUT")
		reqBody := []byte("")
		req := httptest.NewRequest("PUT", "/user/login", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CheckLogIn")
	})
	t.Run("Database Error", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("CheckLogIn", service.LogInRequest{
			UserId:   "gooddy20",
			Password: "correctPassword",
		}).Return(&service.LogInResponse{}, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/login", hdlr.LogIn).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id":  "gooddy20",
			"password": "correctPassword",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/user/login", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		type ResponseBody struct {
			IsLogIn bool `json:"IsLogIn"`
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("CreateUser", service.NewUserRequest{
			UserId:   "gooddy21",
			Password: "correctPassword",
			Username: "GoodDyZa",
			Weight:   70,
			Protein:  100,
			Fat:      50,
			Carb:     120,
		}).Return(nil)
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/", hdlr.CreateUser).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id":  "gooddy21",
			"password": "correctPassword",
			"username": "GoodDyZa",
			"weight":   70,
			"protein":  100,
			"fat":      50,
			"carb":     120,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/user/", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/", hdlr.CreateUser).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id":  "gooddy21",
			"password": "correctPassword",
			"username": "GoodDyZa",
			"weight":   70,
			"protein":  100,
			"fat":      50,
			"carb":     120,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/user/", bytes.NewReader(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateUser")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/", hdlr.CreateUser).Methods("POST")
		reqBody := []byte("")
		req := httptest.NewRequest("POST", "/user/", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateUser")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("CreateUser", service.NewUserRequest{
			UserId:   "gooddy21",
			Password: "correctPassword",
			Username: "GoodDyZa",
			Weight:   70,
			Protein:  100,
			Fat:      50,
			Carb:     120,
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/", hdlr.CreateUser).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id":  "gooddy21",
			"password": "correctPassword",
			"username": "GoodDyZa",
			"weight":   70,
			"protein":  100,
			"fat":      50,
			"carb":     120,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/user/", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestGetUserDetail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("GetUserDetail", "gooddy20").Return(&service.UserResponse{
			Username:       "GoodDy",
			Weight:         70,
			Protein:        120,
			Fat:            50,
			Carb:           120,
			FavoriteMenues: "9,10",
		}, nil)
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/{user_id}", hdlr.GetUserDetail).Methods("GET")
		req := httptest.NewRequest("GET", "/user/gooddy20", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		resultBody := service.UserResponse{}
		_ = json.Unmarshal(res.Body.Bytes(), &resultBody)
		expectedBody := service.UserResponse{
			Username:       "GoodDy",
			Weight:         70,
			Protein:        120,
			Fat:            50,
			Carb:           120,
			FavoriteMenues: "9,10",
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, resultBody)
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("GetUserDetail", "gooddy20").Return(&service.UserResponse{}, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/{user_id}", hdlr.GetUserDetail).Methods("GET")
		req := httptest.NewRequest("GET", "/user/gooddy20", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestUpdateUserDetail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("UpdateUser", service.UpdateUserRequest{
			UserId:  "gooddy20",
			Weight:  69,
			Protein: 100,
			Fat:     45,
			Carb:    100,
		}).Return(nil)
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/userdetail", hdlr.UpdateUserDetail).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id": "gooddy20",
			"weight":  69,
			"protein": 100,
			"fat":     45,
			"carb":    100,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/user/userdetail", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/userdetail", hdlr.UpdateUserDetail).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id": "gooddy20",
			"weight":  69,
			"protein": 100,
			"fat":     45,
			"carb":    100,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/user/userdetail", bytes.NewReader(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "UpdateUser")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/userdetail", hdlr.UpdateUserDetail).Methods("PUT")
		reqBody := []byte("")
		req := httptest.NewRequest("PUT", "/user/userdetail", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "UpdateUser")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewUserServiceMock()
		srv.On("UpdateUser", service.UpdateUserRequest{
			UserId:  "gooddy20",
			Weight:  69,
			Protein: 100,
			Fat:     45,
			Carb:    100,
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewUserHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/user/userdetail", hdlr.UpdateUserDetail).Methods("PUT")
		preReqBody := map[string]interface{}{
			"user_id": "gooddy20",
			"weight":  69,
			"protein": 100,
			"fat":     45,
			"carb":    100,
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/user/userdetail", bytes.NewReader(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}
