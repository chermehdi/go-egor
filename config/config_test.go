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

var configurationPath = getConfigPath()

func getConfigPath() string {
	tempDir, _ := ioutil.TempDir("", "temp")

	return path.Join(tempDir, "config.yaml")
}

func createFakeConfigFile() error {
	configPath := configurationPath

	var buffer bytes.Buffer
	configuration := createDefaultConfiguration()

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
	config, err := LoadConfiguration(configurationPath)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, config.Lang.Default, createDefaultConfiguration().Lang.Default)
	assert.Equal(t, config.Server.Port, createDefaultConfiguration().Server.Port)
}

func deleteFakeConfigFile() {
	configPath := configurationPath
	_ = os.Remove(configPath)
}

func TestGetConfigurationValue(t *testing.T) {
	config := createDefaultConfiguration()

	value, err := GetConfigurationValue(config, "server.port")
	assert.NoError(t, err, "No error should be thrown when getting an existing key")
	assert.Equal(t, value, "1200")

	value, err = GetConfigurationValue(config, "cpp.lib.location")
	assert.NoError(t, err, "No error should be thrown when getting cpp library location")

	_, err = GetConfigurationValue(config, "unknown.key")
	assert.Error(t, err, "An error is returned if the configuration key is not known")

    value, err = GetConfigurationValue(config, "config.templates")
    assert.NoError(t, err,"No error should be thrown when getting custom template")
}

func TestConfig_HasCppLibrary(t *testing.T) {
	config := &Config{
		CppLibraryLocation: "/include",
	}
	assert.True(t, config.HasCppLibrary())

	config = &Config{}
	assert.False(t, config.HasCppLibrary())
}
