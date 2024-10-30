package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	tests := map[string]struct {
		input        string
		expectsError bool
	}{
		"successful connect": {
			input: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				"localhost",
				"courses-db-user",
				"courses-db-password",
				"coursesDB",
				"5432"),
			expectsError: false,
		},
		"unsuccessful connect": {
			input: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				"",
				"courses-db-user",
				"courses-db-password",
				"coursesDB",
				"5432"),
			expectsError: true,
		},
	}
	for testName, testConditions := range tests {
		t.Run(testName, func(t *testing.T) {
			_, err := NewDatabase(testConditions.input)
			if testConditions.expectsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
