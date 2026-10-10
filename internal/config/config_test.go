package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAppConfig(t *testing.T) {
	expectedConfig := AppConfig{
		NetAddress: NetAddress{
			Host: "localhost",
			Port: 1234,
		},
		BaseURL: "http://localhost:1234",
	}

	cases := []struct {
		name           string
		envArgs        map[string]string
		cmdArgs        map[string]string
		expectedConfig AppConfig
	}{
		{
			name:           "no env or cmd",
			envArgs:        map[string]string{},
			cmdArgs:        map[string]string{},
			expectedConfig: GetAppConfigDefaults(),
		},
		{
			name: "with env only",
			envArgs: map[string]string{
				"SERVER_ADDRESS": "localhost:1234",
				"BASE_URL":       "http://localhost:1234",
			},
			cmdArgs:        map[string]string{},
			expectedConfig: expectedConfig,
		},
		{
			name:    "with cmd only",
			envArgs: map[string]string{},
			cmdArgs: map[string]string{
				"-a": "localhost:1234",
				"-b": "http://localhost:1234",
			},
			expectedConfig: expectedConfig,
		},
		{
			name: "with env and cmd",
			envArgs: map[string]string{
				"SERVER_ADDRESS": "localhost:1234",
				"BASE_URL":       "http://localhost:1234",
			},
			cmdArgs: map[string]string{
				"-a": "localhost:4321",
				"-b": "http://localhost:4321",
			},
			expectedConfig: expectedConfig,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			oldArgs := os.Args
			oldFlags := flag.CommandLine

			t.Cleanup(func() {
				os.Args = oldArgs
				flag.CommandLine = oldFlags
			})

			flag.CommandLine = flag.NewFlagSet("shortener", flag.ContinueOnError)

			os.Args = []string{"shortener"}

			for k, v := range tc.cmdArgs {
				os.Args = append(os.Args, k, v)
			}

			for k, v := range tc.envArgs {
				t.Setenv(k, v)
			}

			res, err := GetAppConfig()

			assert.NoError(t, err)
			tcNetAddr := tc.expectedConfig.NetAddress
			resNetAddr := res.NetAddress

			assert.Equal(t, tc.expectedConfig.BaseURL, res.BaseURL, "Base url does not match")
			assert.Equal(t, tcNetAddr.Port, resNetAddr.Port, "Port does not match")
			assert.Equal(t, tcNetAddr.Host, resNetAddr.Host, "Host does not match")
		})
	}
}
