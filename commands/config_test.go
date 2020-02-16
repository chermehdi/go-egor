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
	config := getDefaultConfiguration()

	UpdateConfiguration(config, "server.port", "1245")

	assert.Equal(t, config.Server.Port, 1245)
	assert.Equal(t, config.Lang.Default, "cpp")

	UpdateConfiguration(config, "lang.default", "java")
	assert.Equal(t, config.Lang.Default, "java")
}

func TestSetConfigurationUnknownKey(t *testing.T) {
	config := getDefaultConfiguration()
	err := UpdateConfiguration(config, "unkown.key", "123")

	assert.Error(t, err, "Error should be returned if the key is unknown")
}
