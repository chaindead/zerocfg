package env_test

import (
	"os"
	"testing"

	zfg "github.com/chaindead/zerocfg"
	"github.com/chaindead/zerocfg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		awaited map[string]bool
		found   map[string]string
	}{
		{
			name: "simple env",
			envs: map[string]string{
				"ENV1": "bar",
			},
			awaited: map[string]bool{"env1": true},
			found:   map[string]string{"env1": "bar"},
		},
		{
			name: "no env",
			envs: map[string]string{
				"ENV1": "bar",
			},
			awaited: map[string]bool{},
			found:   map[string]string{},
		},
		{
			name: "composite env",
			envs: map[string]string{
				"A_B_C_D": "bar",
			},
			awaited: map[string]bool{"a.b.c.d": true},
			found:   map[string]string{"a.b.c.d": "bar"},
		},
		{
			name: "composite env",
			envs: map[string]string{
				"CAMELCASE_DASH_UNDERWEAR": "bar",
			},
			awaited: map[string]bool{"camelCase.da-sh.under_wear": true},
			found:   map[string]string{"camelCase.da-sh.under_wear": "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := env.New()

			t.Cleanup(func() {
				for k := range tt.envs {
					require.NoError(t, os.Unsetenv(k))
				}
			})
			for k, v := range tt.envs {
				_, found := os.LookupEnv(k)
				require.False(t, found)

				err := os.Setenv(k, v)
				require.NoError(t, err)
			}

			found, unknown, err := p.Parse(tt.awaited, zfg.ToString)
			require.NoError(t, err)
			assert.Empty(t, unknown)

			assert.Equal(t, tt.found, found)

		})
	}
}
