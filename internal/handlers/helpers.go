package handlers

//helpers.go defines miscellaneous helper functions used in ../handlers/* and ../services/*

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// turns name into firstName and lastName for querying. Will return an error if there's not two names or name is empty.
func formatName(name string) (firstName string, lastName string, err error) {
	firstName, lastName = "", ""
	if name != "" {
		nameArr := strings.Fields(name)
		if len(nameArr) != 2 {
			return "", "", fmt.Errorf("must have a first and last name")
		}
		firstName = nameArr[0]
		lastName = nameArr[1]
		return firstName, lastName, nil
	}
	return "", "", fmt.Errorf("name is empty")
}

// custom validation function for course object
func ValidateType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "professor" || value == "student"
}

// returns whether all values in input are unique
func areUnique(input []int) bool {
	result := make(map[int]int)
	for _, val := range input {
		result[val]++
		if result[val] != 1 {
			return false
		}
	}
	return true
}
func logError(r *http.Request, message string, status int) {
	log.Printf(strconv.Itoa(status) + " ERROR: " + message + " at: " + r.Method + " " + r.URL.Path)
}
