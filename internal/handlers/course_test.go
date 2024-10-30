package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"tech-challenge/internal/models"
	"tech-challenge/internal/services"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCourses(t *testing.T) {

	//Test cases
	//ExistsSuccess
	//DoesNotExistSuccess
	//InternalServerError

	testCases := map[string]struct {
		serviceReturn    []models.Course
		serviceErr       error
		expectedReturn   []models.Course
		expectedHTTPCode int
	}{
		"exists success": {
			serviceReturn: []models.Course{
				{ID: 1, Name: "Class 1"},
				{ID: 2, Name: "Class 2"},
			},
			serviceErr: nil,
			expectedReturn: []models.Course{
				{ID: 1, Name: "Class 1"},
				{ID: 2, Name: "Class 2"},
			},
			expectedHTTPCode: http.StatusOK,
		},
		"does not exist success": {
			serviceReturn:    []models.Course{},
			serviceErr:       nil,
			expectedReturn:   []models.Course{},
			expectedHTTPCode: http.StatusOK,
		},
		"internal error": {
			serviceReturn:    []models.Course{},
			serviceErr:       errors.New("an error occured!"),
			expectedReturn:   []models.Course(nil),
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}
	for test, testVars := range testCases {
		log.Println("Testing " + test + "...")
		mockService := new(services.MockCourseService)
		handler := &CourseHandler{CourseService: mockService}
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/course/", nil)
		assert.NoError(t, err)

		mockService.On("GetAllCourses").Return(testVars.serviceReturn, testVars.serviceErr)
		handler.GetAllCourses(rr, req)
		var responseCourses []models.Course

		if testVars.expectedHTTPCode == http.StatusOK {
			err = json.NewDecoder(rr.Body).Decode(&responseCourses)
			assert.NoError(t, err)
			assert.Equal(t, testVars.expectedReturn, responseCourses)
		}

		assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

		mockService.AssertExpectations(t)
	}
}
func TestGetCourse(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/course/1", nil)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	mockService := new(services.MockCourseService)
	handler := &CourseHandler{CourseService: mockService}

	rr := httptest.NewRecorder()

	returnedCourse := models.Course{ID: 1, Name: "TestCourse"}
	mockService.On("GetCourse", 1).Return(returnedCourse, nil)

	fmt.Println("--- In Test ---")
	fmt.Println(req.URL)
	fmt.Println("id: " + chi.URLParam(req, "id"))
	handler.GetCourse(rr, req)

	var responseCourses models.Course
	err = json.NewDecoder(rr.Body).Decode(&responseCourses)
	assert.NoError(t, err)
	assert.Equal(t, returnedCourse, responseCourses)
	assert.Equal(t, http.StatusOK, rr.Code)

	mockService.AssertExpectations(t)
}
