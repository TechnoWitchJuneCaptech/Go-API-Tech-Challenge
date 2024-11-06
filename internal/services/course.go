package services

//course.go defines the service functions and logic used by RealCourseService structs to query a db, and a CourseService interface for testing.

import (
	"database/sql"
	"fmt"
	"tech-challenge/internal/models"
)

type CourseService interface {
	GetAllCourses() ([]models.Course, error)
	GetCourse(int) (models.Course, error)
	UpdateCourse(int, models.Course) (models.Course, error)
	CreateCourse(models.Course) (int, error)
	DeleteCourse(int) (int64, error)
}

type RealCourseService struct {
	db *sql.DB
}

func NewCourseService(db *sql.DB) *RealCourseService {
	return &RealCourseService{
		db: db,
	}
}

func (c *RealCourseService) GetAllCourses() ([]models.Course, error) {
	rows, err := c.db.Query(`SELECT * FROM "course"`)
	if err != nil {
		return []models.Course{}, fmt.Errorf("failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err = rows.Scan(&course.ID, &course.Name)
		if err != nil {
			return []models.Course{}, fmt.Errorf("failed to scan course from row: %w", err)
		}
		courses = append(courses, course)
	}
	if err = rows.Err(); err != nil {
		return []models.Course{}, fmt.Errorf("failed to scan courses: %w", err)
	}
	return courses, nil
}
func (c *RealCourseService) GetCourse(id int) (models.Course, error) {
	row, err := c.db.Query(`SELECT * FROM "course" 
							WHERE "id" = $1 
							LIMIT 1`, id)
	if err != nil {
		return models.Course{}, fmt.Errorf("failed to get course: %w", err)
	}
	defer row.Close()

	var course models.Course
	if isEmpty := !row.Next(); isEmpty {
		return models.Course{}, fmt.Errorf("course not found")
	}
	err = row.Scan(&course.ID, &course.Name)
	if err != nil {
		return models.Course{}, fmt.Errorf("failed to scan course from row: %w", err)
	}
	return course, nil
}
func (c *RealCourseService) UpdateCourse(id int, course models.Course) (models.Course, error) {
	row, err := c.db.Exec(`UPDATE "course" 
						SET "name" = $1
						WHERE "id" = $2`,
		course.Name,
		id,
	)
	if err != nil {
		return models.Course{}, fmt.Errorf("failed to update course: %w", err)
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return models.Course{}, fmt.Errorf("failed to update course: %w", err)
	}
	if rowsAffected == 0 {
		return models.Course{}, fmt.Errorf("course not found")
	}
	course.ID = id
	return course, nil
}
func (c *RealCourseService) CreateCourse(course models.Course) (int, error) {
	row, err := c.db.Query(`INSERT INTO "course" (name)
							VALUES ($1) RETURNING id`,
		course.Name)
	if err != nil {
		return -1, fmt.Errorf("failed to create course: %w", err)
	}
	var lastInsertedID = -1
	row.Next()
	err = row.Scan(&lastInsertedID)
	if err != nil {
		return -1, fmt.Errorf("internal error accessing inserted id: %w", err)
	}
	return lastInsertedID, nil
}
func (c *RealCourseService) DeleteCourse(id int) (int64, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	rows, err := tx.Exec(`DELETE FROM "person_course"
						WHERE "course_id" = $1`,
		id)
	if err != nil {
		return -1, fmt.Errorf("failed to delete course relations: %w", err)
	}
	rows, err = tx.Exec(`DELETE FROM "course"
						WHERE "id" = $1`,
		id)
	if err != nil {
		return -1, fmt.Errorf("failed to delete course with ID: %v. %w", id, err)
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("failed to get the number of affected rows: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return rowsAffected, nil
}
