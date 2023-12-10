package service_test

import (
	"database/sql"
	"fmt"
	"go-nutritioncalculator2/errs"
	repository "go-nutritioncalculator2/repositories"
	service "go-nutritioncalculator2/services"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAllRecordsByUserId(t *testing.T) {
	t.Run("Success Case 1", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordsByUserId", "gooddy20").Return([]repository.Record{
			{Id: 1,
				UserId:           "gooddy20",
				List:             "9,9,10",
				Menues:           "Moo Yang-2, Sticky Rice-1 ",
				Note:             "Breakfast",
				Weight:           70,
				Protein:          40,
				Fat:              10,
				Carb:             20,
				EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
				Status:           1,
				IsUpdated:        1,
				CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
			},
			{Id: 2,
				UserId:           "gooddy20",
				List:             "14",
				Menues:           "Momo Buffet-1",
				Note:             "My BD",
				Weight:           71,
				Protein:          50,
				Fat:              35,
				Carb:             40,
				EventTimestamp:   time.Date(2023, 12, 5, 12, 30, 56, 0, time.UTC).UTC(),
				Status:           1,
				IsUpdated:        1,
				CreatedTimestamp: time.Date(2023, 12, 5, 19, 0, 2, 0, time.UTC).UTC(),
			},
		}, nil)
		srv := service.NewRecordService(repo)
		result, _ := srv.GetAllRecordsByUserId("gooddy20")
		expected := []service.RecordResponse{
			{Id: 1,
				List:           "9,9,10",
				Menues:         "Moo Yang-2, Sticky Rice-1 ",
				Note:           "Breakfast",
				Weight:         70,
				Protein:        40,
				Fat:            10,
				Carb:           20,
				EventTimestamp: time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
				IsUpdated:      1,
			},
			{Id: 2,
				List:           "14",
				Menues:         "Momo Buffet-1",
				Note:           "My BD",
				Weight:         71,
				Protein:        50,
				Fat:            35,
				Carb:           40,
				EventTimestamp: time.Date(2023, 12, 5, 12, 30, 56, 0, time.UTC).UTC(),
				IsUpdated:      1,
			},
		}
		assert.Equal(t, expected, result)
	})
	t.Run("Success Case 2", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordsByUserId", "gooddy20").Return([]repository.Record{}, sql.ErrNoRows)
		srv := service.NewRecordService(repo)
		result, _ := srv.GetAllRecordsByUserId("gooddy20")
		expected := []service.RecordResponse{}
		assert.Equal(t, expected, result)
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordsByUserId", "gooddy20").Return([]repository.Record{}, sql.ErrConnDone)
		srv := service.NewRecordService(repo)
		_, err := srv.GetAllRecordsByUserId("gooddy20")
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestCreateRecord(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("CreateRecord", repository.Record{
			UserId:           "gooddy20",
			List:             "9,9,10,11",
			Note:             "Lunch",
			Weight:           70,
			EventTimestamp:   time.Date(2023, 12, 5, 12, 30, 56, 0, time.UTC).UTC(),
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.Record{
			UserId:           "gooddy20",
			List:             "9,9,10,11",
			Note:             "Lunch",
			Weight:           70,
			EventTimestamp:   time.Date(2023, 12, 5, 12, 30, 56, 0, time.UTC).UTC(),
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}, nil)
		srv := service.NewRecordService(repo)
		err := srv.CreateRecord(service.NewRecordRequest{
			UserId:         "gooddy20",
			List:           "9,9,10,11",
			Note:           "Lunch",
			Weight:         70,
			EventTimestamp: "2023-12-05 12:30:56",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Parse Event Timestamp (String to Datetime) Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		srv := service.NewRecordService(repo)
		err := srv.CreateRecord(service.NewRecordRequest{
			UserId:         "gooddy20",
			List:           "9,9,10,11",
			Note:           "Lunch",
			Weight:         70,
			EventTimestamp: "2023-12-05 12:3x;56",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "CreateRecord")
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("CreateRecord", repository.Record{
			UserId:           "gooddy20",
			List:             "9,9,10,11",
			Note:             "Lunch",
			Weight:           70,
			EventTimestamp:   time.Date(2023, 12, 5, 12, 30, 56, 0, time.UTC).UTC(),
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.Record{}, sql.ErrConnDone)
		srv := service.NewRecordService(repo)
		err := srv.CreateRecord(service.NewRecordRequest{
			UserId:         "gooddy20",
			List:           "9,9,10,11",
			Note:           "Lunch",
			Weight:         70,
			EventTimestamp: "2023-12-05 12:30:56",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestDeleteRecord(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Breakfast",
			Weight:           70,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateRecord", repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Breakfast",
			Weight:           70,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
			Status:           0,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewRecordService(repo)
		err := srv.DeleteRecord(1)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("No The Record Id", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{}, sql.ErrNoRows)
		srv := service.NewRecordService(repo)
		err := srv.DeleteRecord(1)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Record Id - ", 1, " is not found")})
		repo.AssertNotCalled(t, "UpdateRecord")
	})
	t.Run("Get Record Database Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{}, sql.ErrConnDone)
		srv := service.NewRecordService(repo)
		err := srv.DeleteRecord(1)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateRecord")
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Breakfast",
			Weight:           70,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateRecord", repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Breakfast",
			Weight:           70,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
			Status:           0,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}).Return(sql.ErrConnDone)
		srv := service.NewRecordService(repo)
		err := srv.DeleteRecord(1)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestUpdateRecord(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Breakfast",
			Weight:           70,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateRecord", repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Extra Lunch",
			Weight:           74,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 5, 12, 00, 00, 0, time.UTC).UTC(),
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewRecordService(repo)
		err := srv.UpdateRecord(service.UpdateRecordRequest{
			Id:             1,
			List:           "9,9,9,10",
			Note:           "Extra Lunch",
			Weight:         74,
			EventTimestamp: "2023-12-05 12:00:00",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("No The Record Id", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{}, sql.ErrNoRows)
		srv := service.NewRecordService(repo)
		err := srv.UpdateRecord(service.UpdateRecordRequest{
			Id:             1,
			List:           "9,9,9,10",
			Note:           "Extra Lunch",
			Weight:         74,
			EventTimestamp: "2023-12-05 12:00:00",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Record Id - ", 1, " is not found")})
		repo.AssertNotCalled(t, "UpdateRecord")
	})
	t.Run("Get Record Database Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{}, sql.ErrConnDone)
		srv := service.NewRecordService(repo)
		err := srv.UpdateRecord(service.UpdateRecordRequest{
			Id:             1,
			List:           "9,9,9,10",
			Note:           "Extra Lunch",
			Weight:         74,
			EventTimestamp: "2023-12-05 12:00:00",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateRecord")
	})
	t.Run("Parse (String to Datetime) Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Breakfast",
			Weight:           70,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}, nil)
		srv := service.NewRecordService(repo)
		err := srv.UpdateRecord(service.UpdateRecordRequest{
			Id:             1,
			List:           "9,9,9,10",
			Note:           "Extra Lunch",
			Weight:         74,
			EventTimestamp: "2023-12-05 12;0x:0x",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateRecord")
	})
	t.Run("Update Record Database Error", func(t *testing.T) {
		repo := repository.NewRecordRepositoryMock()
		repo.On("GetRecordById", 1).Return(&repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Breakfast",
			Weight:           70,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 4, 10, 30, 12, 0, time.UTC).UTC(),
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateRecord", repository.Record{
			Id:               1,
			UserId:           "gooddy20",
			List:             "9,9,9,10",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			Note:             "Extra Lunch",
			Weight:           74,
			Protein:          40,
			Fat:              10,
			Carb:             20,
			EventTimestamp:   time.Date(2023, 12, 5, 12, 00, 00, 0, time.UTC).UTC(),
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 12, 4, 19, 30, 19, 0, time.UTC).UTC(),
		}).Return(sql.ErrConnDone)
		srv := service.NewRecordService(repo)
		err := srv.UpdateRecord(service.UpdateRecordRequest{
			Id:             1,
			List:           "9,9,9,10",
			Note:           "Extra Lunch",
			Weight:         74,
			EventTimestamp: "2023-12-05 12:00:00",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}
