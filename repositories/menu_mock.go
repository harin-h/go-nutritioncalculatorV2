package repository

import "github.com/stretchr/testify/mock"

type menuRepositoryMock struct {
	mock.Mock
}

func NewMenuRepositoryMock() *menuRepositoryMock {
	return &menuRepositoryMock{}
}

func (r *menuRepositoryMock) CreateMenu(menu Menu) (*Menu, error) {
	args := r.Called(menu)
	return args.Get(0).(*Menu), args.Error(1)
}

func (r *menuRepositoryMock) GetAllMenues() ([]Menu, error) {
	args := r.Called()
	return args.Get(0).([]Menu), args.Error(1)
}

func (r *menuRepositoryMock) GetMenuById(menuId int) (*Menu, error) {
	args := r.Called(menuId)
	return args.Get(0).(*Menu), args.Error(1)
}

func (r *menuRepositoryMock) UpdateMenu(menu Menu) error {
	args := r.Called(menu)
	return args.Error(0)
}
