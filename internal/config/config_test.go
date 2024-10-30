package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := map[string]struct {
		input        map[string]string
		output       Config
		expectsError bool
	}{
		"success": {
			input: map[string]string{
				"ENV":               "development",
				"DATABASE_NAME":     "test_db",
				"DATABASE_USER":     "test_user",
				"DATABASE_PASSWORD": "test_password",
				"DATABASE_HOST":     "localhost",
				"DATABASE_PORT":     "5432",
				"HTTP_DOMAIN":       "localhost",
				"HTTP_PORT":         "8000",
			},
			output: Config{
				Env:                  "development",
				DBName:               "test_db",
				DBUser:               "test_user",
				DBPassword:           "test_password",
				DBHost:               "localhost",
				DBPort:               "5432",
				HTTPDomain:           "localhost",
				HTTPPort:             "8000",
				HTTPShutdownDuration: 10,
			},
			expectsError: false},
		"missing required field": {
			input: map[string]string{
				"ENV":               "development",
				"DATABASE_NAME":     "test_db",
				"DATABASE_USER":     "",
				"DATABASE_PASSWORD": "test_password",
				"DATABASE_HOST":     "localhost",
				"DATABASE_PORT":     "5432",
				"HTTP_DOMAIN":       "localhost",
				"HTTP_PORT":         "8000",
			},
			output:       Config{},
			expectsError: true},
	}

	for name, testConditions := range tests {
		//setting up and cleaning up env vars for testing
		t.Run(name, func(t *testing.T) {
			for key, value := range testConditions.input {
				t.Setenv(key, value)

			}

			cfg, err := NewConfig()

			if testConditions.expectsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, cfg, testConditions.output)
		})
	}
}
