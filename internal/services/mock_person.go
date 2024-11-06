package services

//mock_person.go is used for testing purposes in ../handlers/person_test.go

import (
	"tech-challenge/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockPersonService struct {
	mock.Mock
}

func (s *MockPersonService) GetAllPeople(age int, firstName string, lastName string) ([]models.Person, error) {
	args := s.Called(age, firstName, lastName)
	return args.Get(0).([]models.Person), args.Error(1)
}
func (s *MockPersonService) GetPerson(firstName string, lastName string) (models.Person, error) {
	args := s.Called(firstName, lastName)
	return args.Get(0).(models.Person), args.Error(1)
}
func (s *MockPersonService) UpdatePerson(firstName string, lastName string, person models.Person) (models.Person, error) {
	args := s.Called(firstName, lastName, person)
	return args.Get(0).(models.Person), args.Error(1)
}
func (s *MockPersonService) CreatePerson(person models.Person) (int, error) {
	args := s.Called(person)
	return args.Get(0).(int), args.Error(1)
}
func (s *MockPersonService) DeletePerson(firstName string, lastName string) (int64, error) {
	args := s.Called(firstName, lastName)
	return args.Get(0).(int64), args.Error(1)
}
