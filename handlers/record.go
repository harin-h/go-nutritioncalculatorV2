package handler

import (
	"encoding/json"
	"go-nutritioncalculator2/errs"
	service "go-nutritioncalculator2/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type recordHandler struct {
	recordSrv service.RecordService
}

func NewRecordHandler(recordSrv service.RecordService) recordHandler {
	return recordHandler{recordSrv: recordSrv}
}

// CreateRecord ... Create a "Record"
// @Summary Create a "Record"
// @Description Create a 'Record'
// @Tags Record
// @Accept json
// @Param request body service.NewRecordRequest true "`Record`'s data detail"
// @Response 200
// @Response 406 "Request Body Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /record/ [post]
func (h recordHandler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request service.NewRecordRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	err = h.recordSrv.CreateRecord(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}

// DeleteRecord ... Delete a "Record"
// @Summary Delete a "Record"
// @Description Delete a 'Record'
// @Tags Record
// @Param record_id path int true "`Record`'s id that you want to delete"
// @Response 200
// @Response 406 "Request parameters Not Acceptable"
// @Response 500 "Internal Server Error"
// @Router /record/{record_id} [delete]
func (h recordHandler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recordId, err := strconv.ParseInt(vars["record_id"], 0, 0)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Parse data type error"})
		return
	}
	err = h.recordSrv.DeleteRecord(int(recordId))
	if err != nil {
		handlerError(w, err)
		return
	}
}

// UpdateRecord ... Update a "Record"
// @Summary Update a "Record"
// @Description Update a 'Record'
// @Tags Record
// @Accept json
// @Param request body service.UpdateRecordRequest true "`Record`'s data detail that you want to change to"
// @Response 200
// @Response 406 "Request Body Not Acceptable or `Record`'s id is not found"
// @Response 500 "Internal Server Error"
// @Router /record/ [put]
func (h recordHandler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Header"})
		return
	}
	var request service.UpdateRecordRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerError(w, errs.AppError{Code: http.StatusNotAcceptable, Message: "Incorrect Request Body"})
		return
	}
	err = h.recordSrv.UpdateRecord(request)
	if err != nil {
		handlerError(w, err)
		return
	}
}

// GetRecordsByUserId ... Get all "Record" of "User"
// @Summary Get all "Record" of "User"
// @Description Get all `Record` of `User` by `User Id`
// @Tags Record
// @Produce json
// @Param user_id path string true "`User Id` that you want to get `Record`"
// @Response 200 {object} []service.RecordResponse
// @Response 500 "Internal Server Error"
// @Router /record/{user_id} [get]
func (h recordHandler) GetRecordsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := h.recordSrv.GetAllRecordsByUserId(vars["user_id"])
	if err != nil {
		handlerError(w, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
