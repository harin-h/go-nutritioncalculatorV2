package service

import (
	"database/sql"
	"go-nutritioncalculator2/errs"
	"go-nutritioncalculator2/logs"
	repository "go-nutritioncalculator2/repositories"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) userService {
	return userService{userRepo: userRepo}
}

func (s userService) CheckLogIn(logInReq LogInRequest) (*LogInResponse, error) {
	user, err := s.userRepo.GetUserById(logInReq.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &LogInResponse{IsLogIn: false}, nil
		}
		logs.Error(err)
		return nil, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	if user.Password == logInReq.Password {
		return &LogInResponse{IsLogIn: true}, nil
	}
	return &LogInResponse{IsLogIn: false}, nil
}

func (s userService) GetUserDetail(userId string) (*userResponse, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is not found"}
		}
		logs.Error(err)
		return nil, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	userRes := userResponse{
		Username:       user.Username,
		Weight:         user.Weight,
		Protein:        user.Protein,
		Fat:            user.Fat,
		Carb:           user.Carb,
		FavoriteMenues: user.FavoriteMenues,
	}
	return &userRes, nil
}

func (s userService) CreateUser(newUser NewUserRequest) error {
	user := repository.User{
		UserId:           newUser.UserId,
		Password:         newUser.Password,
		Username:         newUser.Username,
		Weight:           newUser.Weight,
		Protein:          newUser.Protein,
		Fat:              newUser.Fat,
		Carb:             newUser.Carb,
		FavoriteMenues:   "",
		CreatedTimestamp: time.Now().UTC(),
	}
	var err error
	var isOk bool
	isOk, err = regexp.MatchString(`\w{6,}`, user.UserId)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	} else if len(user.UserId) < 6 || !isOk {
		return errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id need to contain more than 5 letter and alphabet only"}
	}
	isOk, err = regexp.MatchString(`\S{6,}`, user.Password)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	} else if len(user.Password) < 6 || !isOk {
		return errs.AppError{Code: http.StatusNotAcceptable, Message: "Password need to contain more than 5 letter and no whitespace"}
	}
	isOk, err = regexp.MatchString(`\w{6,}`, user.Username)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	} else if len(user.Username) < 6 || !isOk {
		return errs.AppError{Code: http.StatusNotAcceptable, Message: "Username need to contain more than 5 letter and alphabet only"}
	}
	_, err = s.userRepo.GetUserById(user.UserId)
	if err == nil {
		return errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is already used"}
	}
	_, err = s.userRepo.GetUserByUsername(user.Username)
	if err == nil {
		return errs.AppError{Code: http.StatusNotAcceptable, Message: "Username is already used"}
	}
	err = s.userRepo.CreateUser(user)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s userService) UpdateUser(newUpdateUser UpdateUserRequest) error {
	var err error
	var isOk bool
	updateUser := repository.User{
		UserId:         newUpdateUser.UserId,
		Password:       newUpdateUser.Password,
		Username:       newUpdateUser.Username,
		Weight:         newUpdateUser.Weight,
		Protein:        newUpdateUser.Protein,
		Fat:            newUpdateUser.Fat,
		Carb:           newUpdateUser.Carb,
		FavoriteMenues: newUpdateUser.FavoriteMenues,
	}
	user, err := s.userRepo.GetUserById(updateUser.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is not found"}
		}
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	if updateUser.Password != "" {
		isOk, err = regexp.MatchString(`\S{6,}`, updateUser.Password)
		if err != nil {
			logs.Error(err)
			return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
		} else if len(updateUser.Password) < 6 || !isOk {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: "Password need to contain more than 5 letter and no whitespace"}
		}
	} else {
		updateUser.Password = user.Password
	}
	if updateUser.Username != "" {
		isOk, err = regexp.MatchString(`\w{6,}`, updateUser.Username)
		if err != nil {
			logs.Error(err)
			return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
		} else if len(updateUser.Username) < 6 || !isOk {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: "Username need to contain more than 5 letter and alphabet only"}
		}
	} else {
		updateUser.Username = user.Username
	}
	if updateUser.Weight == 0 {
		updateUser.Weight = user.Weight
	}
	if updateUser.Protein == 0 {
		updateUser.Protein = user.Protein
	}
	if updateUser.Fat == 0 {
		updateUser.Fat = user.Fat
	}
	if updateUser.Carb == 0 {
		updateUser.Carb = user.Carb
	}
	if updateUser.FavoriteMenues == "" && updateUser != (repository.User{UserId: newUpdateUser.UserId}) {
		updateUser.FavoriteMenues = user.FavoriteMenues
	}
	err = s.userRepo.UpdateUser(updateUser)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s userService) RecoverFavoriteMenues(userId string, deletedMenuId int) error {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: "User Id is not found"}
		}
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	tempFavoriteMenues := strings.Split(user.FavoriteMenues, ",")
	index := -1
	for i := 0; i < len(tempFavoriteMenues); i++ {
		if tempFavoriteMenues[i] == strconv.Itoa(deletedMenuId) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	if index == 0 {
		tempFavoriteMenues = tempFavoriteMenues[1:]
	} else if index == len(tempFavoriteMenues)-1 {
		tempFavoriteMenues = tempFavoriteMenues[:len(tempFavoriteMenues)-1]
	} else {
		tempFavoriteMenues = append(tempFavoriteMenues[:index], tempFavoriteMenues[index+1:]...)
	}
	user.FavoriteMenues = strings.Join(tempFavoriteMenues, ",")
	err = s.userRepo.UpdateUser(*user)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}
