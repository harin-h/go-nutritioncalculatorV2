package service

import (
	"database/sql"
	"fmt"
	"go-nutritioncalculator2/errs"
	"go-nutritioncalculator2/logs"
	repository "go-nutritioncalculator2/repositories"
	"net/http"
	"time"
)

type recordService struct {
	recordRepo repository.RecordRepository
}

func NewRecordService(recordRepo repository.RecordRepository) recordService {
	return recordService{recordRepo: recordRepo}
}

func (s recordService) GetAllRecordsByUserId(userId string) ([]RecordResponse, error) {
	records, err := s.recordRepo.GetRecordsByUserId(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []RecordResponse{}, nil
		}
		logs.Error(err)
		return nil, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	recordsRes := []RecordResponse{}
	for i := 0; i < len(records); i++ {
		record := RecordResponse{
			Id:             records[i].Id,
			List:           records[i].List,
			Note:           records[i].Note,
			Weight:         records[i].Weight,
			Protein:        records[i].Protein,
			Fat:            records[i].Fat,
			Carb:           records[i].Carb,
			EventTimestamp: records[i].EventTimestamp,
			IsUpdated:      records[i].IsUpdated,
		}
		recordsRes = append(recordsRes, record)
	}
	return recordsRes, nil
}

func (s recordService) CreateRecord(newRecordReq NewRecordRequest) error {
	tempEventTimestamp, err := time.Parse("2006-01-02 15:04:05", newRecordReq.EventTimestamp)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	newRecord := repository.Record{
		UserId:           newRecordReq.UserId,
		List:             newRecordReq.List,
		Note:             newRecordReq.Note,
		EventTimestamp:   tempEventTimestamp,
		Status:           1,
		CreatedTimestamp: time.Now().UTC(),
	}
	_, err = s.recordRepo.CreateRecord(newRecord)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s recordService) DeleteRecord(recordId int) error {
	record, err := s.recordRepo.GetRecordById(recordId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Record Id - ", recordId, " is not found")}
		}
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	record.Status = 0
	err = s.recordRepo.UpdateRecord(*record)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s recordService) UpdateRecord(updateRecordReq UpdateRecordRequest) error {
	record, err := s.recordRepo.GetRecordById(updateRecordReq.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Record Id - ", updateRecordReq.Id, " is not found")}
		}
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	if updateRecordReq.List != "" {
		record.List = updateRecordReq.List
	}
	if updateRecordReq.Note != "" {
		record.Note = updateRecordReq.Note
	}
	if updateRecordReq.Weight != 0 {
		record.Weight = updateRecordReq.Weight
	}
	if updateRecordReq.EventTimestamp != "" {
		tempEventTimestamp, err := time.Parse("2006-01-02 15:04:05", updateRecordReq.EventTimestamp)
		if err != nil {
			logs.Error(err)
			return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
		}
		record.EventTimestamp = tempEventTimestamp
	}
	err = s.recordRepo.UpdateRecord(*record)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}
