package services

//person_test.go tests ./person.go. This file does not fully utilize table based tests (TBT) because the number of SQL queries is highly variable for some service functions.
//While TBTs would reduce repeated code, they would contain an overabundance of if statements and be less accessible to understand.

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"tech-challenge/internal/models"
	"tech-challenge/internal/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Tests for GetAllPeople()
// NameAndAgeSuccess
// AgeSuccess
// NameSuccess
// NoArgsSuccess
// EmptyListSuccess
// FailedToGetPeople
// FailedToGetCoursesForPerson
func (s *testSuit) TestGetAllPeopleNameAgeSuccess() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[3:])
	returnRowsMapQuery := testutil.MustStructsToRows(person_course[9:])
	firstName := "Bubbles"
	lastName := "Thane"
	age := 18
	returnFinal := []models.Person{{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18, Courses: []int{1, 2, 3}}}

	query := `SELECT * FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) AND age = $3`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName, age).WillReturnRows(returnRowsPersonQuery).WillReturnError(nil)
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery).WillReturnError(nil)

	result, err := s.personService.GetAllPeople(age, firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, err, nil)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetAllPeopleNameSuccess() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[3:])
	returnRowsMapQuery := testutil.MustStructsToRows(person_course[9:])
	firstName := "Bubbles"
	lastName := "Thane"
	age := -1
	returnFinal := []models.Person{{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18, Courses: []int{1, 2, 3}}}

	query := `SELECT * FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2)`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(returnRowsPersonQuery).WillReturnError(nil)
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery).WillReturnError(nil)
	result, err := s.personService.GetAllPeople(age, firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, err, nil)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetAllPeopleAgeSuccess() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[:2])
	returnRowsMapQuery1 := testutil.MustStructsToRows(person_course[:3])
	returnRowsMapQuery2 := testutil.MustStructsToRows(person_course[3:6])
	firstName := ""
	lastName := ""
	age := 22
	returnFinal := []models.Person{{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22, Courses: []int{1, 2, 3}},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22, Courses: []int{1, 2, 3}},
	}

	query := `SELECT * FROM "person" WHERE age = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(age).WillReturnRows(returnRowsPersonQuery).WillReturnError(nil)
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery1).WillReturnError(nil)
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery2).WillReturnError(nil)
	result, err := s.personService.GetAllPeople(age, firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, err, nil)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetAllPeopleNoArgsSuccess() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people)
	returnRowsMapQuery1 := testutil.MustStructsToRows(person_course[:3])
	returnRowsMapQuery2 := testutil.MustStructsToRows(person_course[3:6])
	returnRowsMapQuery3 := testutil.MustStructsToRows(person_course[6:9])
	returnRowsMapQuery4 := testutil.MustStructsToRows(person_course[9:12])
	firstName := ""
	lastName := ""
	age := -1
	returnFinal := []models.Person{{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22, Courses: []int{1, 2, 3}},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22, Courses: []int{1, 2, 3}},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222, Courses: []int{1, 2, 3}},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18, Courses: []int{1, 2, 3}},
	}

	query := `SELECT * FROM "person"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsPersonQuery).WillReturnError(nil)
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery1).WillReturnError(nil)
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery2).WillReturnError(nil)
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery3).WillReturnError(nil)
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery4).WillReturnError(nil)
	result, err := s.personService.GetAllPeople(age, firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, err, nil)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetAllPeopleEmptySuccess() {
	t := s.T()

	firstName := ""
	lastName := ""
	age := -1
	returnFinal := []models.Person(nil)

	query := `SELECT * FROM "person"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(&sqlmock.Rows{}).WillReturnError(nil)
	result, err := s.personService.GetAllPeople(age, firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, err, nil)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetAllPeopleGetPeopleFailure() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[:2])
	firstName := ""
	lastName := ""
	age := 22
	returnFinal := []models.Person{}
	returnErr := fmt.Errorf("failed to get people: %w", errors.New("can't get people"))

	query := `SELECT * FROM "person" WHERE age = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(age).WillReturnRows(returnRowsPersonQuery).WillReturnError(errors.New("can't get people"))
	result, err := s.personService.GetAllPeople(age, firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, err, returnErr)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetAllPeopleGetCoursesFailure() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[:2])
	returnRowsMapQuery1 := testutil.MustStructsToRows(person_course[:3])
	returnRowsMapQuery2 := testutil.MustStructsToRows(person_course[3:6])
	firstName := ""
	lastName := ""
	age := 22
	returnFinal := []models.Person{}
	returnErr := fmt.Errorf("failed to get courses for person: %w", errors.New("can't get courses"))

	query := `SELECT * FROM "person" WHERE age = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(age).WillReturnRows(returnRowsPersonQuery).WillReturnError(nil)
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery1).WillReturnError(nil)
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery2).WillReturnError(errors.New("can't get courses"))
	result, err := s.personService.GetAllPeople(age, firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Tests for GetPerson()
// Get Person Exists Success
// Get Person Empty Success
// Get Person Get Person Fail
// Get Person Get Course Fail
func (s *testSuit) TestGetPersonExistsSuccess() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[3:])
	returnRowsMapQuery := testutil.MustStructsToRows(person_course[9:])
	firstName := "Bubbles"
	lastName := "Thane"
	returnFinal := models.Person{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18, Courses: []int{1, 2, 3}}

	query := `SELECT * FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(returnRowsPersonQuery).WillReturnError(nil)
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery).WillReturnError(nil)
	result, err := s.personService.GetPerson(firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, nil, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetPersonDoesntExistSuccess() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "NotAProfessor"
	returnFinal := models.Person{}

	query := `SELECT * FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(&sqlmock.Rows{}).WillReturnError(nil)
	result, err := s.personService.GetPerson(firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, nil, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetPersonGetPersonFailure() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[3:])
	firstName := "Bubbles"
	lastName := "Thane"
	returnFinal := models.Person{}
	returnErr := fmt.Errorf("failed to get person: %w", errors.New("can't get person"))

	query := `SELECT * FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(returnRowsPersonQuery).WillReturnError(errors.New("can't get person"))
	result, err := s.personService.GetPerson(firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestGetPersonGetCourseFailure() {
	t := s.T()

	people := []PersonDTO{
		{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
		{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
		{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	}

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	returnRowsPersonQuery := testutil.MustStructsToRows(people[3:])
	returnRowsMapQuery := testutil.MustStructsToRows(person_course[9:])
	firstName := "Bubbles"
	lastName := "Thane"
	returnFinal := models.Person{}
	returnErr := fmt.Errorf("failed to get courses for person: %w", errors.New("can't get course"))

	query := `SELECT * FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(returnRowsPersonQuery).WillReturnError(nil)
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(returnRowsMapQuery).WillReturnError(errors.New("can't get course"))
	result, err := s.personService.GetPerson(firstName, lastName)

	assert.Equal(t, returnFinal, result)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Tests for UpdatePerson()
// UpdatePersonChangeAllSuccess
// UpdatePersonPersonNotFoundFailure
// UpdatePersonUpdatePersonFailure
// UpdatePersonGetIDFailure
// UpdatePersonGetMapListFailure
// UpdatePersonDeleteCoursesFailure
// UpdatePersonGetCourseListFailure
// UpdatePersonCourseNotFoundFailure
// UpdatePersonUpdateCoursesFailure
// UpdatePersonTransactionBeginFailure
// UpdatePersonTransactionCommitFailure
func (s *testSuit) TestUpdatePersonChangeAllSuccess() {
	t := s.T()

	// people := []PersonDTO{
	// 	{ID: 0, FirstName: "Tim", LastName: "Rogers", Type: "student", Age: 22},
	// 	{ID: 1, FirstName: "Jill", LastName: "Rogers", Type: "student", Age: 22},
	// 	{ID: 2, FirstName: "Jack", LastName: "Daniels", Type: "student", Age: 222},
	// 	{ID: 3, FirstName: "Bubbles", LastName: "Thane", Type: "professor", Age: 18},
	// }

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	// courses := []models.Course{
	// 	{ID: 1, Name: "Go for Gophers!"},
	// 	{ID: 2, Name: "Final Fantasy XIV 101"},
	// 	{ID: 3, Name: "Ramen 401"},
	// 	{ID: 4, Name: "Data Structures & Algorithms"},
	// 	{ID: 5, Name: "Intermediate Gerontology"},
	// }

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnPerson := inputPerson

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}}))
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnRows(testutil.MustStructsToRows(person_course[9:]))
	query = `DELETE FROM "person_course" WHERE person_id = $1 AND course_id = ANY ($2::int[])`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(3, pq.Array([]int{1, 2})).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))
	query = `INSERT INTO "person_course" (person_id, course_id) VALUES (3, 4), (3, 5)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(1, 1))
	s.dbMock.ExpectCommit()

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.NoError(t, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonNotFoundFailure() {
	t := s.T()

	personInput := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnErr := fmt.Errorf("person not found")
	returnPerson := models.Person{}

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 0))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", personInput)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonUpdatePersonFailure() {
	t := s.T()

	personInput := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnErr := fmt.Errorf("failed to update person: %w", errors.New("can't update person"))
	returnPerson := models.Person{}

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 0)).WillReturnError(errors.New("can't update person"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", personInput)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonGetIDFailure() {
	t := s.T()

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnErr := fmt.Errorf("failed to retreive id: %w", errors.New("can't get ID"))
	returnPerson := models.Person{}

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}})).WillReturnError(errors.New("can't get ID"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonGetMapListFailure() {
	t := s.T()

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnPerson := models.Person{}
	returnErr := fmt.Errorf("failed to retreive course list: %w", errors.New("can't get map"))

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}}))
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnRows(testutil.MustStructsToRows(person_course[9:])).WillReturnError(errors.New("can't get map"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Error(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonDeleteCoursesFailure() {
	t := s.T()

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnPerson := models.Person{}
	returnErr := fmt.Errorf("failed to update course list: %w", errors.New("can't delete"))

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}}))
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnRows(testutil.MustStructsToRows(person_course[9:]))
	query = `DELETE FROM "person_course" WHERE person_id = $1 AND course_id = ANY ($2::int[])`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(3, pq.Array([]int{1, 2})).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(errors.New("can't delete"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonGetCourseListFailure() {
	t := s.T()

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnPerson := models.Person{}
	returnErr := fmt.Errorf("failed to retreive course list: %w", errors.New("can't get courses"))

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}}))
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnRows(testutil.MustStructsToRows(person_course[9:]))
	query = `DELETE FROM "person_course" WHERE person_id = $1 AND course_id = ANY ($2::int[])`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(3, pq.Array([]int{1, 2})).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}})).WillReturnError(errors.New("can't get courses"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonCourseNotFoundFailure() {
	t := s.T()

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5, 6}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnPerson := models.Person{}
	returnErr := fmt.Errorf("course not found, trying to join a course that doesn't exist")

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}}))
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnRows(testutil.MustStructsToRows(person_course[9:]))
	query = `DELETE FROM "person_course" WHERE person_id = $1 AND course_id = ANY ($2::int[])`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(3, pq.Array([]int{1, 2})).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonUpdateCoursesFailure() {
	t := s.T()

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnPerson := models.Person{}
	returnErr := fmt.Errorf("failed to update course list: %w", errors.New("can't update courses"))

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}}))
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnRows(testutil.MustStructsToRows(person_course[9:]))
	query = `DELETE FROM "person_course" WHERE person_id = $1 AND course_id = ANY ($2::int[])`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(3, pq.Array([]int{1, 2})).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))
	query = `INSERT INTO "person_course" (person_id, course_id) VALUES (3, 4), (3, 5)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(errors.New("can't update courses"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonTransactionBeginFailure() {
	t := s.T()

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	resultPerson := models.Person{}
	resultErr := fmt.Errorf("failed to begin transaction: %w", errors.New("can't begin transaction"))

	s.dbMock.ExpectBegin().WillReturnError(errors.New("can't begin transaction"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, resultPerson, updatedPerson)
	assert.Equal(t, resultErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestUpdatePersonTransactionCommitFailure() {
	t := s.T()

	person_course := []Person_Course{
		{PersonID: 0, CourseID: 1},
		{PersonID: 0, CourseID: 2},
		{PersonID: 0, CourseID: 3},
		{PersonID: 1, CourseID: 1},
		{PersonID: 1, CourseID: 2},
		{PersonID: 1, CourseID: 3},
		{PersonID: 2, CourseID: 1},
		{PersonID: 2, CourseID: 2},
		{PersonID: 2, CourseID: 3},
		{PersonID: 3, CourseID: 1},
		{PersonID: 3, CourseID: 2},
		{PersonID: 3, CourseID: 3},
	}

	inputPerson := models.Person{ID: 3, FirstName: "Bubbly", LastName: "Thane", Type: "student", Age: 19, Courses: []int{3, 4, 5}}
	updateInput := []driver.Value{"Bubbly", "Thane", "student", 19, "Bubbles", "Thane"}
	returnPerson := models.Person{}
	returnErr := fmt.Errorf("failed to commit transaction: %w", errors.New("commit failed"))

	s.dbMock.ExpectBegin()
	query := `UPDATE "person" SET "first_name" = $1, "last_name" = $2, "type" = $3, "age" = $4 WHERE LOWER(first_name) = LOWER($5) AND LOWER(last_name) = LOWER($6)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(updateInput...).WillReturnResult(sqlmock.NewResult(3, 1))
	query = `SELECT id FROM "person" WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Bubbly", "Thane").WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 3}}))
	query = `SELECT * FROM "person_course" WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnRows(testutil.MustStructsToRows(person_course[9:]))
	query = `DELETE FROM "person_course" WHERE person_id = $1 AND course_id = ANY ($2::int[])`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(3, pq.Array([]int{1, 2})).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))
	query = `INSERT INTO "person_course" (person_id, course_id) VALUES (3, 4), (3, 5)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(1, 1))
	s.dbMock.ExpectCommit().WillReturnError(errors.New("commit failed"))

	updatedPerson, err := s.personService.UpdatePerson("Bubbles", "Thane", inputPerson)
	assert.Equal(t, returnPerson, updatedPerson)
	assert.Equal(t, returnErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Tests for CreatePerson()
// CreatePersonSuccess
// CreatePersonFailure
// CreatePersonGetCoursesFailure
// CreatePersonCourseNotFoundFailure
// CreatePersonUpdateCourseFailure
// CreatePersonTransactionBeginFailure
// CreatePersonTransactionCommitFailure
func (s *testSuit) TestCreatePersonSuccess() {
	t := s.T()

	inputPerson := models.Person{FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3, 4, 5}}
	expectedInsertedID := 4

	s.dbMock.ExpectBegin()
	query := `INSERT INTO "person" (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(inputPerson.FirstName, inputPerson.LastName, inputPerson.Type, inputPerson.Age).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedInsertedID))

	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))
	query = `INSERT INTO "person_course" (person_id, course_id) VALUES (4, 1), (4, 2), (4, 3), (4, 4), (4, 5)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(1, 5))
	s.dbMock.ExpectCommit()

	insertedID, err := s.personService.CreatePerson(inputPerson)
	assert.Equal(t, expectedInsertedID, insertedID)
	assert.NoError(t, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreatePersonFailure() {
	t := s.T()

	inputPerson := models.Person{FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3, 4, 5}}
	expectedInsertedID := -1
	expectedErr := fmt.Errorf("failed to create person: %w", errors.New("can't create person"))

	s.dbMock.ExpectBegin()
	query := `INSERT INTO "person" (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(inputPerson.FirstName, inputPerson.LastName, inputPerson.Type, inputPerson.Age).
		WillReturnRows(sqlmock.NewRows([]string{"id"})).WillReturnError(errors.New("can't create person"))

	insertedID, err := s.personService.CreatePerson(inputPerson)
	assert.Equal(t, expectedInsertedID, insertedID)
	assert.Equal(t, expectedErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreatePersonGetCoursesFailure() {
	t := s.T()

	inputPerson := models.Person{FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3, 4, 5}}
	expectedInsertedID := -1
	expectedErr := fmt.Errorf("failed to retreive course list: %w", errors.New("can't get courses"))

	s.dbMock.ExpectBegin()
	query := `INSERT INTO "person" (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(inputPerson.FirstName, inputPerson.LastName, inputPerson.Type, inputPerson.Age).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedInsertedID))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}})).WillReturnError(errors.New("can't get courses"))

	insertedID, err := s.personService.CreatePerson(inputPerson)
	assert.Equal(t, expectedInsertedID, insertedID)
	assert.Equal(t, expectedErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreatePersonCourseNotFoundFailure() {
	t := s.T()

	inputPerson := models.Person{FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{8}}
	expectedInsertedID := -1
	expectedErr := fmt.Errorf("course not found, trying to join a course that doesn't exist")

	s.dbMock.ExpectBegin()
	query := `INSERT INTO "person" (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(inputPerson.FirstName, inputPerson.LastName, inputPerson.Type, inputPerson.Age).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))

	insertedID, err := s.personService.CreatePerson(inputPerson)
	assert.Equal(t, expectedInsertedID, insertedID)
	assert.Equal(t, expectedErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreatePersonUpdateCourseFailure() {
	t := s.T()

	inputPerson := models.Person{FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3, 4, 5}}
	expectedInsertedID := -1
	expectedErr := fmt.Errorf("failed to update course list: %w", errors.New("can't update courses"))

	s.dbMock.ExpectBegin()
	query := `INSERT INTO "person" (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(inputPerson.FirstName, inputPerson.LastName, inputPerson.Type, inputPerson.Age).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))
	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))
	query = `INSERT INTO "person_course" (person_id, course_id) VALUES (4, 1), (4, 2), (4, 3), (4, 4), (4, 5)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(int64(4), 5)).WillReturnError(errors.New("can't update courses"))

	insertedID, err := s.personService.CreatePerson(inputPerson)
	assert.Equal(t, expectedInsertedID, insertedID)
	assert.Equal(t, expectedErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreatePersonTransactionBeginFailure() {
	t := s.T()

	inputPerson := models.Person{FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3, 4, 5}}
	expectedErr := fmt.Errorf("failed to begin transaction: %w", errors.New("can't begin transaction"))

	s.dbMock.ExpectBegin().WillReturnError(errors.New("can't begin transaction"))
	insertedID, err := s.personService.CreatePerson(inputPerson)
	assert.Equal(t, -1, insertedID)
	assert.Equal(t, expectedErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestCreatePersonTransactionCommitFailure() {
	t := s.T()

	inputPerson := models.Person{FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3, 4, 5}}
	expectedInsertedID := -1
	expectedErr := fmt.Errorf("failed to commit transaction: %w", errors.New("can't commit transaction"))

	s.dbMock.ExpectBegin()
	query := `INSERT INTO "person" (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(inputPerson.FirstName, inputPerson.LastName, inputPerson.Type, inputPerson.Age).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))

	query = `SELECT id FROM "course"`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(testutil.MustStructsToRows([]ID{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}))
	query = `INSERT INTO "person_course" (person_id, course_id) VALUES (4, 1), (4, 2), (4, 3), (4, 4), (4, 5)`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(1, 5))
	s.dbMock.ExpectCommit().WillReturnError(errors.New("can't commit transaction"))

	insertedID, err := s.personService.CreatePerson(inputPerson)
	assert.Equal(t, expectedInsertedID, insertedID)
	assert.Equal(t, expectedErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Tests for DeletePerson()
// DeletePersonSuccess
// DeletePersonGetIDFailure
// DeletePersonNotFoundFailure
// DeletePersonDeleteCourseFailure
// DeletePersonFailure

// DeletePersonTransactionBeginFailure
// DeletePersonTransactionCommitFailure
func (s *testSuit) TestDeletePersonSuccess() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "Thane"
	personID := 2
	expectedRowsAffected := int64(1)
	queryReturn := testutil.MustStructsToRows([]ID{{ID: personID}})

	s.dbMock.ExpectBegin()
	query := `SELECT id FROM "person" WHERE LOWER("first_name") = LOWER($1) AND LOWER("last_name") = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(queryReturn)
	query = `DELETE FROM "person_course" WHERE "person_id" = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(personID).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `DELETE FROM "person" WHERE "id" = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(personID).WillReturnResult(sqlmock.NewResult(1, 1))
	s.dbMock.ExpectCommit()

	rowsAffected, err := s.personService.DeletePerson(firstName, lastName)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, expectedRowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeletePersonGetIDFailure() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "Thane"
	personID := 2
	queryReturn := testutil.MustStructsToRows([]ID{{ID: personID}})
	expectedErr := fmt.Errorf("failed to query database for id %w", errors.New("can't get IDs"))
	expectedRowsAffected := int64(-1)

	s.dbMock.ExpectBegin()
	query := `SELECT id FROM "person" WHERE LOWER("first_name") = LOWER($1) AND LOWER("last_name") = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(queryReturn).WillReturnError(errors.New("can't get IDs"))

	rowsAffected, err := s.personService.DeletePerson(firstName, lastName)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, rowsAffected, expectedRowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeletePersonNotFoundFailure() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "Thane"
	expectedRowsAffected := int64(-1)
	queryReturn := &sqlmock.Rows{}
	expectedErr := fmt.Errorf("person not found")

	s.dbMock.ExpectBegin()
	query := `SELECT id FROM "person" WHERE LOWER("first_name") = LOWER($1) AND LOWER("last_name") = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(queryReturn)

	rowsAffected, err := s.personService.DeletePerson(firstName, lastName)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, rowsAffected, expectedRowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeletePersonDeleteCourseFailure() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "Thane"
	personID := 2
	queryReturn := testutil.MustStructsToRows([]ID{{ID: personID}})
	expectedErr := fmt.Errorf("failed to delete course relations: %w", errors.New("can't delete courses"))
	expectedRowsAffected := int64(-1)

	s.dbMock.ExpectBegin()
	query := `SELECT id FROM "person" WHERE LOWER("first_name") = LOWER($1) AND LOWER("last_name") = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(queryReturn)
	query = `DELETE FROM "person_course" WHERE "person_id" = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(personID).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(errors.New("can't delete courses"))

	rowsAffected, err := s.personService.DeletePerson(firstName, lastName)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, rowsAffected, expectedRowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeletePersonFailure() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "Thane"
	personID := 2
	queryReturn := testutil.MustStructsToRows([]ID{{ID: personID}})
	expectedErr := fmt.Errorf("failed to delete person with ID: %v. %w", personID, errors.New("can't delete person"))
	expectedRowsAffected := int64(-1)

	s.dbMock.ExpectBegin()
	query := `SELECT id FROM "person" WHERE LOWER("first_name") = LOWER($1) AND LOWER("last_name") = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(queryReturn)
	query = `DELETE FROM "person_course" WHERE "person_id" = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(personID).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `DELETE FROM "person" WHERE "id" = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(personID).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(errors.New("can't delete person"))

	rowsAffected, err := s.personService.DeletePerson(firstName, lastName)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, rowsAffected, expectedRowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeletePersonTransactionBeginFailure() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "Thane"
	expectedErr := fmt.Errorf("failed to begin transaction: %w", errors.New("can't begin transaction"))
	expectedRowsAffected := int64(-1)

	s.dbMock.ExpectBegin().WillReturnError(errors.New("can't begin transaction"))
	rowsAffected, err := s.personService.DeletePerson(firstName, lastName)

	assert.Equal(t, expectedRowsAffected, rowsAffected)
	assert.Equal(t, expectedErr, err)
	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func (s *testSuit) TestDeletePersonTransactionCommitFailure() {
	t := s.T()

	firstName := "Bubbles"
	lastName := "Thane"
	personID := 2
	expectedRowsAffected := int64(-1)
	expectedErr := fmt.Errorf("failed to commit transaction: %w", errors.New("can't commit transaction"))
	queryReturn := testutil.MustStructsToRows([]ID{{ID: personID}})

	s.dbMock.ExpectBegin()
	query := `SELECT id FROM "person" WHERE LOWER("first_name") = LOWER($1) AND LOWER("last_name") = LOWER($2) LIMIT 1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(firstName, lastName).WillReturnRows(queryReturn)
	query = `DELETE FROM "person_course" WHERE "person_id" = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(personID).WillReturnResult(sqlmock.NewResult(1, 1))
	query = `DELETE FROM "person" WHERE "id" = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(personID).WillReturnResult(sqlmock.NewResult(1, 1))
	s.dbMock.ExpectCommit().WillReturnError(errors.New("can't commit transaction"))

	rowsAffected, err := s.personService.DeletePerson(firstName, lastName)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, rowsAffected, expectedRowsAffected)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
