package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"tech-challenge/internal/models"
	"tech-challenge/internal/services"
	"testing"

	"github.com/go-chi/chi/v5"
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
		t.Run(test, func(t *testing.T) {
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
		})

	}
}
func TestGetCourse(t *testing.T) {
	testCases := map[string]struct {
		id               string
		serviceReturn    models.Course
		serviceErr       error
		expectedReturn   models.Course
		expectedHTTPCode int
	}{
		"success": {
			id:               "1",
			serviceReturn:    models.Course{ID: 1, Name: "TestCourse"},
			serviceErr:       nil,
			expectedReturn:   models.Course{ID: 1, Name: "TestCourse"},
			expectedHTTPCode: http.StatusOK,
		},
		"can't parse": {
			id:               "hi",
			serviceReturn:    models.Course{},
			serviceErr:       errors.New("Not Found"),
			expectedReturn:   models.Course{},
			expectedHTTPCode: http.StatusBadRequest,
		},
		"internal error": {
			id:               "555",
			serviceReturn:    models.Course{},
			serviceErr:       errors.New("Something went wrong"),
			expectedReturn:   models.Course{},
			expectedHTTPCode: http.StatusInternalServerError,
		},
		"course not found": {
			id:               "555",
			serviceReturn:    models.Course{},
			serviceErr:       nil,
			expectedReturn:   models.Course{},
			expectedHTTPCode: http.StatusNotFound,
		},
	}

	for test, testVars := range testCases {
		t.Run(test, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/course/"+testVars.id, nil)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", testVars.id)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			mockService := new(services.MockCourseService)
			handler := &CourseHandler{CourseService: mockService}

			rr := httptest.NewRecorder()

			intId, _ := strconv.Atoi(testVars.id)
			if test != "can't parse" {
				mockService.On("GetCourse", intId).Return(testVars.serviceReturn, testVars.serviceErr)
			}
			handler.GetCourse(rr, req)

			var responseCourses models.Course
			json.NewDecoder(rr.Body).Decode(&responseCourses)
			assert.Equal(t, testVars.expectedReturn, responseCourses)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestUpdateCourse(t *testing.T) {
	testCases := map[string]struct {
		id               string
		requestBody      models.Course
		serviceReturn    models.Course
		serviceErr       error
		expectedReturn   models.Course
		expectedHTTPCode int
	}{
		"success": {
			id:               "1",
			requestBody:      models.Course{ID: 2, Name: "UpdatedCourse"},
			serviceReturn:    models.Course{ID: 1, Name: "UpdatedCourse"},
			serviceErr:       nil,
			expectedReturn:   models.Course{ID: 1, Name: "UpdatedCourse"},
			expectedHTTPCode: http.StatusOK,
		},
		"can't parse": {
			id:               "abcd",
			requestBody:      models.Course{ID: 2, Name: "UpdatedCourse"},
			serviceReturn:    models.Course{},
			serviceErr:       nil,
			expectedReturn:   models.Course{},
			expectedHTTPCode: http.StatusBadRequest,
		},
		"bad validation": {
			id:               "1",
			requestBody:      models.Course{ID: 2},
			serviceReturn:    models.Course{},
			serviceErr:       nil,
			expectedReturn:   models.Course{},
			expectedHTTPCode: http.StatusBadRequest,
		},
		"internal error": {
			id:               "1",
			requestBody:      models.Course{ID: 2, Name: "UpdatedCourse"},
			serviceReturn:    models.Course{},
			serviceErr:       errors.New("couldn't update!"),
			expectedReturn:   models.Course{},
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for test, testVars := range testCases {
		t.Run(test, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(testVars.requestBody)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPut, "/api/course/"+testVars.id, buf)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", testVars.id)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			mockService := new(services.MockCourseService)
			handler := &CourseHandler{CourseService: mockService}

			rr := httptest.NewRecorder()

			intId, _ := strconv.Atoi(testVars.id)
			if test != "can't parse" && test != "bad validation" {
				mockService.On("UpdateCourse", intId, testVars.requestBody).Return(testVars.serviceReturn, testVars.serviceErr)
			}
			handler.UpdateCourse(rr, req)

			var responseCourses models.Course
			json.NewDecoder(rr.Body).Decode(&responseCourses)
			assert.Equal(t, testVars.expectedReturn, responseCourses)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestCreateCourse(t *testing.T) {
	testCases := map[string]struct {
		requestBody      models.Course
		serviceReturn    int
		serviceErr       error
		expectedReturn   int
		expectedHTTPCode int
	}{
		"success": {
			requestBody:      models.Course{ID: 2, Name: "CreatedCourse"},
			serviceReturn:    1,
			serviceErr:       nil,
			expectedReturn:   1,
			expectedHTTPCode: http.StatusOK,
		},
		"bad validation": {
			requestBody:      models.Course{ID: 2},
			serviceReturn:    0,
			serviceErr:       nil,
			expectedReturn:   0,
			expectedHTTPCode: http.StatusBadRequest,
		},
		"internal error": {
			requestBody:      models.Course{ID: 2, Name: "CreatedCourse"},
			serviceReturn:    -1,
			serviceErr:       errors.New("Couldn't make course!"),
			expectedReturn:   0,
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for test, testVars := range testCases {
		t.Run(test, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(testVars.requestBody)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/api/course/", buf)
			assert.NoError(t, err)

			mockService := new(services.MockCourseService)
			handler := &CourseHandler{CourseService: mockService}

			rr := httptest.NewRecorder()

			if test != "bad validation" {
				mockService.On("CreateCourse", testVars.requestBody).Return(testVars.serviceReturn, testVars.serviceErr)
			}
			handler.CreateCourse(rr, req)

			var responseBody int
			json.NewDecoder(rr.Body).Decode(&responseBody)
			assert.Equal(t, testVars.expectedReturn, responseBody)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestDeleteCourse(t *testing.T) {
	testCases := map[string]struct {
		id               string
		serviceReturn    int64
		serviceErr       error
		expectedReturn   string
		expectedHTTPCode int
	}{
		"success": {
			id:               "22",
			serviceReturn:    1,
			serviceErr:       nil,
			expectedReturn:   "course successfully deleted",
			expectedHTTPCode: http.StatusOK,
		},
		"can't parse": {
			id:               "hahahahaha",
			serviceReturn:    0,
			serviceErr:       nil,
			expectedReturn:   "",
			expectedHTTPCode: http.StatusBadRequest,
		},
		"internal error": {
			id:               "4",
			serviceReturn:    -1,
			serviceErr:       errors.New("can't delete!"),
			expectedReturn:   "",
			expectedHTTPCode: http.StatusInternalServerError,
		},
		"course not found": {
			id:               "4",
			serviceReturn:    0,
			serviceErr:       nil,
			expectedReturn:   "",
			expectedHTTPCode: http.StatusNotFound,
		},
	}

	for test, testVars := range testCases {
		t.Run(test, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "/api/course/"+testVars.id, nil)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", testVars.id)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			mockService := new(services.MockCourseService)
			handler := &CourseHandler{CourseService: mockService}

			rr := httptest.NewRecorder()

			intId, _ := strconv.Atoi(testVars.id)
			if test != "can't parse" {
				mockService.On("DeleteCourse", intId).Return(testVars.serviceReturn, testVars.serviceErr)
			}
			handler.DeleteCourse(rr, req)

			var responseCourses string
			json.NewDecoder(rr.Body).Decode(&responseCourses)
			assert.Equal(t, testVars.expectedReturn, responseCourses)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
