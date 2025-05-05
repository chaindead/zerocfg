package json_test

import (
	"os"
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
			input: `{"tags": ["tag1", "tag2", "tag3"]}`,
			awaited: map[string]bool{
				"tags": true,
			},
			found: map[string]string{
				"tags": `["tag1","tag2","tag3"]`,
			},
		},
		{
			name:  "map within awaited key",
			input: `{"db": {"host": "localhost", "port": 5432}}`,
			awaited: map[string]bool{
				"db": true,
			},
			found: map[string]string{
				"db": `{"host":"localhost","port":5432}`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "awaited key within map",
			input: `{"db": {"host": "localhost", "port": 5432}}`,
			awaited: map[string]bool{
				"db.host": true,
			},
			found: map[string]string{
				"db.host": `localhost`,
			},
			unknown: map[string]string{
				"db.port": "5432",
			},
		},
		{
			name:  "bool and null",
			input: `{"flag": true, "value": null}`,
			awaited: map[string]bool{
				"flag":  true,
				"value": true,
			},
			found: map[string]string{
				"flag":  `true`,
				"value": `null`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "numbers (float, exp, negative)",
			input: `{"pi": 3.1415, "big": 1e6, "neg": -0.5}`,
			awaited: map[string]bool{
				"pi":  true,
				"big": true,
				"neg": true,
			},
			found: map[string]string{
				"pi":  `3.1415`,
				"big": `1e6`,
				"neg": `-0.5`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "empty structures",
			input: `{"emptyObj": {}, "emptyArr": []}`,
			awaited: map[string]bool{
				"emptyObj": true,
				"emptyArr": true,
			},
			found: map[string]string{
				"emptyObj": `{}`,
				"emptyArr": `[]`,
			},
			unknown: map[string]string{},
		},
		{
			name: "escaped strings",
			input: `{
                "path": "C:\\\\Users\\\\Alice",
                "quote": "\"To be\"",
                "euro": "\u20AC"
            }`,
			awaited: map[string]bool{
				"path":  true,
				"quote": true,
				"euro":  true,
			},
			found: map[string]string{
				"path":  `C:\\Users\\Alice`,
				"quote": `"To be"`,
				"euro":  `€`,
			},
			unknown: map[string]string{},
		},
		{
			name: "nested arrays & objects",
			input: `{
                "users": [
                    {"id": 1, "name": "Alice"},
                    {"id": 2, "name": "Bob"}
                ]
            }`,
			awaited: map[string]bool{
				"users": true,
			},
			found: map[string]string{
				"users": `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "array element by index (await whole array)",
			input: `{"tags": ["go", "json", "test"]}`,
			awaited: map[string]bool{
				"tags": true,
			},
			found: map[string]string{
				"tags": `["go","json","test"]`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "duplicate keys",
			input: `{"a": 1, "a": 2}`,
			awaited: map[string]bool{
				"a": true,
			},
			found: map[string]string{
				"a": `2`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "key with dot character (await key itself)",
			input: `{"a.b": {"c": 3}}`,
			awaited: map[string]bool{
				`a.b`: true,
			},
			found: map[string]string{
				`a.b`: `{"c":3}`,
			},
			unknown: map[string]string{},
		},
		{
			name:  "large integer",
			input: `{"big": 9007199254740993}`,
			awaited: map[string]bool{
				"big": true,
			},
			found: map[string]string{
				"big": `9007199254740993`,
			},
			unknown: map[string]string{},
		},
		{
			name:    "json only bool",
			input:   `true`,
			wantErr: true,
		},
		{
			name:    "json only string",
			input:   `"hello"`,
			wantErr: true,
		},
		{
			name:    "json only number",
			input:   `123`,
			wantErr: true,
		},
		{
			name:    "json only null",
			input:   `null`,
			wantErr: true,
		},
		{
			name:    "json only array",
			input:   `[1,2,3]`,
			wantErr: true,
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

			if tt.input == "" {

				require.NotEqual(t, "", path, "tempFile should return a path even for empty content")
			}
			p := json.New(&path)

			found, unknown, err := p.Parse(tt.awaited, zfg.ToString)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.found, found, "Found map mismatch")
				assert.Equal(t, tt.unknown, unknown, "Unknown map mismatch")
			}
		})
	}
}

func TestParse_Error(t *testing.T) {
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
	assert.Contains(t, err.Error(), "is empty")
}

func TestParse_FileNotExist(t *testing.T) {
	nonExistentPath := "no_such_file.json"
	p := json.New(&nonExistentPath)

	found, unknown, err := p.Parse(map[string]bool{"some.key": true}, zfg.ToString)
	require.NoError(t, err)
	assert.Empty(t, found)
	assert.Empty(t, unknown)
}

func tempFile(t *testing.T, data string) string {
	f, err := os.CreateTemp("", "test-*.json")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	if data != "" {
		_, err = f.WriteString(data)
		require.NoError(t, err)
	}

	require.NoError(t, f.Close())

	return f.Name()
}
