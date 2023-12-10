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

func TestGetFavListsByUserId(t *testing.T) {
	t.Run("Success Case: Got Favorite Lists", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListsByUserId", "gooddy20").Return([]repository.FavList{
			{Id: 1, UserId: "gooddy20", Name: "Daily Breakfast", Menues: "Moo Yang-2, Sticky Rice-1 ", List: "9,9,10", Protein: 40, Fat: 10, Carb: 20, Status: 1, IsUpdated: 1, CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC()},
			{Id: 2, UserId: "gooddy20", Name: "Daily Breakfast", Menues: "Omelet-2 ", List: "1,1", Protein: 10, Fat: 2, Carb: 0, Status: 1, IsUpdated: 1, CreatedTimestamp: time.Date(2023, 13, 12, 10, 31, 15, 0, time.UTC).UTC()},
		}, nil)
		srv := service.NewFavListService(repo)
		result, _ := srv.GetFavListsByUserId("gooddy20")
		expected := []service.FavListResponse{
			{Id: 1, Name: "Daily Breakfast", Menues: "Moo Yang-2, Sticky Rice-1 ", List: "9,9,10", Protein: 40, Fat: 10, Carb: 20, IsUpdated: 1},
			{Id: 2, Name: "Daily Breakfast", Menues: "Omelet-2 ", List: "1,1", Protein: 10, Fat: 2, Carb: 0, IsUpdated: 1},
		}
		assert.Equal(t, expected, result)
	})
	t.Run("Success Case: No Favorite Lists", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListsByUserId", "gooddy20").Return([]repository.FavList{}, sql.ErrNoRows)
		srv := service.NewFavListService(repo)
		result, _ := srv.GetFavListsByUserId("gooddy20")
		expected := []service.FavListResponse{}
		assert.Equal(t, expected, result)
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListsByUserId", "gooddy20").Return([]repository.FavList{}, sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		_, err := srv.GetFavListsByUserId("gooddy20")
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestCreateFavList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("CreateFavList", repository.FavList{
			UserId:           "gooddy20",
			Name:             "Daily Breakfast V2",
			List:             "1,1,3",
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.FavList{
			Id:               3,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast V2",
			Menues:           "Omelet-2, Boiled Egg-1 ",
			List:             "1,1,3",
			Protein:          14,
			Fat:              2,
			Carb:             0,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}, nil)
		srv := service.NewFavListService(repo)
		err := srv.CreateFavList(service.NewFavListRequest{UserId: "gooddy20", Name: "Daily Breakfast V2", List: "1,1,3"})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("CreateFavList", repository.FavList{
			UserId:           "gooddy20",
			Name:             "Daily Breakfast V2",
			List:             "1,1,3",
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.FavList{}, sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		err := srv.CreateFavList(service.NewFavListRequest{UserId: "gooddy20", Name: "Daily Breakfast V2", List: "1,1,3"})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestDeleteFavList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           0,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewFavListService(repo)
		err := srv.DeleteFavList(1)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Get Favorite List Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		err := srv.DeleteFavList(1)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Update Favorite List Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           0,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		err := srv.DeleteFavList(1)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestUpdateFavList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast V2",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewFavListService(repo)
		err := srv.UpdateFavList(service.UpdateFavListRequest{
			Id:   1,
			Name: "Daily Breakfast V2",
			List: "9,9,9,10",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("No The Favorite List Id", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, sql.ErrNoRows)
		srv := service.NewFavListService(repo)
		err := srv.UpdateFavList(service.UpdateFavListRequest{
			Id:   1,
			Name: "Daily Breakfast V2",
			List: "9,9,9,10",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Favorite List Id - ", 1, "is not found")})
		repo.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Get Favorite List Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		err := srv.UpdateFavList(service.UpdateFavListRequest{
			Id:   1,
			Name: "Daily Breakfast V2",
			List: "9,9,9,10",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Update Favorite List Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast V2",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        1,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		err := srv.UpdateFavList(service.UpdateFavListRequest{
			Id:   1,
			Name: "Daily Breakfast V2",
			List: "9,9,9,10",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestRecoverFavList(t *testing.T) {
	t.Run("Success Case 1", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(1, 10, 0)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case 2", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "10,11,11",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(1, 9, 11)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case 3", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 2).Return(&repository.FavList{
			Id:               2,
			UserId:           "gooddy20",
			Name:             "Daily Lunch",
			Menues:           "Moo Yang-2, Sticky Rice-1 ,Orange Juice-1 ",
			List:             "9,9,10,11",
			Protein:          40,
			Fat:              10,
			Carb:             40,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 15, 12, 10, 23, 38, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               2,
			UserId:           "gooddy20",
			Name:             "Daily Lunch",
			Menues:           "Moo Yang-2, Sticky Rice-1 ,Orange Juice-1 ",
			List:             "9,9,11",
			Protein:          40,
			Fat:              10,
			Carb:             40,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 15, 12, 10, 23, 38, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(2, 10, 0)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case: No Deleted Menu in Favorite Lists", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 2).Return(&repository.FavList{
			Id:               2,
			UserId:           "gooddy20",
			Name:             "Daily Lunch",
			Menues:           "Moo Yang-2, Sticky Rice-1 ,Orange Juice-1 ",
			List:             "9,9,10,11",
			Protein:          40,
			Fat:              10,
			Carb:             40,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 15, 12, 10, 23, 38, 0, time.UTC).UTC(),
		}, nil)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(2, 12, 0)
		assert.ErrorIs(t, err, nil)
		repo.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("No The Favorite List Id", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 2).Return(&repository.FavList{}, sql.ErrNoRows)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(2, 12, 0)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Favorite List Id - ", 2, "is not found")})
		repo.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Get Favorite List Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 2).Return(&repository.FavList{}, sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(2, 12, 0)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Parse List (String to Int) Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 2).Return(&repository.FavList{
			Id:               2,
			UserId:           "gooddy20",
			Name:             "Daily Lunch",
			Menues:           "Moo Yang-2, Sticky Rice-1 ,Orange Juice-1 ",
			List:             "9,9,10,11.5",
			Protein:          40,
			Fat:              10,
			Carb:             40,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 15, 12, 10, 23, 38, 0, time.UTC).UTC(),
		}, nil)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(2, 12, 0)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateFavList")
	})
	t.Run("Update Favorite List Database Error", func(t *testing.T) {
		repo := repository.NewFavListRepositoryMock()
		repo.On("GetFavListById", 1).Return(&repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "9,9,10",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateFavList", repository.FavList{
			Id:               1,
			UserId:           "gooddy20",
			Name:             "Daily Breakfast",
			Menues:           "Moo Yang-2, Sticky Rice-1 ",
			List:             "10,11,11",
			Protein:          40,
			Fat:              10,
			Carb:             20,
			Status:           1,
			IsUpdated:        0,
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(sql.ErrConnDone)
		srv := service.NewFavListService(repo)
		err := srv.RecoverFavList(1, 9, 11)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}
