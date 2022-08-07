package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Logger struct {
		Level string `yaml:"level"`
		Path  string `yaml:"path"`
	} `yaml:"logger"`
	StorageMode string `yaml:"storageMode"`

	Database struct {
		Driver   string `yaml:"driver"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"http"`
}

func New(path string) (*Config, error) {
	if path == "" {
		return nil, ErrFilePathEmpty
	}

	fh, err := os.Open(path)
	if err != nil {
		return nil, ErrOpenFailed
	}
	defer func() {
		fh.Close()
	}()

	cfg := &Config{}

	decoder := yaml.NewDecoder(fh)
	if err := decoder.Decode(cfg); err != nil {
		return nil, ErrReadFile
	}

	return cfg, nil
}
