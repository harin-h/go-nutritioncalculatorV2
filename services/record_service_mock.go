package service

import "github.com/stretchr/testify/mock"

type recordServiceMock struct {
	mock.Mock
}

func NewRecordServiceMock() *recordServiceMock {
	return &recordServiceMock{}
}

func (s *recordServiceMock) GetAllRecordsByUserId(userId string) ([]RecordResponse, error) {
	args := s.Called(userId)
	return args.Get(0).([]RecordResponse), args.Error(1)
}

func (s *recordServiceMock) CreateRecord(newRecordReq NewRecordRequest) error {
	args := s.Called(newRecordReq)
	return args.Error(0)
}

func (s *recordServiceMock) DeleteRecord(recordId int) error {
	args := s.Called(recordId)
	return args.Error(0)
}

func (s *recordServiceMock) UpdateRecord(updateRecordReq UpdateRecordRequest) error {
	args := s.Called(updateRecordReq)
	return args.Error(0)
}
