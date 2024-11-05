package services

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuit struct {
	suite.Suite
	realCourseService *RealCourseService
	personService     *RealPersonService
	dbMock            sqlmock.Sqlmock
}
type Person_Course struct {
	PersonID int
	CourseID int
}
type PersonDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Type      string `json:"type" validate:"required,ValidateType"`
	Age       int    `json:"age" validate:"required,gt=0"`
}
type ID struct {
	ID int
}

// calls SetupSuite() runs all test functions in this .go file that begin with "Test", then calls TearDownSuite()
func TestTestSuit(t *testing.T) {
	suite.Run(t, new(testSuit))
}
func (s *testSuit) SetupSuite() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)

	s.dbMock = mock
	s.realCourseService = NewCourseService(db)
	s.personService = NewPersonService(db)
}
func (s *testSuit) TearDownSuite() {
	s.realCourseService.db.Close()
	s.personService.db.Close()
}
