package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const LatestVersion = "0.2.0"

// The configuration of the CLI
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	}
	Lang struct {
		Default string `yaml:"default"`
	}
	ConfigFileName     string `yaml:"config_file_name"`
	Version            string `yaml:"version"`
	Author             string `yaml:"author"`
	CppLibraryLocation string `yaml:"cpp_lib_location"`
	// Should contain template locations per programming language code
	// The expected format is: cpp -> /path/to/template_cpp.template
	CustomTemplate map[string]string `yaml:"custom_template"`
}

func (conf *Config) HasCppLibrary() bool {
	return conf.CppLibraryLocation != ""
}

func getDefaultConfigLocation() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, "egor.yaml"), nil
}

func createDefaultConfiguration() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Server: struct {
			Port int `yaml:"port"`
		}{
			Port: 1200,
		},
		Lang: struct {
			Default string `yaml:"default"`
		}{
			Default: "cpp",
		},
		Version:            LatestVersion,
		ConfigFileName:     "egor-meta.json",
		CppLibraryLocation: path.Join(homeDir, "include"),
		CustomTemplate:     make(map[string]string),
	}
}

// This function is called when the configuration file does not exist already
// This will create the configuration file in the user config dir, with a minimalistic
// default configuration
func SaveConfiguration(config *Config) error {
	location, err := getDefaultConfigLocation()
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(location, buffer.Bytes(), 0777)
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
	// Check if the current version (maybe the user already has a configuration file)
	// is an older version. and update accordingly
	if config.Version < LatestVersion {
		config.Version = LatestVersion
		_ = SaveConfiguration(&config)
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
			config := createDefaultConfiguration()
			if err := SaveConfiguration(config); err != nil {
				return nil, err
			}
		}
	}
	return LoadConfiguration(location)
}

// Gets the configuration value associated with the given key
func GetConfigurationValue(config *Config, key string) (string, error) {
	lowerKey := strings.ToLower(key)
	if lowerKey == "server.port" {
		return strconv.Itoa(config.Server.Port), nil
	} else if lowerKey == "lang.default" {
		return config.Lang.Default, nil
	} else if lowerKey == "author" {
		return config.Author, nil
	} else if lowerKey == "cpp.lib.location" {
		return config.CppLibraryLocation, nil
	} else {
		return "", errors.New(fmt.Sprintf("Unknown config key %s", key))
	}
}
