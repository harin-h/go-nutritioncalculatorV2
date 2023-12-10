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
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateRecord(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("CreateRecord", service.NewRecordRequest{
			UserId:         "gooddy20",
			List:           "9,9,10",
			Note:           "Breakfast",
			Weight:         70,
			EventTimestamp: "2023-12-05 10:00:00",
		}).Return(nil)
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.CreateRecord).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"list":            "9,9,10",
			"note":            "Breakfast",
			"weight":          70,
			"event_timestamp": "2023-12-05 10:00:00",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/record/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.CreateRecord).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"list":            "9,9,10",
			"note":            "Breakfast",
			"weight":          70,
			"event_timestamp": "2023-12-05 10:00:00",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/record/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateRecord")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.CreateRecord).Methods("POST")
		reqBody := []byte("")
		req := httptest.NewRequest("POST", "/record/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "CreateRecord")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("CreateRecord", service.NewRecordRequest{
			UserId:         "gooddy20",
			List:           "9,9,10",
			Note:           "Breakfast",
			Weight:         70,
			EventTimestamp: "2023-12-05 10:00:00",
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.CreateRecord).Methods("POST")
		preReqBody := map[string]interface{}{
			"user_id":         "gooddy20",
			"list":            "9,9,10",
			"note":            "Breakfast",
			"weight":          70,
			"event_timestamp": "2023-12-05 10:00:00",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("POST", "/record/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestDeleteRecord(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("DeleteRecord", 1).Return(nil)
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/{record_id}", hdlr.DeleteRecord).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/record/1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Parse Id (String to Int) Error", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/{record_id}", hdlr.DeleteRecord).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/record/1.1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Parse data type error", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "DeleteRecord")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("DeleteRecord", 1).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/{record_id}", hdlr.DeleteRecord).Methods("DELETE")
		req := httptest.NewRequest("DELETE", "/record/1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestUpdateRecord(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("UpdateRecord", service.UpdateRecordRequest{
			Id:             1,
			List:           "9,9,10,11",
			Note:           "Breakfast + Juice",
			Weight:         0,
			EventTimestamp: "2023-12-05 10:20:00",
		}).Return(nil)
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.UpdateRecord).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":              1,
			"list":            "9,9,10,11",
			"note":            "Breakfast + Juice",
			"weight":          0,
			"event_timestamp": "2023-12-05 10:20:00",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/record/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Incorrect Request Header", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.UpdateRecord).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":              1,
			"list":            "9,9,10,11",
			"note":            "Breakfast + Juice",
			"weight":          0,
			"event_timestamp": "2023-12-05 10:20:00",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/record/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Header", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "UpdateRecord")
	})
	t.Run("Incorrect Request Body", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.UpdateRecord).Methods("PUT")
		reqBody := []byte("")
		req := httptest.NewRequest("PUT", "/record/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotAcceptable, res.Code)
		assert.Equal(t, "Incorrect Request Body", strings.Replace(res.Body.String(), "\n", "", -1))
		srv.AssertNotCalled(t, "UpdateRecord")
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("UpdateRecord", service.UpdateRecordRequest{
			Id:             1,
			List:           "9,9,10,11",
			Note:           "Breakfast + Juice",
			Weight:         0,
			EventTimestamp: "2023-12-05 10:20:00",
		}).Return(errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/", hdlr.UpdateRecord).Methods("PUT")
		preReqBody := map[string]interface{}{
			"id":              1,
			"list":            "9,9,10,11",
			"note":            "Breakfast + Juice",
			"weight":          0,
			"event_timestamp": "2023-12-05 10:20:00",
		}
		reqBody, _ := json.Marshal(preReqBody)
		req := httptest.NewRequest("PUT", "/record/", bytes.NewBuffer(reqBody))
		req.Header.Add("content-type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}

func TestGetRecordsByUserId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("GetAllRecordsByUserId", "gooddy20").Return([]service.RecordResponse{
			{Id: 1,
				List:           "9,9,10",
				Menues:         "Moo Yang-2 ,Sticky Rice-1 ",
				Note:           "Breakfast",
				Weight:         70,
				Protein:        40,
				Fat:            10,
				Carb:           20,
				EventTimestamp: time.Date(2023, 12, 5, 10, 0, 0, 0, time.UTC).UTC(),
				IsUpdated:      1},
			{Id: 2,
				List:           "1,9",
				Menues:         "Ramyeon-1 ,Moo Yang-1 ",
				Note:           "Lunch",
				Weight:         70,
				Protein:        25,
				Fat:            25,
				Carb:           65,
				EventTimestamp: time.Date(2023, 12, 5, 12, 30, 0, 0, time.UTC).UTC(),
				IsUpdated:      1},
		}, nil)
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/{user_id}", hdlr.GetRecordsByUserId).Methods("GET")
		req := httptest.NewRequest("GET", "/record/gooddy20", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		resultBody := []service.RecordResponse{}
		_ = json.Unmarshal(res.Body.Bytes(), &resultBody)
		expectedBody := []service.RecordResponse{
			{Id: 1,
				List:           "9,9,10",
				Menues:         "Moo Yang-2 ,Sticky Rice-1 ",
				Note:           "Breakfast",
				Weight:         70,
				Protein:        40,
				Fat:            10,
				Carb:           20,
				EventTimestamp: time.Date(2023, 12, 5, 10, 0, 0, 0, time.UTC).UTC(),
				IsUpdated:      1},
			{Id: 2,
				List:           "1,9",
				Menues:         "Ramyeon-1 ,Moo Yang-1 ",
				Note:           "Lunch",
				Weight:         70,
				Protein:        25,
				Fat:            25,
				Carb:           65,
				EventTimestamp: time.Date(2023, 12, 5, 12, 30, 0, 0, time.UTC).UTC(),
				IsUpdated:      1},
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expectedBody, resultBody)
	})
	t.Run("Service Error", func(t *testing.T) {
		srv := service.NewRecordServiceMock()
		srv.On("GetAllRecordsByUserId", "gooddy20").Return([]service.RecordResponse{}, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		hdlr := handler.NewRecordHandler(srv)
		r := mux.NewRouter()
		r.HandleFunc("/record/{user_id}", hdlr.GetRecordsByUserId).Methods("GET")
		req := httptest.NewRequest("GET", "/record/gooddy20", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "Unexpected error", strings.Replace(res.Body.String(), "\n", "", -1))
	})
}
