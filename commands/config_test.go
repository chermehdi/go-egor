package commands

import (
	"testing"

	"github.com/chermehdi/egor/config"
	"github.com/stretchr/testify/assert"
)

// TODO(chermehdi): This probably should be a test fixture.
func getDefaultConfiguration() *config.Config {
	return &config.Config{
		Server: struct {
			Port int `yaml:"port"`
		}{Port: 1200},
		Lang: struct {
			Default string `yaml:"default"`
		}{Default: "cpp"},
	}
}

func TestSetConfiguration(t *testing.T) {
	configuration := getDefaultConfiguration()

	_ = UpdateConfiguration(configuration, "server.port", "1245")

	assert.Equal(t, configuration.Server.Port, 1245)
	assert.Equal(t, configuration.Lang.Default, "cpp")

	_ = UpdateConfiguration(configuration, "lang.default", "java")
	assert.Equal(t, configuration.Lang.Default, "java")
}

func TestSetConfigurationUnknownKey(t *testing.T) {
	configuration := getDefaultConfiguration()
	err := UpdateConfiguration(configuration, "unkown.key", "123")

	assert.Error(t, err, "Error should be returned if the key is unknown")
}
