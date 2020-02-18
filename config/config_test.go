package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func getDefaultConfiguration() *Config {
	return &Config{
		Server: struct {
			Port int `yaml:"port"`
		}{Port: 1200},
		Lang: struct {
			Default string `yaml:"default"`
		}{Default: "cpp"},
	}
}

func getConfigPath() string {
	tempDir := os.TempDir()
	return path.Join(tempDir, "config.yaml")
}

func createFakeConfigFile() error {
	configPath := getConfigPath()

	var buffer bytes.Buffer
	configuration := getDefaultConfiguration()

	encoder := yaml.NewEncoder(&buffer)
	err := encoder.Encode(configuration)
	if err != nil {
		return err
	}
	// write the fake configuration yaml to the file
	err = ioutil.WriteFile(configPath, buffer.Bytes(), 777)
	if err != nil {
		return err
	}
	return nil
}

func TestLoadConfiguration(t *testing.T) {
	_ = createFakeConfigFile()
	defer deleteFakeConfigFile()
	config, err := LoadConfiguration(getConfigPath())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, config.Lang.Default, getDefaultConfiguration().Lang.Default)
	assert.Equal(t, config.Server.Port, getDefaultConfiguration().Server.Port)
}

func deleteFakeConfigFile() {
	configPath := getConfigPath()
	_ = os.Remove(configPath)
}

func TestGetConfigurationValue(t *testing.T) {
	config := createDefaultConfiguration()

	value, err := GetConfigurationValue(config, "server.port")
	assert.NoError(t, err, "No error should be thrown when getting an existing key")
	assert.Equal(t, value, "1200")

	_, err = GetConfigurationValue(config, "unknown.key")
	assert.Error(t, err, "An error is returned if the configuration key is not known")
}
