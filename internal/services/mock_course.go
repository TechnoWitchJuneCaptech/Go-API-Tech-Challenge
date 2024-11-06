package services

//mock_course.go is used for testing purposes in ../handlers/course_test.go

import (
	"tech-challenge/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockCourseService struct {
	mock.Mock
}

func (s *MockCourseService) GetAllCourses() ([]models.Course, error) {
	args := s.Called()
	return args.Get(0).([]models.Course), args.Error(1)
}
func (s *MockCourseService) GetCourse(id int) (models.Course, error) {
	args := s.Called(id)
	return args.Get(0).(models.Course), args.Error(1)
}
func (s *MockCourseService) UpdateCourse(id int, course models.Course) (models.Course, error) {
	args := s.Called(id, course)
	return args.Get(0).(models.Course), args.Error(1)
}
func (s *MockCourseService) CreateCourse(course models.Course) (int, error) {
	args := s.Called(course)
	return args.Get(0).(int), args.Error(1)
}
func (s *MockCourseService) DeleteCourse(id int) (int64, error) {
	args := s.Called(id)
	return args.Get(0).(int64), args.Error(1)
}
