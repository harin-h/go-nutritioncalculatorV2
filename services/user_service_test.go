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

func TestCheckLogIn(t *testing.T) {
	type testCase struct {
		Name     string
		Request  service.LogInRequest
		Expected *service.LogInResponse
	}
	cases := []testCase{
		{Name: "Success Case: Correct Password", Request: service.LogInRequest{UserId: "gooddy20", Password: "correctPassword"}, Expected: &service.LogInResponse{IsLogIn: true}},
		{Name: "Success Case: Incorrect Password", Request: service.LogInRequest{UserId: "gooddy20", Password: "whatPassword"}, Expected: &service.LogInResponse{IsLogIn: false}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			repo := repository.NewUserRepositoryMock()
			repo.On("GetUserById", c.Request.UserId).Return(&repository.User{UserId: c.Request.UserId, Password: "correctPassword"}, nil)
			srv := service.NewUserService(repo)
			result, _ := srv.CheckLogIn(c.Request)
			assert.Equal(t, c.Expected, result)
		})
	}
	t.Run("Success Case: No The User Id", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy19").Return(&repository.User{}, sql.ErrNoRows)
		srv := service.NewUserService(repo)
		result, _ := srv.CheckLogIn(service.LogInRequest{UserId: "gooddy19", Password: "correctPassword"})
		assert.Equal(t, &service.LogInResponse{IsLogIn: false}, result)
	})

	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, sql.ErrConnDone)
		srv := service.NewUserService(repo)
		_, err := srv.CheckLogIn(service.LogInRequest{UserId: "gooddy20", Password: "correctPassword"})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestGetUserDetail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(
			&repository.User{
				UserId:           "gooddy20",
				Password:         "correctPassword",
				Username:         "GoodDy",
				Weight:           71,
				Protein:          120,
				Fat:              60,
				Carb:             120,
				FavoriteMenues:   "11,12",
				CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
			}, nil)
		srv := service.NewUserService(repo)
		result, _ := srv.GetUserDetail("gooddy20")
		expected := &service.UserResponse{
			Username:       "GoodDy",
			Weight:         71,
			Protein:        120,
			Fat:            60,
			Carb:           120,
			FavoriteMenues: "11,12",
		}
		assert.Equal(t, expected, result)
	})
	t.Run("No The User Id", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy19").Return(&repository.User{}, sql.ErrNoRows)
		srv := service.NewUserService(repo)
		_, err := srv.GetUserDetail("gooddy19")
		assert.Equal(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is not found"})
	})
	t.Run("Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, sql.ErrConnDone)
		srv := service.NewUserService(repo)
		_, err := srv.GetUserDetail("gooddy20")
		assert.Equal(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestCreateUser(t *testing.T) {
	type TestCase struct {
		Name     string
		Request  service.NewUserRequest
		Expected error
	}
	invalid_parameter_cases := []TestCase{
		{Name: "Invalid User Id", Request: service.NewUserRequest{UserId: "good", Password: "correctPassword", Username: "GoodDyZa", Weight: 68, Protein: 0, Fat: 0, Carb: 0}, Expected: errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id need to contain more than 5 letter and alphabet only"}},
		{Name: "Invalid Password", Request: service.NewUserRequest{UserId: "gooddy21", Password: "p a s s", Username: "GoodDyZa", Weight: 68, Protein: 0, Fat: 0, Carb: 0}, Expected: errs.AppError{Code: http.StatusNotAcceptable, Message: "Password need to contain more than 5 letter and no whitespace"}},
		{Name: "Invalid Username", Request: service.NewUserRequest{UserId: "gooddy21", Password: "correctPassword", Username: "Good", Weight: 68, Protein: 0, Fat: 0, Carb: 0}, Expected: errs.AppError{Code: http.StatusNotAcceptable, Message: "Username need to contain more than 5 letter and alphabet only"}},
	}
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, nil)
		repo.On("GetUserById", "gooddy21").Return(&repository.User{}, sql.ErrNoRows)
		repo.On("GetUserByUsername", "GoodDy").Return(&repository.User{}, nil)
		repo.On("GetUserByUsername", "GoodDyZa").Return(&repository.User{}, sql.ErrNoRows)
		repo.On("CreateUser", repository.User{
			UserId:           "gooddy21",
			Password:         "correctPassword",
			Username:         "GoodDyZa",
			Weight:           68,
			Protein:          0,
			Fat:              0,
			Carb:             0,
			FavoriteMenues:   "",
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.CreateUser(service.NewUserRequest{UserId: "gooddy21", Password: "correctPassword", Username: "GoodDyZa", Weight: 68, Protein: 0, Fat: 0, Carb: 0})
		assert.ErrorIs(t, err, nil)
	})
	for _, c := range invalid_parameter_cases {
		t.Run(c.Name, func(t *testing.T) {
			repo := repository.NewUserRepositoryMock()
			srv := service.NewUserService(repo)
			err := srv.CreateUser(c.Request)
			assert.ErrorIs(t, err, c.Expected)
			repo.AssertNotCalled(t, "GetUserById")
			repo.AssertNotCalled(t, "GetUserByUsername")
			repo.AssertNotCalled(t, "CreateUser")
		})
	}
	t.Run("User Id is already used", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, nil)
		repo.On("GetUserById", "gooddy21").Return(&repository.User{}, sql.ErrNoRows)
		srv := service.NewUserService(repo)
		err := srv.CreateUser(service.NewUserRequest{UserId: "gooddy20", Password: "correctPassword", Username: "GoodDyZa", Weight: 68, Protein: 0, Fat: 0, Carb: 0})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is already used"})
		repo.AssertNotCalled(t, "GetUserByUsername")
		repo.AssertNotCalled(t, "CreateUser")
	})
	t.Run("Get User Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy21").Return(&repository.User{}, sql.ErrConnDone)
		srv := service.NewUserService(repo)
		err := srv.CreateUser(service.NewUserRequest{UserId: "gooddy21", Password: "correctPassword", Username: "GoodDyZa", Weight: 68, Protein: 0, Fat: 0, Carb: 0})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "GetUserByUsername")
		repo.AssertNotCalled(t, "CreateUser")
	})
	t.Run("Username is already used", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, nil)
		repo.On("GetUserById", "gooddy21").Return(&repository.User{}, sql.ErrNoRows)
		repo.On("GetUserByUsername", "GoodDy").Return(&repository.User{}, nil)
		repo.On("GetUserByUsername", "GoodDyZa").Return(&repository.User{}, sql.ErrNoRows)
		srv := service.NewUserService(repo)
		err := srv.CreateUser(service.NewUserRequest{UserId: "gooddy21", Password: "correctPassword", Username: "GoodDy", Weight: 68, Protein: 0, Fat: 0, Carb: 0})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "Username is already used"})
		repo.AssertNotCalled(t, "CreateUser")
	})
	t.Run("Get Username Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy21").Return(&repository.User{}, sql.ErrNoRows)
		repo.On("GetUserByUsername", "GoodDyZa").Return(&repository.User{}, sql.ErrConnDone)
		srv := service.NewUserService(repo)
		err := srv.CreateUser(service.NewUserRequest{UserId: "gooddy21", Password: "correctPassword", Username: "GoodDyZa", Weight: 68, Protein: 0, Fat: 0, Carb: 0})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "CreateUser")
	})
	t.Run("Create User Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy21").Return(&repository.User{}, sql.ErrNoRows)
		repo.On("GetUserByUsername", "GoodDyZa").Return(&repository.User{}, sql.ErrNoRows)
		repo.On("CreateUser", repository.User{
			UserId:           "gooddy21",
			Password:         "correctPassword",
			Username:         "GoodDyZa",
			Weight:           68,
			Protein:          0,
			Fat:              0,
			Carb:             0,
			FavoriteMenues:   "",
			CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
		}).Return(sql.ErrConnDone)
		srv := service.NewUserService(repo)
		err := srv.CreateUser(service.NewUserRequest{UserId: "gooddy21", Password: "correctPassword", Username: "GoodDyZa", Weight: 68, Protein: 0, Fat: 0, Carb: 0})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Success Case: Update New Password", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:       "correctPasswordV2",
			Username:       "GoodDy",
			Weight:         71,
			Protein:        120,
			Fat:            60,
			Carb:           120,
			FavoriteMenues: "11,12",
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:   "gooddy20",
			Password: "correctPasswordV2",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case: Update New Username", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:       "correctPassword",
			Username:       "GoodDyInwZa20",
			Weight:         71,
			Protein:        120,
			Fat:            60,
			Carb:           120,
			FavoriteMenues: "11,12",
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:   "gooddy20",
			Username: "GoodDyInwZa20",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case: Update Weight and Nutrition", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:       "correctPassword",
			Username:       "GoodDy",
			Weight:         69,
			Protein:        115,
			Fat:            50,
			Carb:           100,
			FavoriteMenues: "11,12",
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:  "gooddy20",
			Weight:  69,
			Protein: 115,
			Fat:     50,
			Carb:    100,
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case: Update Favorite Menues Case 1", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:       "correctPassword",
			Username:       "GoodDy",
			Weight:         71,
			Protein:        120,
			Fat:            60,
			Carb:           120,
			FavoriteMenues: "11,12,13",
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:         "gooddy20",
			FavoriteMenues: "11,12,13",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case: Update Favorite Menues Case 2", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:       "correctPassword",
			Username:       "GoodDy",
			Weight:         71,
			Protein:        120,
			Fat:            60,
			Carb:           120,
			FavoriteMenues: "",
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId: "gooddy20",
		})
		assert.ErrorIs(t, err, nil)
	})
	t.Run("No The User Id", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, sql.ErrNoRows)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:         "gooddy20",
			FavoriteMenues: "11,12,13",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is not found"})
		repo.AssertNotCalled(t, "UpdateUser")
	})
	t.Run("Get User Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, sql.ErrConnDone)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:         "gooddy20",
			FavoriteMenues: "11,12,13",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateUser")
	})
	t.Run("Invalid Password", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:   "gooddy20",
			Password: "p a s s",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "Password need to contain more than 5 letter and no whitespace"})
	})
	t.Run("Invalid Username", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:   "gooddy20",
			Username: "good",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "Username need to contain more than 5 letter and alphabet only"})
	})
	t.Run("Update User Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:       "correctPassword",
			Username:       "GoodDyInwZa20",
			Weight:         71,
			Protein:        120,
			Fat:            60,
			Carb:           120,
			FavoriteMenues: "11,12",
		}).Return(sql.ErrConnDone)
		srv := service.NewUserService(repo)
		err := srv.UpdateUser(service.UpdateUserRequest{
			UserId:   "gooddy20",
			Username: "GoodDyInwZa20",
		})
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}

