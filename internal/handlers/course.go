package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"tech-challenge/internal/models"
	"tech-challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CourseHandler struct {
	CourseService services.CourseService
}

func (c *CourseHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := c.CourseService.GetAllCourses()
	if err != nil {
		logError(r, "internal error: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(courses)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func (c *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		logError(r, "bad request: cannot parse id to int", http.StatusBadRequest)
		http.Error(w, "bad request: cannot parse id to int", http.StatusBadRequest)
		return
	}
	course, err := c.CourseService.GetCourse(idInt)
	if err != nil {
		logError(r, "internal error: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if reflect.DeepEqual(course, models.Course{}) {
		logError(r, "course not found", http.StatusNotFound)
		http.Error(w, "course not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(course)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func (c *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		logError(r, "bad request: cannot parse id to int", http.StatusBadRequest)
		http.Error(w, "bad request: cannot parse id to int", http.StatusBadRequest)
		return
	}
	var course models.Course
	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		logError(r, err.Error(), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(course)
	if err != nil {
		logError(r, "validation for course object failed", http.StatusBadRequest)
		http.Error(w, "validation for course object failed", http.StatusBadRequest)
		return
	}
	updatedCourse, err := c.CourseService.UpdateCourse(idInt, course)
	if err != nil && err.Error() == "course not found" {
		logError(r, err.Error(), http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		logError(r, "error updating course: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "error updating course: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(updatedCourse)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func (c *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		logError(r, err.Error(), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(course)
	if err != nil {
		logError(r, "validation for course object failed: "+err.Error(), http.StatusBadRequest)
		http.Error(w, "validation for course object failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	insertedID, err := c.CourseService.CreateCourse(course)
	if err != nil {
		logError(r, "failed to create course: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "failed to create course: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(insertedID)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func (c *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		logError(r, "bad request: cannot parse id to int", http.StatusBadRequest)
		http.Error(w, "bad request: cannot parse id to int", http.StatusBadRequest)
		return
	}
	deletedCourseCount, err := c.CourseService.DeleteCourse(idInt)
	if err != nil {
		logError(r, "could not delete course: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "could not delete course: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if deletedCourseCount == 0 {
		logError(r, "course not found", http.StatusNotFound)
		http.Error(w, "course not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode("course successfully deleted")
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
