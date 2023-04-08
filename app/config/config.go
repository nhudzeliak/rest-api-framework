package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/robfig/config"
)

var (
	// envCache stores env flag value.
	envCache string
	// configCache stores the value of LoadConfig after its first call.
	configCache map[string]string
	// basePathCache stores the value of BasePath after its first call.
	basePathCache string
)

// BasePath fetches project path, caching the result.
func BasePath() string {
	if basePathCache == "" {
		basePathCache = evaluateBasePath()
	}
	return basePathCache
}

// evaluateBasePath evaluates the base path for project.
func evaluateBasePath() string {
	if projectPath := os.Getenv("PROJECT_PATH"); projectPath != "" {
		return projectPath
	}
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		gopath = filepath.Join(strings.Split(gopath, string(os.PathListSeparator))[0], "src")
		return gopath + "/github.com/nataliia_hudzeliak/rest-api-framework"
	}
	return "."
}

// Config fetches project configs based on current ENV, caching the result.
func Config() (map[string]string, error) {
	if configCache == nil {
		cfg, err := evaluateConfig(env())
		if err != nil {
			return nil, err
		}
		configCache = cfg
	}
	return configCache, nil
}

// MustConfig wraps Config and panics if an error is returned.
func MustConfig() map[string]string {
	cfg, err := Config()
	if err != nil {
		panic(err)
	}
	return cfg
}

// evaluateConfig loads configs based on ENV flag.
func evaluateConfig(env string) (map[string]string, error) {
	cfg, err := config.ReadDefault(BasePath() + "/app/config/app.conf")
	if err != nil {
		return nil, err
	}
	options, err := cfg.Options(env)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	for _, o := range options {
		value, err := cfg.String(env, o)
		if err != nil {
			return nil, err
		}
		params[o] = strings.Trim(value, "\"")
	}
	return params, nil
}

// env fetches current env.
func env() string {
	if envCache == "" {
		envCache = evaluateEnv()
	}
	return envCache
}

// evaluateEnv fetches ENV flag from environment variables. Returns "default" if one is fails to be found.
func evaluateEnv() string {
	if env := os.Getenv("ENV"); env != "" {
		return env
	}
	return "default"
}
