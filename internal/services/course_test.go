package services

//course_test.go tests ./course.go. This file does not fully utilize table based tests (TBT) because the number of SQL queries is highly variable for some service functions.
//While TBTs would reduce repeated code, they would contain an overabundance of if statements and be less accessible to understand.

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"tech-challenge/internal/models"
	"tech-challenge/internal/testutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (s *testSuit) TestGetAllCourses() {
	t := s.T()

	courses := []models.Course{
		{ID: 0, Name: "My fun GO class"},
		{ID: 1, Name: "Unit Testing 101"},
		{ID: 2, Name: "Table Driven Testing"},
		{ID: 3, Name: "Database Transactions and Hot Chocolate"},
	}
	testCases := map[string]struct {
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		expectedReturn []models.Course
		expectedErr    error
	}{
		"GetAllSuccess": {
			mockReturn:     testutil.MustStructsToRows(courses),
			mockReturnErr:  nil,
			expectedReturn: courses,
			expectedErr:    nil,
		},
		"GetAllEmptyDb": {
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  nil,
			expectedReturn: []models.Course(nil),
			expectedErr:    nil,
		},
		"QueryError": {
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  errors.New("can't query"),
			expectedReturn: []models.Course{},
			expectedErr:    fmt.Errorf("failed to get courses: %w", errors.New("can't query")),
		},
	}
	for testName, testConditions := range testCases {
		t.Run(testName, func(t *testing.T) {
			query := `SELECT * FROM "course"`
			s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testConditions.mockReturn).WillReturnError(testConditions.mockReturnErr)

			actualReturn, err := s.realCourseService.GetAllCourses()
			assert.Equal(t, testConditions.expectedErr, err)
			assert.Equal(t, testConditions.expectedReturn, actualReturn)
			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func (s *testSuit) TestGetCourse() {
	t := s.T()

	courses := []models.Course{
		{ID: 0, Name: "My fun GO class"},
		{ID: 1, Name: "Unit Testing 101"},
		{ID: 2, Name: "Table Driven Testing"},
		{ID: 3, Name: "Database Transactions and Hot Chocolate"},
	}
	testCases := map[string]struct {
		mockReturn     *sqlmock.Rows
		mockId         int
		mockReturnErr  error
		expectedReturn models.Course
		expectedErr    error
	}{
		"ServerError": {
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  errors.New("can't query"),
			mockId:         0,
			expectedErr:    fmt.Errorf("failed to get course: %w", errors.New("can't query")),
			expectedReturn: models.Course{},
		},
		"GetSuccess": {
			mockReturn:     testutil.MustStructsToRows([]models.Course{courses[0]}),
			mockReturnErr:  nil,
			mockId:         0,
			expectedReturn: courses[0],
			expectedErr:    nil,
		},
		"NotFound": {
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  nil,
			mockId:         5,
			expectedReturn: models.Course{},
			expectedErr:    fmt.Errorf("course not found"),
		},
	}
	for testName, testConditions := range testCases {
		t.Run(testName, func(t *testing.T) {

			query := `SELECT * FROM "course" 
							WHERE "id" = $1 
							LIMIT 1`
			s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testConditions.mockReturn).WillReturnError(testConditions.mockReturnErr)

			actualReturn, err := s.realCourseService.GetCourse(testConditions.mockId)
			assert.Equal(t, testConditions.expectedErr, err, testName)
			assert.Equal(t, testConditions.expectedReturn, actualReturn, testName)
			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func (s *testSuit) TestUpdateCourse() {
	t := s.T()

	courseInput := models.Course{Name: "My Fun GO Class"}
	courseOutput := models.Course{ID: 0, Name: "My Fun GO Class"}
	testCases := map[string]struct {
		mockInputArgs  []driver.Value
		mockReturn     driver.Result
		mockReturnErr  error
		inputID        int
		inputCourse    models.Course
		expectedReturn models.Course
		expectedErr    error
	}{
		"ServerError": {
			mockInputArgs:  []driver.Value{courseInput.Name, 0},
			mockReturn:     nil,
			mockReturnErr:  errors.New("can't update"),
			inputID:        0,
			inputCourse:    courseInput,
			expectedReturn: models.Course{},
			expectedErr:    fmt.Errorf("failed to update course: %w", errors.New("can't update")),
		},
		"CourseNotFound": {
			mockInputArgs:  []driver.Value{courseInput.Name, 9},
			mockReturn:     sqlmock.NewResult(0, 0),
			mockReturnErr:  nil,
			inputID:        9,
			inputCourse:    courseInput,
			expectedReturn: models.Course{},
			expectedErr:    fmt.Errorf("course not found"),
		},
		"Success": {
			mockInputArgs:  []driver.Value{courseInput.Name, 0},
			mockReturn:     sqlmock.NewResult(1, 1),
			mockReturnErr:  nil,
			inputID:        0,
			inputCourse:    courseInput,
			expectedReturn: courseOutput,
			expectedErr:    nil,
		},
	}
	for testName, testConditions := range testCases {
		t.Run(testName, func(t *testing.T) {

			query := `UPDATE "course" 
						SET "name" = $1
						WHERE "id" = $2`
			s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(testConditions.mockInputArgs...).WillReturnResult(testConditions.mockReturn).WillReturnError(testConditions.mockReturnErr)

			actualReturn, err := s.realCourseService.UpdateCourse(testConditions.inputID, testConditions.inputCourse)
			assert.Equal(t, testConditions.expectedErr, err, testName)
			assert.Equal(t, testConditions.expectedReturn, actualReturn, testName)
			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

// For TestDeleteCourse(), we do not use table based testing because some tests require multiple queries while others
// only one. There are multiple queries in this transaction.
func (s *testSuit) TestDeleteCourseSuccess() {
	t := s.T()

	courseID := 1

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE "course_id" = $1`)).WithArgs(courseID).WillReturnResult(sqlmock.NewResult(1, 5)).WillReturnError(nil)

	s.dbMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "course" WHERE "id" = $1`)).WithArgs(courseID).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)
	s.dbMock.ExpectCommit()

	rowsAffected, err := s.realCourseService.DeleteCourse(courseID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, int64(1))

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeleteCourseTransactionFailure() {
	t := s.T()

	s.dbMock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

	rowsAffected, err := s.realCourseService.DeleteCourse(1)
	assert.Equal(t, err, fmt.Errorf("failed to begin transaction: %w", errors.New("transaction begin error")))
	assert.Equal(t, int64(-1), rowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeleteCourseDeleteRelationsFailure() {
	t := s.T()

	courseID := 1

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE "course_id" = $1`)).WithArgs(courseID).WillReturnResult(sqlmock.NewResult(1, 5)).WillReturnError(errors.New("can't delete relations"))
	s.dbMock.ExpectRollback()

	rowsAffected, err := s.realCourseService.DeleteCourse(courseID)
	assert.Equal(t, err, fmt.Errorf("failed to delete course relations: %w", errors.New("can't delete relations")))
	assert.Equal(t, int64(-1), rowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeleteCourseFailure() {
	t := s.T()

	courseID := 1

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE "course_id" = $1`)).WithArgs(courseID).WillReturnResult(sqlmock.NewResult(1, 5)).WillReturnError(nil)
	s.dbMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "course" WHERE "id" = $1`)).WithArgs(courseID).WillReturnResult(sqlmock.NewResult(int64(courseID), 1)).WillReturnError(errors.New("can't delete course"))

	s.dbMock.ExpectRollback()

	rowsAffected, err := s.realCourseService.DeleteCourse(courseID)
	assert.Equal(t, err, fmt.Errorf("failed to delete course with ID: %v. %w", courseID, errors.New("can't delete course")))
	assert.Equal(t, int64(-1), rowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeleteCourseCommitTransactionFailure() {
	t := s.T()

	courseID := 1

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE "course_id" = $1`)).WithArgs(courseID).WillReturnResult(sqlmock.NewResult(1, 5)).WillReturnError(nil)
	s.dbMock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "course" WHERE "id" = $1`)).WithArgs(courseID).WillReturnResult(sqlmock.NewResult(int64(courseID), 1)).WillReturnError(nil)
	s.dbMock.ExpectCommit().WillReturnError(errors.New("can't commit"))

	rowsAffected, err := s.realCourseService.DeleteCourse(courseID)

	assert.Equal(t, err, fmt.Errorf("failed to commit transaction: %w", errors.New("can't commit")))
	assert.Equal(t, int64(-1), rowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreateCourseSuccess() {

	t := s.T()

	insertCourse := models.Course{Name: "new course"}
	expectedReturnCourseID := 1

	query := `INSERT INTO "course" (name) VALUES ($1) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(insertCourse.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedReturnCourseID))

	returnedCourse, err := s.realCourseService.CreateCourse(insertCourse)

	assert.Equal(t, expectedReturnCourseID, returnedCourse)
	assert.NoError(t, err)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreateCourseFailure() {

	t := s.T()

	insertCourse := models.Course{Name: "new course"}
	expectedReturnCourseID := -1
	expectedError := fmt.Errorf("failed to create course: %w", errors.New("can't create course"))

	query := `INSERT INTO "course" (name) VALUES ($1) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(insertCourse.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).WillReturnError(errors.New("can't create course"))

	returnedCourse, err := s.realCourseService.CreateCourse(insertCourse)

	assert.Equal(t, expectedReturnCourseID, returnedCourse)
	assert.Equal(t, expectedError, err)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
