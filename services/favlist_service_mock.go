package service

import (
	"github.com/stretchr/testify/mock"
)

type favListServiceMock struct {
	mock.Mock
}

func NewFavListServiceMock() *favListServiceMock {
	return &favListServiceMock{}
}

func (s *favListServiceMock) GetFavListsByUserId(userId string) ([]FavListResponse, error) {
	args := s.Called(userId)
	return args.Get(0).([]FavListResponse), args.Error(1)
}

func (s *favListServiceMock) CreateFavList(newFavListReq NewFavListRequest) error {
	args := s.Called(newFavListReq)
	return args.Error(0)
}

func (s *favListServiceMock) DeleteFavList(favListId int) error {
	args := s.Called(favListId)
	return args.Error(0)
}

func (s *favListServiceMock) UpdateFavList(updateFavListReq UpdateFavListRequest) error {
	args := s.Called(updateFavListReq)
	return args.Error(0)
}

func (s *favListServiceMock) RecoverFavList(favListId int, oldMenuId int, newMenuId int) error {
	args := s.Called(favListId, oldMenuId, newMenuId)
	return args.Error(0)
}
