package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultConfigPath = "internal/config/profile/"
	defaultScope      = "local"
)

type Reader func() ([]byte, error)

type AppConfiguration struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type HTTPConfiguration struct {
	Port string `json:"port"`
}

type LogConfiguration struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

type SwaggerConfiguration struct {
	Enabled bool `json:"enabled"`
}

type EndpointConfiguration struct {
	BaseURL             string  `json:"base_url"`
	Timeout             int     `json:"timeout"`
	CircuitBreakerRatio float64 `json:"circuit_breaker_ratio"`
}

type Configuration struct {
	App        AppConfiguration      `json:"app"`
	HTTP       HTTPConfiguration     `json:"http"`
	Log        LogConfiguration      `json:"log"`
	Swagger    SwaggerConfiguration  `json:"swagger"`
	ExampleAPI EndpointConfiguration `json:"example_api"`
}

func FileReaderImp() ([]byte, error) {
	configDir := os.Getenv("CONFIG_DIR")

	if configDir == "" {
		configDir = filepath.Join("internal", "config", "profile")
	}

	fileName := "configurations.json"
	filePath := filepath.Join(configDir, fileName)

	absPath, err := filepath.Abs(filePath)
	if err == nil {
		filePath = absPath
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file at %s: %w", filePath, err)
	}

	return bytes, nil
}

func New(reader Reader) (Configuration, error) {
	scope := os.Getenv("SCOPE")
	if scope == "" {
		scope = defaultScope
	}

	raw, err := reader()
	if err != nil {
		return Configuration{}, fmt.Errorf("could not get configuration: %w", err)
	}

	var cfg Configuration
	if err = json.Unmarshal(raw, &cfg); err != nil {
		return Configuration{}, fmt.Errorf("could not unmarshal configuration: %w", err)
	}

	if err = cfg.validate(); err != nil {
		return Configuration{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Configuration) validate() error {
	var errs []error

	if c.App.Name == "" {
		errs = append(errs, errors.New("app.name is required"))
	}
	if c.HTTP.Port == "" {
		errs = append(errs, errors.New("http.port is required"))
	}
	if c.Log.Level == "" {
		errs = append(errs, errors.New("log.level is required"))
	}
	if c.ExampleAPI.BaseURL != "" && c.ExampleAPI.Timeout <= 0 {
		errs = append(errs, errors.New("example_api.timeout must be greater than 0 when base_url is set"))
	}

	return errors.Join(errs...)
}
