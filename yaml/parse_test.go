package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name: "simple key-value",
			input: `
str: name
int: 1`,
			expected: map[string]string{
				"str": "name",
				"int": "1",
			},
		},
		{
			name: "nested structure",
			input: `
host: localhost
database:
  port: 5432
  credentials:
    username: admin`,
			expected: map[string]string{
				"host":                          "localhost",
				"database.port":                 "5432",
				"database.credentials.username": "admin",
			},
		},
		{
			name: "array values",
			input: `
tags:
  - tag1
  - tag2
  - tag3
settings:
  numbers: [1, 2, 3]`,
			expected: map[string]string{
				"tags":             "tag1,tag2,tag3",
				"settings.numbers": "1,2,3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parse([]byte(tt.input))
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParse_Error(t *testing.T) {
	_, err := parse([]byte(`invalid: [yaml: "missing closing quote`))
	assert.Error(t, err)
}
