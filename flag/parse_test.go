package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected map[string]string
	}{
		{
			name: "single short flag with value",
			args: []string{"-name", "value"},
			expected: map[string]string{
				"name": "value",
			},
		},
		{
			name: "single long flag with value",
			args: []string{"--name", "value"},
			expected: map[string]string{
				"name": "value",
			},
		},
		{
			name: "multiple flags with values",
			args: []string{"-host", "localhost", "--port", "5432", "-user", "admin"},
			expected: map[string]string{
				"host": "localhost",
				"port": "5432",
				"user": "admin",
			},
		},
		{
			name: "flag without value",
			args: []string{"-debug"},
			expected: map[string]string{
				"debug": "",
			},
		},
		{
			name: "mixed flags with and without values",
			args: []string{"-verbose", "--config", "config.yaml", "-force"},
			expected: map[string]string{
				"verbose": "",
				"config":  "config.yaml",
				"force":   "",
			},
		},
		{
			name: "ignore non-flag arguments",
			args: []string{"positional", "-name", "value", "another", "-debug"},
			expected: map[string]string{
				"name":  "value",
				"debug": "",
			},
		},
		{
			name:     "empty args",
			args:     []string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parse(tt.args)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, result)
		})
	}
}
