package service_test

import (
	"database/sql"
	"go-nutritioncalculator2/errs"
	repository "go-nutritioncalculator2/repositories"
	service "go-nutritioncalculator2/services"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateMenu(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("CreateMenu", repository.Menu{
			Name:             "Omelet",
			Protein:          5,
			Fat:              1,
			Carb:             0,
			CreatorId:        "gooddy20",
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.Menu{
			Id:               1,
			Name:             "Omelet",
			Protein:          5,
			Fat:              1,
			Carb:             0,
			CreatorId:        "gooddy20",
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}, nil)
		srv := service.NewMenuService(repo)
		err := srv.CreateMenu(service.NewMenuRequest{
			Name:      "Omelet",
			Protein:   5,
			Fat:       1,
			Carb:      0,
			CreatorId: "gooddy20",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("CreateMenu", repository.Menu{
			Name:             "Omelet",
			Protein:          5,
			Fat:              1,
			Carb:             0,
			CreatorId:        "gooddy20",
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.Menu{}, sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		err := srv.CreateMenu(service.NewMenuRequest{
			Name:      "Omelet",
			Protein:   5,
			Fat:       1,
			Carb:      0,
			CreatorId: "gooddy20",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestGetAllMenues(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("GetAllMenues").Return([]repository.Menu{
			{Id: 1, Name: "Omelet", Protein: 5, Fat: 1, Carb: 0, CreatorId: "gooddy20", CreatorName: "GoodDy", Like: 2, Status: 1, CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC()},
			{Id: 2, Name: "Fried Egg", Protein: 5, Fat: 2, Carb: 0, CreatorId: "gooddy20", CreatorName: "GoodDy", Like: 0, Status: 0, CreatedTimestamp: time.Date(2023, 11, 14, 15, 12, 35, 0, time.UTC).UTC()},
			{Id: 3, Name: "Khai Tom", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1, CreatedTimestamp: time.Date(2023, 11, 14, 18, 06, 11, 0, time.UTC).UTC()},
		}, nil)
		srv := service.NewMenuService(repo)
		result, _ := srv.GetAllMenues()
		expected := []service.MenuResponse{
			{Id: 1, Name: "Omelet", Protein: 5, Fat: 1, Carb: 0, CreatorId: "gooddy20", CreatorName: "GoodDy", Like: 2, Status: 1},
			{Id: 2, Name: "Fried Egg", Protein: 5, Fat: 2, Carb: 0, CreatorId: "gooddy20", CreatorName: "GoodDy", Like: 0, Status: 0},
			{Id: 3, Name: "Khai Tom", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1},
		}
		assert.Equal(t, expected, result)
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("GetAllMenues").Return([]repository.Menu{}, sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		_, err := srv.GetAllMenues()
		assert.Equal(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestUpdateMenu(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("UpdateMenu", repository.Menu{Id: 1}).Return(nil)
		repo.On("GetMenuById", 1).Return(&repository.Menu{Id: 1, Name: "Omelet", Protein: 5, Fat: 1, Carb: 0, CreatorId: "gooddy20", CreatorName: "GoodDy", Like: 2, Status: 0, CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC()}, nil)
		repo.On("CreateMenu", repository.Menu{
			Id:               0,
			Name:             "Omelet",
			Protein:          5.5,
			Fat:              0.5,
			Carb:             1,
			CreatorId:        "gooddy20",
			CreatorName:      "GoodDy",
			Like:             2,
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.Menu{
			Id:               4,
			Name:             "Omelet",
			Protein:          5.5,
			Fat:              0.5,
			Carb:             1,
			CreatorId:        "gooddy20",
			CreatorName:      "GoodDy",
			Like:             2,
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}, nil)
		srv := service.NewMenuService(repo)
		err := srv.UpdateMenu(service.UpdateMenuRequest{Id: 1, Name: "Omelet", Protein: 5.5, Fat: 0.5, Carb: 1})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Update Menu Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("UpdateMenu", repository.Menu{Id: 1}).Return(sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		err := srv.UpdateMenu(service.UpdateMenuRequest{Id: 1, Name: "Omelet", Protein: 5.5, Fat: 0.5, Carb: 1})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "GetMenuById")
		repo.AssertNotCalled(t, "CreateMenu")
	})
	t.Run("Get Menu Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("UpdateMenu", repository.Menu{Id: 1}).Return(nil)
		repo.On("GetMenuById", 1).Return(&repository.Menu{Id: 1, Name: "Omelet", Protein: 5, Fat: 1, Carb: 0, CreatorId: "gooddy20", CreatorName: "GoodDy", Like: 2, Status: 0, CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC()},
			sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		err := srv.UpdateMenu(service.UpdateMenuRequest{Id: 1, Name: "Omelet", Protein: 5.5, Fat: 0.5, Carb: 1})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "CreateMenu")
	})
	t.Run("Create Menu Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("UpdateMenu", repository.Menu{Id: 1}).Return(nil)
		repo.On("GetMenuById", 1).Return(&repository.Menu{Id: 1, Name: "Omelet", Protein: 5, Fat: 1, Carb: 0, CreatorId: "gooddy20", CreatorName: "GoodDy", Like: 2, Status: 0, CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC()}, nil)
		repo.On("CreateMenu", repository.Menu{
			Id:               0,
			Name:             "Omelet",
			Protein:          5.5,
			Fat:              0.5,
			Carb:             1,
			CreatorId:        "gooddy20",
			CreatorName:      "GoodDy",
			Like:             2,
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(&repository.Menu{
			Id:               4,
			Name:             "Omelet",
			Protein:          5.5,
			Fat:              0.5,
			Carb:             1,
			CreatorId:        "gooddy20",
			CreatorName:      "GoodDy",
			Like:             2,
			Status:           1,
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}, sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		err := srv.UpdateMenu(service.UpdateMenuRequest{Id: 1, Name: "Omelet", Protein: 5.5, Fat: 0.5, Carb: 1})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestRecoverMenu(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("GetMenuById", 3).Return(&repository.Menu{Id: 3, Name: "Khai Tom", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1, CreatedTimestamp: time.Date(2023, 11, 14, 18, 06, 11, 0, time.UTC).UTC()}, nil)
		repo.On("CreateMenu", repository.Menu{Id: 0, Name: "Boiled Egg", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1, CreatedTimestamp: time.Now().UTC().Truncate(time.Second)}).Return(&repository.Menu{Id: 4, Name: "Boiled Egg", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1, CreatedTimestamp: time.Now().UTC().Truncate(time.Second)}, nil)
		srv := service.NewMenuService(repo)
		result, _ := srv.RecoverMenu(3, "Boiled Egg")
		expected := &service.MenuResponse{Id: 4, Name: "Boiled Egg", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1}
		assert.Equal(t, expected, result)
	})
	t.Run("No The Menu Id", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("GetMenuById", 3).Return(&repository.Menu{}, sql.ErrNoRows)
		srv := service.NewMenuService(repo)
		_, err := srv.RecoverMenu(3, "Boiled Egg")
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "Menu Id is not found"})
	})
	t.Run("Get Menu Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("GetMenuById", 3).Return(&repository.Menu{}, sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		_, err := srv.RecoverMenu(3, "Boiled Egg")
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
	t.Run("Create Menu Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("GetMenuById", 3).Return(&repository.Menu{Id: 3, Name: "Khai Tom", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1, CreatedTimestamp: time.Date(2023, 11, 14, 18, 06, 11, 0, time.UTC).UTC()}, nil)
		repo.On("CreateMenu", repository.Menu{Id: 0, Name: "Boiled Egg", Protein: 4, Fat: 0, Carb: 0, CreatorId: "kornkoko", CreatorName: "Kornkoko", Like: 1, Status: 1, CreatedTimestamp: time.Now().UTC().Truncate(time.Second)}).Return(&repository.Menu{}, sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		_, err := srv.RecoverMenu(3, "Boiled Egg")
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestDeleteMenu(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("UpdateMenu", repository.Menu{Id: 1}).Return(nil)
		srv := service.NewMenuService(repo)
		err := srv.DeleteMenu(1)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewMenuRepositoryMock()
		repo.On("UpdateMenu", repository.Menu{Id: 1}).Return(sql.ErrConnDone)
		srv := service.NewMenuService(repo)
		err := srv.DeleteMenu(1)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}
