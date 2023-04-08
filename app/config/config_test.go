package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasePath(t *testing.T) {
	cases := []struct {
		setup     func() (map[string]string, error)
		base      func() string
		assertion func(string)
		cleanup   func(map[string]string) error
	}{
		// PROJECT_PATH present.
		{
			setup: func() (map[string]string, error) {
				memory := make(map[string]string)
				for _, k := range []string{"PROJECT_PATH"} {
					memory[k] = os.Getenv(k)
				}
				basePathCache = ""
				err := os.Setenv("PROJECT_PATH", "test-project-path")
				if err != nil {
					return nil, err
				}
				return memory, nil
			},
			base: func() string {
				return BasePath()
			},
			assertion: func(value string) {
				assert.Equal(t, "test-project-path", value)
			},
			cleanup: func(memory map[string]string) error {
				for k, v := range memory {
					err := os.Setenv(k, v)
					if err != nil {
						return err
					}
				}
				basePathCache = ""
				return nil
			},
		},
		// PROJECT_PATH missing, GOPATH present.
		{
			setup: func() (map[string]string, error) {
				memory := make(map[string]string)
				for _, k := range []string{"PROJECT_PATH", "GOPATH"} {
					memory[k] = os.Getenv(k)
				}
				basePathCache = ""
				err := os.Setenv("PROJECT_PATH", "")
				if err != nil {
					return nil, err
				}
				err = os.Setenv("GOPATH", "test-gopath")
				if err != nil {
					return nil, err
				}
				return memory, nil
			},
			base: func() string {
				return BasePath()
			},
			assertion: func(value string) {
				assert.Equal(t, "test-gopath/src/github.com/nataliia_hudzeliak/rest-api-framework", value)
			},
			cleanup: func(memory map[string]string) error {
				for k, v := range memory {
					err := os.Setenv(k, v)
					if err != nil {
						return err
					}
				}
				basePathCache = ""
				return nil
			},
		},
		// Both missing.
		{
			setup: func() (map[string]string, error) {
				memory := make(map[string]string)
				for _, k := range []string{"PROJECT_PATH", "GOPATH"} {
					memory[k] = os.Getenv(k)
				}
				basePathCache = ""
				err := os.Setenv("PROJECT_PATH", "")
				if err != nil {
					return nil, err
				}
				err = os.Setenv("GOPATH", "")
				if err != nil {
					return nil, err
				}
				return memory, nil
			},
			base: func() string {
				return BasePath()
			},
			assertion: func(value string) {
				assert.Equal(t, ".", value)
			},
			cleanup: func(memory map[string]string) error {
				for k, v := range memory {
					err := os.Setenv(k, v)
					if err != nil {
						return err
					}
				}
				basePathCache = ""
				return nil
			},
		},
	}

	for _, c := range cases {
		memory, err := c.setup()
		assert.Nil(t, err)
		result := c.base()
		c.assertion(result)
		err = c.cleanup(memory)
		assert.Nil(t, err)
	}
}

func TestConfig(t *testing.T) {
	cases := []struct {
		setup     func() (string, error)
		base      func() (string, error)
		assertion func(string)
		cleanup   func(string) error
	}{
		// Default.
		{
			setup: func() (string, error) {
				memory := os.Getenv("ENV")
				configCache = nil
				envCache = ""
				err := os.Setenv("ENV", "")
				if err != nil {
					return "", err
				}
				return memory, nil
			},
			base: func() (string, error) {
				cfg, err := Config()
				if err != nil {
					return "", err
				}
				return cfg["envname"], nil
			},
			assertion: func(value string) {
				assert.Equal(t, "default", value)
			},
			cleanup: func(memory string) error {
				configCache = nil
				envCache = ""
				return os.Setenv("ENV", memory)
			},
		},
		// Test.
		{
			setup: func() (string, error) {
				memory := os.Getenv("ENV")
				configCache = nil
				envCache = ""
				err := os.Setenv("ENV", "test")
				if err != nil {
					return "", err
				}
				return memory, nil
			},
			base: func() (string, error) {
				cfg, err := Config()
				if err != nil {
					return "", err
				}
				return cfg["envname"], nil
			},
			assertion: func(value string) {
				assert.Equal(t, "test", value)
			},
			cleanup: func(memory string) error {
				configCache = nil
				envCache = ""
				return os.Setenv("ENV", memory)
			},
		},
		// Production.
		{
			setup: func() (string, error) {
				memory := os.Getenv("ENV")
				configCache = nil
				envCache = ""
				err := os.Setenv("ENV", "production")
				if err != nil {
					return "", err
				}
				return memory, nil
			},
			base: func() (string, error) {
				cfg, err := Config()
				if err != nil {
					return "", err
				}
				return cfg["envname"], nil
			},
			assertion: func(value string) {
				assert.Equal(t, "production", value)
			},
			cleanup: func(memory string) error {
				configCache = nil
				envCache = ""
				return os.Setenv("ENV", memory)
			},
		},
	}

	for _, c := range cases {
		memory, err := c.setup()
		assert.Nil(t, err)
		value, err := c.base()
		assert.Nil(t, err)
		c.assertion(value)
		err = c.cleanup(memory)
		assert.Nil(t, err)
	}
}
