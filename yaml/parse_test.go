package yaml_test

import (
	"os"
	"testing"

	zfg "github.com/chaindead/zeroconf"
	"github.com/chaindead/zeroconf/yaml"
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
	}{
		{
			name: "simple key-value",
			input: `
str: name
int: 1`,
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
			input: `
host: localhost
database:
  port: 5432
  credentials:
    username: admin`,
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
			name: "array",
			input: `
tags:
  - tag1
  - tag2
  - tag3`,
			awaited: map[string]bool{
				"tags": true,
			},
			found: map[string]string{
				"tags": `["tag1","tag2","tag3"]`,
			},
		},
		{
			name: "map",
			input: `
tags:
  k1: 1
  k2: 2`,
			awaited: map[string]bool{
				"tags": true,
			},
			found: map[string]string{
				"tags": `{"k1":1,"k2":2}`,
			},
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

			name := tempFile(t, tt.input)
			p := yaml.New(&name)

			found, unknown, err := p.Parse(tt.awaited, zfg.ToString)
			require.NoError(t, err)

			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.unknown, unknown)
		})
	}
}

func TestParse_Error(t *testing.T) {
	name := tempFile(t, `invalid: [yaml: "missing closing quote`)
	p := yaml.New(&name)

	_, _, err := p.Parse(map[string]bool{}, zfg.ToString)
	assert.Error(t, err)
}

func tempFile(t *testing.T, data string) string {
	f, err := os.CreateTemp("", "tmpfile-")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, f.Close())
	})

	_, err = f.WriteString(data)
	require.NoError(t, err)

	return f.Name()
}