func TestRecoverFavoriteMenues(t *testing.T) {
	t.Run("Success Case 1", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "12,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.RecoverFavoriteMenues("gooddy20", 11)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case 2", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12,14",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.RecoverFavoriteMenues("gooddy20", 17)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case 3", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.RecoverFavoriteMenues("gooddy20", 12)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("Success Case: No The Deleted Menu Id in Favorite Menues", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(nil)
		srv := service.NewUserService(repo)
		err := srv.RecoverFavoriteMenues("gooddy20", 18)
		assert.ErrorIs(t, err, nil)
	})
	t.Run("No The User Id", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, sql.ErrNoRows)
		srv := service.NewUserService(repo)
		err := srv.RecoverFavoriteMenues("gooddy20", 18)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is not found"})
		repo.AssertNotCalled(t, "UpdateUser")
	})
	t.Run("Get User Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{}, sql.ErrConnDone)
		srv := service.NewUserService(repo)
		err := srv.RecoverFavoriteMenues("gooddy20", 18)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
		repo.AssertNotCalled(t, "UpdateUser")
	})
	t.Run("Update User Database Error", func(t *testing.T) {
		repo := repository.NewUserRepositoryMock()
		repo.On("GetUserById", "gooddy20").Return(&repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,12,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}, nil)
		repo.On("UpdateUser", repository.User{UserId: "gooddy20",
			Password:         "correctPassword",
			Username:         "GoodDy",
			Weight:           71,
			Protein:          120,
			Fat:              60,
			Carb:             120,
			FavoriteMenues:   "11,14,17",
			CreatedTimestamp: time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC).UTC(),
		}).Return(sql.ErrConnDone)
		srv := service.NewUserService(repo)
		err := srv.RecoverFavoriteMenues("gooddy20", 12)
		assert.ErrorIs(t, err, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"})
	})
}
