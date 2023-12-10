package service

import "github.com/stretchr/testify/mock"

type menuServiceMock struct {
	mock.Mock
}

func NewMenuServiceMock() *menuServiceMock {
	return &menuServiceMock{}
}

func (s *menuServiceMock) CreateMenu(newMenu NewMenuRequest) error {
	args := s.Called(newMenu)
	return args.Error(0)
}

func (s *menuServiceMock) GetAllMenues() ([]MenuResponse, error) {
	args := s.Called()
	return args.Get(0).([]MenuResponse), args.Error(1)
}

func (s *menuServiceMock) UpdateMenu(updateMenu UpdateMenuRequest) error {
	args := s.Called(updateMenu)
	return args.Error(0)
}

func (s *menuServiceMock) RecoverMenu(menuId int, name string) (*MenuResponse, error) {
	args := s.Called(menuId, name)
	return args.Get(0).(*MenuResponse), args.Error(1)
}

func (s *menuServiceMock) DeleteMenu(menuId int) error {
	args := s.Called(menuId)
	return args.Error(0)
}
