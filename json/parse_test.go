package json_test

import (
	"os"
	"strings"
	"testing"

	zfg "github.com/chaindead/zerocfg"
	"github.com/chaindead/zerocfg/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		awaited map[string]bool
		found   map[string]string
		unknown map[string]string
		wantErr bool
		errText string
	}{
		{
			name:  "simple key-value",
			input: `{"str": "name", "int": 1}`,
			awaited: map[string]bool{
				"str": true,
			},
			found: map[string]string{
				"str": `name`,
			},
			unknown: map[string]string{
				"int": `1`,
			},
		},
		{
			name: "nested",
			input: `{
                "host": "localhost",
                "database": {
                    "port": 5432,
                    "credentials": {
                        "username": "admin"
                    }
                }
            }`,
			awaited: map[string]bool{
				"database.credentials.username": true,
			},
			found: map[string]string{
				"database.credentials.username": `admin`,
			},
			unknown: map[string]string{
				"host":          `localhost`,
				"database.port": `5432`,
			},
		},
		{
			name:  "array",
			input: `{"tags": ["a", "b"]}`,
			awaited: map[string]bool{
				"tags": true,
			},
			found: map[string]string{
				"tags": `["a","b"]`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "map",
			input: `{"options": {"k1": 1, "k2": "v2"}}`,
			awaited: map[string]bool{
				"options": true,
			},
			found: map[string]string{
				"options": `{"k1":1,"k2":"v2"}`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "null value is skipped",
			input: `{"key1": "value1", "key2": null, "key3": 123}`,
			awaited: map[string]bool{
				"key1": true,
				"key2": true,
				"key3": true,
			},
			found: map[string]string{
				"key1": `value1`,
				"key3": `123`,
			},
			unknown: map[string]string{},
		},
		{
			name:    "json root is not an object - array",
			input:   `[1, 2, 3]`,
			wantErr: true,
			errText: "json root is not an object",
		},
		{
			name:    "json root is not an object - string",
			input:   `"just a string"`,
			wantErr: true,
			errText: "json root is not an object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.found == nil {
				tt.found = map[string]string{}
			}
			if tt.unknown == nil {
				tt.unknown = map[string]string{}
			}

			path := tempFile(t, tt.input)
			p := json.New(&path)

			found, unknown, err := p.Parse(tt.awaited, zfg.ToString)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errText != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errText), "Error message mismatch: expected to contain '%s', got '%s'", tt.errText, err.Error())
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.found, found, "Found map mismatch")
				assert.Equal(t, tt.unknown, unknown, "Unknown map mismatch")
			}
		})
	}
}

func TestParse_MalformedError(t *testing.T) {
	path := tempFile(t, `{"invalid": "json",}`)
	p := json.New(&path)

	_, _, err := p.Parse(map[string]bool{}, zfg.ToString)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal json")
}

func TestParse_EmptyInputError(t *testing.T) {
	path := tempFile(t, ``)
	p := json.New(&path)

	_, _, err := p.Parse(map[string]bool{}, zfg.ToString)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty json input")
}

func TestParse_FileNotExistError(t *testing.T) {
	nonExistentPath := "no_such_file.json"
	p := json.New(&nonExistentPath)

	_, _, err := p.Parse(map[string]bool{"some.key": true}, zfg.ToString)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read json file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func tempFile(t *testing.T, data string) string {
	f, err := os.CreateTemp("", "tmpjson-")
	require.NoError(t, err)
	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})

	_, err = f.WriteString(data)
	require.NoError(t, err)

	require.NoError(t, f.Close())

	return f.Name()
}
