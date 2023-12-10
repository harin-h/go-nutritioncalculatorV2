package service

import (
	"database/sql"
	"go-nutritioncalculator2/errs"
	"go-nutritioncalculator2/logs"
	repository "go-nutritioncalculator2/repositories"
	"net/http"
	"time"
)

type menuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) menuService {
	return menuService{menuRepo: menuRepo}
}

func (s menuService) CreateMenu(newMenu NewMenuRequest) error {
	menu := repository.Menu{
		Name:             newMenu.Name,
		Protein:          newMenu.Protein,
		Fat:              newMenu.Fat,
		Carb:             newMenu.Carb,
		CreatorId:        newMenu.CreatorId,
		Status:           1,
		CreatedTimestamp: time.Now().UTC().Truncate(time.Second),
	}
	_, err := s.menuRepo.CreateMenu(menu)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s menuService) GetAllMenues() ([]MenuResponse, error) {
	menues, err := s.menuRepo.GetAllMenues()
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	menuesRes := []MenuResponse{}
	for i := 0; i < len(menues); i++ {
		menu := MenuResponse{
			Id:          menues[i].Id,
			Name:        menues[i].Name,
			Protein:     menues[i].Protein,
			Fat:         menues[i].Fat,
			Carb:        menues[i].Carb,
			CreatorId:   menues[i].CreatorId,
			CreatorName: menues[i].CreatorName,
			Like:        menues[i].Like,
			Status:      menues[i].Status,
		}
		menuesRes = append(menuesRes, menu)
	}
	return menuesRes, nil
}

func (s menuService) UpdateMenu(updateMenu UpdateMenuRequest) error {
	err := s.menuRepo.UpdateMenu(repository.Menu{Id: updateMenu.Id})
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	menu, err := s.menuRepo.GetMenuById(updateMenu.Id)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	menu.Id = 0
	menu.Status = 1
	if updateMenu.Name != "" {
		menu.Name = updateMenu.Name
	}
	if updateMenu.Protein != menu.Protein {
		menu.Protein = updateMenu.Protein
	}
	if updateMenu.Fat != menu.Fat {
		menu.Fat = updateMenu.Fat
	}
	if updateMenu.Carb != menu.Carb {
		menu.Carb = updateMenu.Carb
	}
	menu.CreatedTimestamp = time.Now().UTC().Truncate(time.Second)
	_, err = s.menuRepo.CreateMenu(*menu)
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}

func (s menuService) RecoverMenu(menuId int, name string) (*MenuResponse, error) {
	menu, err := s.menuRepo.GetMenuById(menuId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.AppError{Code: http.StatusNotAcceptable, Message: "Menu Id is not found"}
		}
		logs.Error(err)
		return nil, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	menu.Id = 0
	menu.Status = 1
	menu.CreatedTimestamp = time.Now().UTC().Truncate(time.Second)
	if name != "" {
		menu.Name = name
	}
	newMenu, err := s.menuRepo.CreateMenu(*menu)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	menuRes := MenuResponse{
		Id:          newMenu.Id,
		Name:        newMenu.Name,
		Protein:     newMenu.Protein,
		Fat:         newMenu.Fat,
		Carb:        newMenu.Carb,
		CreatorId:   newMenu.CreatorId,
		CreatorName: newMenu.CreatorName,
		Like:        newMenu.Like,
		Status:      newMenu.Status,
	}
	return &menuRes, nil
}

func (s menuService) DeleteMenu(menuId int) error {
	err := s.menuRepo.UpdateMenu(repository.Menu{Id: menuId})
	if err != nil {
		logs.Error(err)
		return errs.AppError{Code: http.StatusInternalServerError, Message: "Unexpected error"}
	}
	return nil
}
