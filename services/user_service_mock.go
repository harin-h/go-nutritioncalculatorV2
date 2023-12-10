package service

import "github.com/stretchr/testify/mock"

type userServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *userServiceMock {
	return &userServiceMock{}
}

func (s *userServiceMock) CheckLogIn(logInReq LogInRequest) (*LogInResponse, error) {
	args := s.Called(logInReq)
	return args.Get(0).(*LogInResponse), args.Error(1)
}

func (s *userServiceMock) GetUserDetail(userId string) (*UserResponse, error) {
	args := s.Called(userId)
	return args.Get(0).(*UserResponse), args.Error(1)
}

func (s *userServiceMock) CreateUser(newUser NewUserRequest) error {
	args := s.Called(newUser)
	return args.Error(0)
}

func (s *userServiceMock) UpdateUser(newUpdateUser UpdateUserRequest) error {
	args := s.Called(newUpdateUser)
	return args.Error(0)
}

func (s *userServiceMock) RecoverFavoriteMenues(userId string, deletedMenuId int) error {
	args := s.Called(userId, deletedMenuId)
	return args.Error(0)
}
