package config

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

// The configuration of the CLI
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	}
	Lang struct {
		Default string `yaml:"default"`
	}
	Version string `yaml:"version"`
}

func getDefaultConfigLocation() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, "egor.yaml"), nil
}

func createDefaultConfiguration() *Config {
	return &Config{
		Server: struct {
			Port int `yaml:"port"`
		}{
			Port: 12,
		},
		Lang: struct {
			Default string `yaml:"default"`
		}{
			Default: "cpp",
		},
		Version: "1.0",
	}
}

// This function is called when the configuration file does not exist already
// This will create the configuration file in the user config dir, with a minimalistic
// default configuration
func saveDefaultConfiguration() error {
	location, err := getDefaultConfigLocation()
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	err = encoder.Encode(createDefaultConfiguration())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(location, buffer.Bytes(), 0644)
}

// Returns the Configuration object associated with
// the path given as a parameter
func LoadConfiguration(location string) (*Config, error) {
	file, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Returns the Configuration object associated with
// the default configuration location
func LoadDefaultConfiguration() (*Config, error) {
	location, err := getDefaultConfigLocation()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(location); err != nil {
		if os.IsNotExist(err) {
			saveDefaultConfiguration()
		}
	}
	return LoadConfiguration(location)
}
