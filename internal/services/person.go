package services

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"tech-challenge/internal/models"

	"github.com/lib/pq"
)

type PersonService struct {
	db *sql.DB
}

func NewPersonService(db *sql.DB) *PersonService {
	return &PersonService{
		db: db,
	}
}
func (p PersonService) GetAllPeople(age int, firstName string, lastName string) ([]models.Person, error) {
	var err error

	var rows *sql.Rows
	if age != -1 && (firstName != "" && lastName != "") {
		rows, err = dbQueryGetPeopleByNameAndAge(firstName, lastName, age, p.db)
	} else if age != -1 {
		rows, err = dbQueryGetPeopleByAge(age, p.db)
	} else if firstName != "" && lastName != "" {
		rows, err = dbQueryGetPeopleByName(firstName, lastName, p.db)
	} else {
		rows, err = dbQueryGetPeople(p.db)
	}

	if err != nil {
		return []models.Person{}, fmt.Errorf("failed to get people: %w", err)
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var person models.Person
		err = rows.Scan(&person.ID,
			&person.FirstName,
			&person.LastName,
			&person.Type,
			&person.Age,
		)
		if err != nil {
			return []models.Person{}, fmt.Errorf("failed to scan person from row: %w", err)
		}

		courseRows, err := p.db.Query(`SELECT * FROM "person_course"
				WHERE person_id = $1`,
			person.ID)
		if err != nil {
			return []models.Person{}, fmt.Errorf("failed to get courses for person: %w", err)
		}
		var personCourses = make([]int, 0)
		for courseRows.Next() {
			var personID int
			var courseID int
			courseRows.Scan(&personID, &courseID)
			personCourses = append(personCourses, courseID)
		}
		person.Courses = personCourses
		people = append(people, person)
	}
	if err = rows.Err(); err != nil {
		return []models.Person{}, fmt.Errorf("failed to scan people: %w", err)
	}
	return people, nil
}
func (p PersonService) GetPerson(firstName string, lastName string) (models.Person, error) {
	rows, err := p.db.Query(`SELECT * FROM "person" 
	WHERE first_name = $1
	AND last_name = $2
	LIMIT 1`,
		firstName,
		lastName)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to get person: %w", err)
	}

	var person models.Person
	if isEmpty := !rows.Next(); isEmpty {
		return person, nil
	}

	err = rows.Scan(&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.Type,
		&person.Age,
	)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to scan person: %w", err)
	}
	courseRows, err := p.db.Query(`SELECT * FROM "person_course"
				WHERE person_id = $1`,
		person.ID)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to get courses for person: %w", err)
	}
	var personCourses = make([]int, 0)
	for courseRows.Next() {
		var personID int
		var courseID int
		courseRows.Scan(&personID, &courseID)
		personCourses = append(personCourses, courseID)
	}
	person.Courses = personCourses

	return person, nil
}

