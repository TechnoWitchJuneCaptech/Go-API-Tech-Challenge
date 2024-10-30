package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"tech-challenge/internal/models"
	"tech-challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func GetAllPeople(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	p := services.NewPersonService(db)
	params := r.URL.Query()
	name := params.Get("name")
	var age = -1

	if params.Has("age") {
		ageString := params.Get("age")
		var err error
		age, err = strconv.Atoi(ageString)
		if err != nil {
			logError(r, "bad request: cannot parse age to int", http.StatusBadRequest)
			http.Error(w, "bad request: cannot parse age to int", http.StatusBadRequest)
			return
		}
	}
	if age <= 0 && params.Has("age") {
		logError(r, "bad request: age must be greater than 0", http.StatusBadRequest)
		http.Error(w, "bad request: age must be greater than 0", http.StatusBadRequest)
		return
	}

	firstName, lastName := "", ""
	if name != "" {
		var err error
		firstName, lastName, err = formatName(name)
		if err != nil {
			logError(r, "bad request: "+err.Error(), http.StatusBadRequest)
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	people, err := p.GetAllPeople(age, firstName, lastName)
	if err != nil {
		logError(r, "internal error: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(people)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func GetPerson(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	p := services.NewPersonService(db)
	name := chi.URLParam(r, "name")

	if name == "" {
		logError(r, "bad request: name required", http.StatusBadRequest)
		http.Error(w, "bad request: name required", http.StatusBadRequest)
		return
	}
	firstName, lastName, err := formatName(name)
	if err != nil {
		logError(r, "bad request: "+err.Error(), http.StatusBadRequest)
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	person, err := p.GetPerson(firstName, lastName)
	if err != nil {
		logError(r, "internal error: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if reflect.DeepEqual(person, models.Person{}) {
		logError(r, "person not found", http.StatusNotFound)
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(person)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func UpdatePerson(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	p := services.NewPersonService(db)
	name := chi.URLParam(r, "name")
	if name == "" {
		logError(r, "bad request: name required", http.StatusBadRequest)
		http.Error(w, "bad request: name required", http.StatusBadRequest)
		return
	}
	firstName, lastName, err := formatName(name)
	if err != nil {
		logError(r, "bad request: "+err.Error(), http.StatusBadRequest)
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	var person models.Person
	err = json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		logError(r, err.Error(), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !areUnique(person.Courses) {
		logError(r, "bad request: class IDs must be unique", http.StatusBadRequest)
		http.Error(w, "bad request: class IDs must be unique", http.StatusBadRequest)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("ValidateType", ValidateType)
	err = validate.Struct(person)
	if err != nil {
		logError(r, "validation for person object failed", http.StatusBadRequest)
		http.Error(w, "validation for person object failed", http.StatusBadRequest)
		return
	}
	updatedPerson, err := p.UpdatePerson(firstName, lastName, person)
	if err != nil && (err.Error() == "person not found" ||
		err.Error() == "course not found, trying to join a coure that doesn't exist") {
		logError(r, err.Error(), http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		logError(r, "error updating person: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "error updating person: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(updatedPerson)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func CreatePerson(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	p := services.NewPersonService(db)
	var person models.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		logError(r, err.Error(), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !areUnique(person.Courses) {
		logError(r, "bad request: class IDs must be unique", http.StatusBadRequest)
		http.Error(w, "bad request: class IDs must be unique", http.StatusBadRequest)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("ValidateType", ValidateType)
	err = validate.Struct(person)
	if err != nil {
		logError(r, "validation for person object failed", http.StatusBadRequest)
		http.Error(w, "validation for person object failed", http.StatusBadRequest)
		return
	}
	insertedID, err := p.CreatePerson(person)
	if err != nil {
		logError(r, "failed to create person: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "failed to create person: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(insertedID)
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
func DeletePerson(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	p := services.NewPersonService(db)
	name := chi.URLParam(r, "name")

	if name == "" {
		logError(r, "bad request: name required", http.StatusBadRequest)
		http.Error(w, "bad request: name required", http.StatusBadRequest)
		return
	}
	firstName, lastName, err := formatName(name)
	if err != nil {
		logError(r, "bad request: "+err.Error(), http.StatusBadRequest)
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	deletedPersonCount, err := p.DeletePerson(firstName, lastName)
	fmt.Println(deletedPersonCount)
	fmt.Println(err)
	if deletedPersonCount == 0 || (err != nil && err.Error() == "person not found") {
		logError(r, "person not found", http.StatusNotFound)
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}
	if err != nil {
		logError(r, "could not delete person: "+err.Error(), http.StatusInternalServerError)
		http.Error(w, "could not delete person: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("person successfully deleted")
	if err != nil {
		logError(r, "internal error", http.StatusInternalServerError)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
