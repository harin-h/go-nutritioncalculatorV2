package repository

import "github.com/stretchr/testify/mock"

type recordRepositoryMock struct {
	mock.Mock
}

func NewRecordRepositoryMock() *recordRepositoryMock {
	return &recordRepositoryMock{}
}

func (r *recordRepositoryMock) GetRecordsByUserId(userId string) ([]Record, error) {
	args := r.Called(userId)
	return args.Get(0).([]Record), args.Error(1)
}

func (r *recordRepositoryMock) GetRecordById(recordId int) (*Record, error) {
	args := r.Called(recordId)
	return args.Get(0).(*Record), args.Error(1)
}

func (r *recordRepositoryMock) CreateRecord(record Record) (*Record, error) {
	args := r.Called(record)
	return args.Get(0).(*Record), args.Error(1)
}

func (r *recordRepositoryMock) UpdateRecord(record Record) error {
	args := r.Called(record)
	return args.Error(0)
}
