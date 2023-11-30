package service

import (
	"database/sql"
	"fmt"
	"go-nutritioncalculator2/errs"
	"go-nutritioncalculator2/logs"
	repository "go-nutritioncalculator2/repositories"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type favListService struct {
	favListRepo repository.FavListRepository
}

func NewFavListService(favListRepo repository.FavListRepository) favListService {
	return favListService{favListRepo: favListRepo}
}

func (s favListService) GetFavListsByUserId(userId string) ([]FavListResponse, error) {
	favLists, err := s.favListRepo.GetFavListsByUserId(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []FavListResponse{}, nil
		}
		logs.Error(err)
		return nil, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	favListsRes := []FavListResponse{}
	for i := 0; i < len(favLists); i++ {
		favList := FavListResponse{
			Id:        favLists[i].Id,
			Name:      favLists[i].Name,
			Menues:    favLists[i].Menues,
			List:      favLists[i].List,
			Protein:   favLists[i].Protein,
			Fat:       favLists[i].Fat,
			Carb:      favLists[i].Carb,
			IsUpdated: favLists[i].IsUpdated,
		}
		favListsRes = append(favListsRes, favList)
	}
	return favListsRes, nil
}

func (s favListService) CreateFavList(newFavListReq NewFavListRequest) error {
	newFavList := repository.FavList{
		UserId:           newFavListReq.UserId,
		Name:             newFavListReq.Name,
		List:             newFavListReq.List,
		Status:           1,
		CreatedTimestamp: time.Now().UTC(),
	}
	_, err := s.favListRepo.CreateFavList(newFavList)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s favListService) DeleteFavList(favListId int) error {
	favList, err := s.favListRepo.GetFavListById(favListId)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	favList.Status = 0
	err = s.favListRepo.UpdateFavList(*favList)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s favListService) UpdateFavList(updateFavListReq UpdateFavListRequest) error {
	favList, err := s.favListRepo.GetFavListById(updateFavListReq.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Favorite List Id - ", updateFavListReq.Id, "is not found")}
		}
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	if updateFavListReq.Name != "" {
		favList.Name = updateFavListReq.Name
	}
	if updateFavListReq.List != "" {
		favList.List = updateFavListReq.List
	}
	err = s.favListRepo.UpdateFavList(*favList)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s favListService) RecoverFavList(favListId int, oldMenuId int, newMenuId int) error {
	favList, err := s.favListRepo.GetFavListById(favListId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AppError{Code: http.StatusNotAcceptable, Message: fmt.Sprint("Favorite List Id - ", favList.Id, "is not found")}
		}
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	tempFavList := strings.Split(favList.List, ",")
	minIndex := -1
	maxIndex := -1
	for i := 0; i < len(tempFavList); i++ {
		tempMenuId, err := strconv.Atoi(tempFavList[i])
		if err != nil {
			logs.Error(err)
			return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
		}
		if tempMenuId < oldMenuId {
			continue
		} else if tempMenuId == oldMenuId {
			if minIndex == -1 {
				minIndex = i
			} else {
				maxIndex = i
			}
		} else {
			break
		}
	}
	if minIndex == -1 {
		return nil
	}
	if maxIndex == -1 {
		maxIndex = minIndex
	}
	if minIndex == 0 {
		tempFavList = tempFavList[maxIndex+1:]
	} else if maxIndex == len(tempFavList)-1 {
		tempFavList = tempFavList[:minIndex]
	} else {
		tempFavList = append(tempFavList[:minIndex], tempFavList[maxIndex+1:]...)
	}
	if newMenuId != 0 {
		tempNewMenuId := strconv.Itoa(newMenuId)
		for i := 0; i <= maxIndex-minIndex; i++ {
			tempFavList = append(tempFavList, tempNewMenuId)
		}
	}
	favList.List = strings.Join(tempFavList, ",")
	err = s.favListRepo.UpdateFavList(*favList)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}
