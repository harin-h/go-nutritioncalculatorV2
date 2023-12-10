package repository

import "github.com/stretchr/testify/mock"

type favListRepositoryMock struct {
	mock.Mock
}

func NewFavListRepositoryMock() *favListRepositoryMock {
	return &favListRepositoryMock{}
}

func (r *favListRepositoryMock) GetFavListsByUserId(userId string) ([]FavList, error) {
	args := r.Called(userId)
	return args.Get(0).([]FavList), args.Error(1)
}

func (r *favListRepositoryMock) GetFavListById(favListId int) (*FavList, error) {
	args := r.Called(favListId)
	return args.Get(0).(*FavList), args.Error(1)
}

func (r *favListRepositoryMock) CreateFavList(favList FavList) (*FavList, error) {
	args := r.Called(favList)
	return args.Get(0).(*FavList), args.Error(1)
}

func (r *favListRepositoryMock) UpdateFavList(favList FavList) error {
	args := r.Called(favList)
	return args.Error(0)
}
