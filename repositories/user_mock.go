package repository

import "github.com/stretchr/testify/mock"

type userRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *userRepositoryMock {
	return &userRepositoryMock{}
}

func (r *userRepositoryMock) GetUserById(userId string) (*User, error) {
	args := r.Called(userId)
	return args.Get(0).(*User), args.Error(1)
}

func (r *userRepositoryMock) GetUserByUsername(username string) (*User, error) {
	args := r.Called(username)
	return args.Get(0).(*User), args.Error(1)
}

func (r *userRepositoryMock) CreateUser(user User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *userRepositoryMock) UpdateUser(user User) error {
	args := r.Called(user)
	return args.Error(0)
}
