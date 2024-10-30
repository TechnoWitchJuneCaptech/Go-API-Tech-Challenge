package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func formatName(name string) (string, string, error) {
	firstName, lastName := "", ""
	if name != "" {
		nameArr := strings.Fields(name)
		if len(nameArr) != 2 {
			return "", "", fmt.Errorf("must have a first and last name")
		}
		//capitalizing names for queries
		firstName = nameArr[0]
		firstName = strings.ToUpper(string(firstName[0])) + strings.ToLower(string(firstName[1:]))
		lastName = nameArr[1]
		lastName = strings.ToUpper(string(lastName[0])) + strings.ToLower(string(lastName[1:]))
		return firstName, lastName, nil
	}
	return "", "", fmt.Errorf("name is empty")
}
func ValidateType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "professor" || value == "student"
}
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