// This is really bad architecture. Because firstName and lastName do not constitute a unique key, this function could update the wrong user.
func (p PersonService) UpdatePerson(firstName string, lastName string, person models.Person) (models.Person, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row, err := tx.Exec(`UPDATE "person" 
	SET "first_name" = $1,
		"last_name" = $2,
		"type" = $3,
		"age" = $4
	WHERE "first_name" = $5
		AND "last_name" = $6`,
		person.FirstName,
		person.LastName,
		person.Type,
		person.Age,
		firstName,
		lastName,
	)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to update person: %w", err)
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to update person: %w", err)
	}
	if rowsAffected == 0 {
		return models.Person{}, fmt.Errorf("person not found")
	}

	//removing and adding courses to person_course

	//1. do a select to get ID
	rows, err := tx.Query(`SELECT id FROM "person"
						WHERE first_name = $1
						AND last_name = $2
						LIMIT 1`,
		person.FirstName,
		person.LastName)
	if err != nil || !rows.Next() {
		return models.Person{}, fmt.Errorf("failed to retreive id: %w", err)
	}
	rows.Scan(&person.ID)
	rows.Close()
	//2. use ID to do select of courses from person_course
	rows, err = tx.Query(`SELECT * FROM "person_course"
						WHERE person_id = $1`,
		person.ID)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to retreive course list: %w", err)
	}
	//3. form an array of all classes they're currently in
	currentCourses := make([]int, 0)
	for rows.Next() {
		var pID int
		var cID int
		rows.Scan(&pID, &cID)
		currentCourses = append(currentCourses, cID)
	}
	rows.Close()
	//4. do a delete query on the ones not in the new person course list
	coursesToDelete := getDifference(currentCourses, person.Courses)
	if len(coursesToDelete) > 0 {
		_, err = tx.Exec(`DELETE FROM "person_course" 
		WHERE person_id = $1 
		AND course_id = ANY ($2::int[])`,
			person.ID,
			pq.Array(coursesToDelete))

		if err != nil {
			return models.Person{}, fmt.Errorf("failed to update course list: %w", err)
		}
	}
	//5. Validate the courses they want to be added to actually exist
	coursesToInsert := getDifference(person.Courses, currentCourses)

	rows, err = tx.Query(`SELECT id FROM "course"`)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to retreive course list: %w", err)
	}

	courseIDs := make(map[int]bool)
	for rows.Next() {
		var courseID int
		rows.Scan(&courseID)
		courseIDs[courseID] = true
	}
	for _, val := range coursesToInsert {
		if !courseIDs[val] {
			return models.Person{}, fmt.Errorf("course not found, trying to join a course that doesn't exist")
		}
	}
	rows.Close()
	//6. do an insert query on the ones not currently in the table
	var sb strings.Builder
	if len(coursesToInsert) > 0 {
		sb.WriteString("(" + strconv.Itoa(person.ID) + ", " + strconv.Itoa(coursesToInsert[0]) + ")")
		for i := 1; i < len(coursesToInsert); i++ {
			sb.WriteString(", (" + strconv.Itoa(person.ID) + ", " + strconv.Itoa(coursesToInsert[i]) + ")")
		}
	}
	if sb.String() != "" {
		query := `INSERT INTO "person_course" (person_id, course_id) VALUES ` + sb.String()
		_, err = tx.Exec(query)
		if err != nil {
			return models.Person{}, fmt.Errorf("failed to update course list: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return models.Person{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return person, nil
}
func (p PersonService) CreatePerson(person models.Person) (int, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	//insert person into table
	row, err := tx.Query(`INSERT INTO "person" (first_name, last_name, type, age)
							VALUES ($1, $2, $3, $4) RETURNING id`,
		person.FirstName,
		person.LastName,
		person.Type,
		person.Age)
	if err != nil {
		return -1, fmt.Errorf("failed to create person: %w", err)
	}
	var lastInsertedID = -1
	row.Next()
	err = row.Scan(&lastInsertedID)
	if err != nil {
		return -1, fmt.Errorf("internal error accessing inserted id: %w", err)
	}
	row.Close()
	//validate all courses to insert exist
	rows, err := tx.Query(`SELECT id FROM "course"`)
	if err != nil {
		return -1, fmt.Errorf("failed to retreive course list: %w", err)
	}

	courseIDs := make(map[int]bool)
	for rows.Next() {
		var courseID int
		rows.Scan(&courseID)
		courseIDs[courseID] = true
	}
	rows.Close()
	for _, val := range person.Courses {
		if !courseIDs[val] {
			return -1, fmt.Errorf("course not found, trying to join a course that doesn't exist")
		}
	}
	//inserting to person_courses. For this iteration, we assume the person is not
	//currently in any courses since we are creating them. This should be error
	//checked on a future iteration.
	var sb strings.Builder
	if len(person.Courses) > 0 {
		sb.WriteString("(" + strconv.Itoa(lastInsertedID) + ", " + strconv.Itoa(person.Courses[0]) + ")")
		for i := 1; i < len(person.Courses); i++ {
			sb.WriteString(", (" + strconv.Itoa(lastInsertedID) + ", " + strconv.Itoa(person.Courses[i]) + ")")
		}
	}
	if sb.String() != "" {
		query := `INSERT INTO "person_course" (person_id, course_id) VALUES ` + sb.String()
		_, err = tx.Exec(query)
		if err != nil {
			return -1, fmt.Errorf("failed to update course list: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return lastInsertedID, nil
}

// This is really bad architecture. Because firstName and lastName do not constitute a unique key, this function could delete multiple users.
func (p PersonService) DeletePerson(firstName string, lastName string) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	//get person's id
	rows, err := tx.Query(`SELECT id FROM "person"
						WHERE "first_name" = $1
						AND "last_name" = $2
						LIMIT 1`, firstName, lastName)

	if err != nil {
		return -1, fmt.Errorf("failed to query database for id %w", err)
	}
	var personID int
	if !rows.Next() {
		return -1, fmt.Errorf("person not found")
	}
	rows.Scan(&personID)
	rows.Close()

	//delete from person_course
	//This version assumes first_name + last_name can be used as a unique identifier.
	//This will cause errors. If two people have the same name, it's possible
	//we will delete the wrong one.
	//In the future, this API should change to using id since it is the table's primary key or have another way to uniquely identify person entities.

	_, err = tx.Exec(`DELETE FROM "person_course"
						WHERE "person_id" = $1`,
		personID)
	if err != nil {
		return -1, fmt.Errorf("failed to delete course relations: %w", err)
	}
	//delete from person

	result, err := tx.Exec(`DELETE FROM "person"
						WHERE "id" = $1`,
		personID)
	if err != nil {
		return -1, fmt.Errorf("failed to delete person with ID: %v. %w", personID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("failed to get the number of affected rows: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return rowsAffected, nil
}
