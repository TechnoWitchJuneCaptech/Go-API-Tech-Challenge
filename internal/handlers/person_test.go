package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"tech-challenge/internal/models"
	"tech-challenge/internal/services"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPeople(t *testing.T) {
	testCases := map[string]struct {
		age              string
		name             string
		serviceReturn    []models.Person
		serviceErr       error
		expectedReturn   []models.Person
		expectedHTTPCode int
	}{"success all": {
		age:  "",
		name: "",
		serviceReturn: []models.Person{
			{ID: 1, FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3}},
			{ID: 2, FirstName: "Jonas", LastName: "Tyroller", Type: "professor", Age: 37, Courses: []int{1, 2, 3}},
			{ID: 3, FirstName: "Blue", LastName: "Pinkman", Type: "student", Age: 18, Courses: []int{1, 2, 3}},
		},
		serviceErr: nil,
		expectedReturn: []models.Person{
			{ID: 1, FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3}},
			{ID: 2, FirstName: "Jonas", LastName: "Tyroller", Type: "professor", Age: 37, Courses: []int{1, 2, 3}},
			{ID: 3, FirstName: "Blue", LastName: "Pinkman", Type: "student", Age: 18, Courses: []int{1, 2, 3}},
		},
		expectedHTTPCode: http.StatusOK},
		"success name": {
			age:  "",
			name: "Juniper Scott",
			serviceReturn: []models.Person{
				{ID: 1, FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3}},
				{ID: 2, FirstName: "Juniper", LastName: "Scott", Type: "professor", Age: 37, Courses: []int{1, 2, 3}},
			},
			serviceErr: nil,
			expectedReturn: []models.Person{
				{ID: 1, FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 25, Courses: []int{1, 2, 3}},
				{ID: 2, FirstName: "Juniper", LastName: "Scott", Type: "professor", Age: 37, Courses: []int{1, 2, 3}},
			},
			expectedHTTPCode: http.StatusOK},
		"success name and age": {
			age:  "28",
			name: "Juniper Scott",
			serviceReturn: []models.Person{
				{ID: 1, FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 28, Courses: []int{1, 2, 3}},
				{ID: 2, FirstName: "Juniper", LastName: "Scott", Type: "professor", Age: 28, Courses: []int{1, 2, 3}},
			},
			serviceErr: nil,
			expectedReturn: []models.Person{
				{ID: 1, FirstName: "Juniper", LastName: "Scott", Type: "student", Age: 28, Courses: []int{1, 2, 3}},
				{ID: 2, FirstName: "Juniper", LastName: "Scott", Type: "professor", Age: 28, Courses: []int{1, 2, 3}},
			},
			expectedHTTPCode: http.StatusOK},
		"success age": {
			age:  "28",
			name: "",
			serviceReturn: []models.Person{
				{ID: 1, FirstName: "Funman", LastName: "McGee", Type: "student", Age: 28, Courses: []int{1, 2, 3}},
				{ID: 2, FirstName: "Juniper", LastName: "Scott", Type: "professor", Age: 28, Courses: []int{1, 2, 3}},
			},
			serviceErr: nil,
			expectedReturn: []models.Person{
				{ID: 1, FirstName: "Funman", LastName: "McGee", Type: "student", Age: 28, Courses: []int{1, 2, 3}},
				{ID: 2, FirstName: "Juniper", LastName: "Scott", Type: "professor", Age: 28, Courses: []int{1, 2, 3}},
			},
			expectedHTTPCode: http.StatusOK},
		"success empty": {
			age:              "28",
			name:             "My FavoritePerson",
			serviceReturn:    []models.Person{},
			serviceErr:       nil,
			expectedReturn:   []models.Person{},
			expectedHTTPCode: http.StatusOK},
		"failure can't parse": {
			age:              "abcd",
			name:             "My FavoritePerson",
			serviceReturn:    []models.Person{},
			serviceErr:       nil,
			expectedReturn:   []models.Person{},
			expectedHTTPCode: http.StatusBadRequest},
		"failure negative age": {
			age:              "-25",
			name:             "My FavoritePerson",
			serviceReturn:    []models.Person{},
			serviceErr:       nil,
			expectedReturn:   []models.Person{},
			expectedHTTPCode: http.StatusBadRequest},
		"failure format one name": {
			age:              "25",
			name:             "MyFavoritePerson",
			serviceReturn:    []models.Person{},
			serviceErr:       nil,
			expectedReturn:   []models.Person{},
			expectedHTTPCode: http.StatusBadRequest},
		"failure format three names": {
			age:              "25",
			name:             "My Favorite Person",
			serviceReturn:    []models.Person{},
			serviceErr:       nil,
			expectedReturn:   []models.Person{},
			expectedHTTPCode: http.StatusBadRequest},
		"failure internal error": {
			age:              "25",
			name:             "My FavoritePerson",
			serviceReturn:    []models.Person{},
			serviceErr:       errors.New("couldn't get people!"),
			expectedReturn:   []models.Person{},
			expectedHTTPCode: http.StatusInternalServerError},
	}

	for test, testVars := range testCases {
		t.Run(test, func(t *testing.T) {
			mockService := new(services.MockPersonService)
			handler := &PersonHandler{PersonService: mockService}
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/person/", nil)
			assert.NoError(t, err)

			if test == "success name" {
				q := req.URL.Query()
				q.Add("name", testVars.name)
				req.URL.RawQuery = q.Encode()
			} else if test == "success name and age" || test == "success empty" || test == "failure negative age" || test == "failure format one name" || test == "failure format three names" || test == "failure internal error" {
				q := req.URL.Query()
				q.Add("name", testVars.name)
				q.Add("age", testVars.age)
				req.URL.RawQuery = q.Encode()
			} else if test == "success age" || test == "failure can't parse" {
				q := req.URL.Query()
				q.Add("age", testVars.age)
				req.URL.RawQuery = q.Encode()
			}
			if test != "failure can't parse" && test != "failure format one name" && test != "failure format three names" {
				firstName, lastName := "", ""
				if test != "success all" && test != "success age" {
					firstName, lastName, err = formatName(testVars.name)
					assert.NoError(t, err)
				}
				ageInt := -1
				if testVars.age != "" {
					ageInt, err = strconv.Atoi(testVars.age)
					assert.NoError(t, err)
				}
				if test != "failure negative age" {
					mockService.On("GetAllPeople", ageInt, firstName, lastName).Return(testVars.serviceReturn, testVars.serviceErr)
				}
			}
			handler.GetAllPeople(rr, req)
			var responsePeople []models.Person

			if testVars.expectedHTTPCode == http.StatusOK {
				err = json.NewDecoder(rr.Body).Decode(&responsePeople)
				assert.NoError(t, err)
				assert.Equal(t, testVars.expectedReturn, responsePeople)
			}

			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestGetPerson(t *testing.T) {
	testCases := map[string]struct {
		name             string
		queryFirstName   string
		queryLastName    string
		serviceReturn    models.Person
		serviceErr       error
		expectedReturn   models.Person
		expectedHTTPCode int
	}{"success": {
		name:             "My Favoriteperson",
		queryFirstName:   "My",
		queryLastName:    "Favoriteperson",
		serviceReturn:    models.Person{ID: 25, Age: 27, FirstName: "My", LastName: "Favoriteperson", Type: "professor", Courses: []int{1}},
		serviceErr:       nil,
		expectedReturn:   models.Person{ID: 25, Age: 27, FirstName: "My", LastName: "Favoriteperson", Type: "professor", Courses: []int{1}},
		expectedHTTPCode: http.StatusOK,
	}, "success case insensitive": {
		name:             "my FaVoRiTePeRsOn",
		queryFirstName:   "my",
		queryLastName:    "FaVoRiTePeRsOn",
		serviceReturn:    models.Person{ID: 25, Age: 27, FirstName: "My", LastName: "FavoritePerson", Type: "professor", Courses: []int{1}},
		serviceErr:       nil,
		expectedReturn:   models.Person{ID: 25, Age: 27, FirstName: "My", LastName: "FavoritePerson", Type: "professor", Courses: []int{1}},
		expectedHTTPCode: http.StatusOK,
	}, "failure format one name": {
		name:             "MYFAVORITEPERSON",
		queryFirstName:   "",
		queryLastName:    "",
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure format three names": {
		name:             "MY FAVORITE PERSON",
		queryFirstName:   "",
		queryLastName:    "",
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing name": {
		name:             "",
		queryFirstName:   "",
		queryLastName:    "",
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure not found": {
		name:             "My Favoriteperson",
		queryFirstName:   "My",
		queryLastName:    "Favoriteperson",
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusNotFound,
	}, "failure internal server error": {
		name:             "My Favoriteperson",
		queryFirstName:   "My",
		queryLastName:    "Favoriteperson",
		serviceReturn:    models.Person{},
		serviceErr:       errors.New("not found!"),
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusInternalServerError,
	}}
	for testName, testVars := range testCases {
		t.Run(testName, func(t *testing.T) {
			mockService := new(services.MockPersonService)
			handler := &PersonHandler{PersonService: mockService}
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/person/"+testVars.name, nil)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", testVars.name)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			firstName, lastName := "", ""
			firstName, lastName, err = formatName(testVars.name)
			if testName != "failure format one name" && testName != "failure format three names" && testName != "failure missing name" {
				assert.NoError(t, err)
				assert.Equal(t, testVars.queryFirstName, firstName)
				assert.Equal(t, testVars.queryLastName, lastName)

				mockService.On("GetPerson", firstName, lastName).Return(testVars.serviceReturn, testVars.serviceErr)

			} else {
				assert.Error(t, err)
			}

			handler.GetPerson(rr, req)
			var responsePerson models.Person

			if testName != "failure format one name" && testName != "failure format three names" && testName != "failure missing name" && testName != "failure not found" && testName != "failure internal server error" {
				err = json.NewDecoder(rr.Body).Decode(&responsePerson)
				assert.NoError(t, err)
			}

			assert.Equal(t, testVars.expectedReturn, responsePerson)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestUpdatePerson(t *testing.T) {
	testCases := map[string]struct {
		name             string
		requestBody      models.Person
		serviceReturn    models.Person
		serviceErr       error
		expectedReturn   models.Person
		expectedHTTPCode int
	}{"success": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "NewName", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    models.Person{ID: 25, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		serviceErr:       nil,
		expectedReturn:   models.Person{ID: 25, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		expectedHTTPCode: http.StatusOK,
	}, "success uppercase": {
		name:             "MY FAVORITEPERSON",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "MY", LastName: "NEWNAME", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    models.Person{ID: 25, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		serviceErr:       nil,
		expectedReturn:   models.Person{ID: 25, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		expectedHTTPCode: http.StatusOK,
	}, "success lowercase": {
		name:             "my favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "my", LastName: "newname", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    models.Person{ID: 25, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		serviceErr:       nil,
		expectedReturn:   models.Person{ID: 25, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		expectedHTTPCode: http.StatusOK,
	}, "failure missing name": {
		name:             "",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "my", LastName: "newname", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing new name": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "", LastName: "", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure name one word": {
		name:             "MyFavoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure name three words": {
		name:             "My Favorite Person",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure duplicate courses": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 5, 5}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing age": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{5, 5, 5}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing firstName": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, LastName: "Newname", Type: "professor", Courses: []int{5, 5, 5}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing lastName": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", Type: "professor", Courses: []int{5, 5, 5}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing type": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Courses: []int{5, 5, 5}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing courses": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor"},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure invalid type": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "madman", Courses: []int{1, 2, 3}},
		serviceReturn:    models.Person{},
		serviceErr:       nil,
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure person not found": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{1, 2, 3}},
		serviceReturn:    models.Person{},
		serviceErr:       fmt.Errorf("person not found"),
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusNotFound,
	}, "failure course not found": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{88888}},
		serviceReturn:    models.Person{},
		serviceErr:       fmt.Errorf("course not found, trying to join a course that doesn't exist"),
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusNotFound,
	}, "failure internal error": {
		name:             "My Favoriteperson",
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "Newname", Type: "professor", Courses: []int{88888}},
		serviceReturn:    models.Person{},
		serviceErr:       fmt.Errorf("new error!"),
		expectedReturn:   models.Person{},
		expectedHTTPCode: http.StatusInternalServerError,
	}}
	for testName, testVars := range testCases {
		t.Run(testName, func(t *testing.T) {
			mockService := new(services.MockPersonService)
			handler := &PersonHandler{PersonService: mockService}

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(testVars.requestBody)
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/api/person/"+testVars.name, buf)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", testVars.name)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			firstName, lastName := "", ""
			firstName, lastName, err = formatName(testVars.name)
			if testName == "failure missing name" || testName == "failure name one word" || testName == "failure name three words" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if testName == "success" || testName == "success uppercase" || testName == "success lowercase" || testName == "failure person not found" || testName == "failure course not found" || testName == "failure internal error" {
				mockService.On("UpdatePerson", firstName, lastName, testVars.requestBody).Return(testVars.serviceReturn, testVars.serviceErr)
			}

			handler.UpdatePerson(rr, req)
			var responsePerson models.Person

			if testName == "success" || testName == "success uppercase" || testName == "success lowercase" {
				err = json.NewDecoder(rr.Body).Decode(&responsePerson)
				assert.NoError(t, err)
			}

			assert.Equal(t, testVars.expectedReturn, responsePerson)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestCreatePerson(t *testing.T) {
	testCases := map[string]struct {
		requestBody      models.Person
		serviceReturn    int
		serviceErr       error
		expectedReturn   int
		expectedHTTPCode int
	}{"success": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "NewName", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    12,
		serviceErr:       nil,
		expectedReturn:   12,
		expectedHTTPCode: http.StatusOK,
	}, "failure missing body": {
		requestBody:      models.Person{},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing age": {
		requestBody:      models.Person{ID: 28, FirstName: "My", LastName: "NewName", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure negative age": {
		requestBody:      models.Person{ID: 28, Age: -1, FirstName: "My", LastName: "NewName", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure zero age": {
		requestBody:      models.Person{ID: 28, Age: 0, FirstName: "My", LastName: "NewName", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing firstName": {
		requestBody:      models.Person{ID: 28, Age: 28, LastName: "NewName", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure empty firstName": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "", LastName: "NewName", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing lastName": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure empty lastName": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "", Type: "professor", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing type": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "NewName", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure incompatible type": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "NewName", Type: "madman", Courses: []int{5, 6, 7}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure missing courses": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "NewName", Type: "professor"},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure duplicate courses": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "NewName", Type: "professor", Courses: []int{5, 4, 3, 2, 5}},
		serviceReturn:    -1,
		serviceErr:       nil,
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusBadRequest,
	}, "failure internal error": {
		requestBody:      models.Person{ID: 28, Age: 28, FirstName: "My", LastName: "NewName", Type: "professor", Courses: []int{5, 4, 3, 2, 1}},
		serviceReturn:    -1,
		serviceErr:       errors.New("an error occured!"),
		expectedReturn:   -1,
		expectedHTTPCode: http.StatusInternalServerError,
	}}

	for testName, testVars := range testCases {
		t.Run(testName, func(t *testing.T) {
			mockService := new(services.MockPersonService)
			handler := &PersonHandler{PersonService: mockService}

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(testVars.requestBody)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/api/person/", buf)
			assert.NoError(t, err)

			if testName == "success" || testName == "failure internal error" {
				mockService.On("CreatePerson", testVars.requestBody).Return(testVars.serviceReturn, testVars.serviceErr)
			}
			handler.CreatePerson(rr, req)

			var response int = -1
			if testName == "success" {
				err = json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
			}

			assert.Equal(t, testVars.expectedReturn, response)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
func TestDeletePerson(t *testing.T) {
	testCases := map[string]struct {
		name             string
		serviceReturn    int64
		serviceErr       error
		expectedReturn   string
		expectedHTTPCode int
	}{
		"success": {
			name:             "My Testperson",
			serviceReturn:    1,
			serviceErr:       nil,
			expectedReturn:   "person successfully deleted",
			expectedHTTPCode: http.StatusOK,
		},
		"failure missing name": {
			name:             "",
			serviceReturn:    -1,
			serviceErr:       nil,
			expectedReturn:   "",
			expectedHTTPCode: http.StatusBadRequest,
		},
		"failure one name": {
			name:             "Johnny",
			serviceReturn:    -1,
			serviceErr:       nil,
			expectedReturn:   "",
			expectedHTTPCode: http.StatusBadRequest,
		},
		"failure three names": {
			name:             "Johnny Silver Bullet",
			serviceReturn:    -1,
			serviceErr:       nil,
			expectedReturn:   "",
			expectedHTTPCode: http.StatusBadRequest,
		},
		"failure not found 1": {
			name:             "Johnny Bullet",
			serviceReturn:    0,
			serviceErr:       nil,
			expectedReturn:   "",
			expectedHTTPCode: http.StatusNotFound,
		},
		"failure not found 2": {
			name:             "Johnny Bullet",
			serviceReturn:    -1,
			serviceErr:       fmt.Errorf("person not found"),
			expectedReturn:   "",
			expectedHTTPCode: http.StatusNotFound,
		},
		"failure internal error": {
			name:             "Johnny Bullet",
			serviceReturn:    -1,
			serviceErr:       fmt.Errorf("new error!"),
			expectedReturn:   "",
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for test, testVars := range testCases {
		t.Run(test, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "/api/person/"+testVars.name, nil)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", testVars.name)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			mockService := new(services.MockPersonService)
			handler := &PersonHandler{PersonService: mockService}

			rr := httptest.NewRecorder()

			firstName, lastName, err := formatName(testVars.name)
			if test == "success" || test == "failure not found 1" || test == "failure not found 2" || test == "failure internal error" {
				assert.NoError(t, err)
				mockService.On("DeletePerson", firstName, lastName).Return(testVars.serviceReturn, testVars.serviceErr)
			} else {
				assert.Error(t, err)
			}
			handler.DeletePerson(rr, req)

			var responseCourses string
			err = json.NewDecoder(rr.Body).Decode(&responseCourses)
			if test == "success" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, testVars.expectedReturn, responseCourses)
			assert.Equal(t, testVars.expectedHTTPCode, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
