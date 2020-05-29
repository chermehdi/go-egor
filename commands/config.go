package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func UpdateConfiguration(config *config.Config, key, value string) error {
	lowerKey := strings.ToLower(key)
	if lowerKey == "server.port" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		config.Server.Port = port
	} else if lowerKey == "lang.default" {
		// suppose that the language is valid (cpp, python, java ...)
		config.Lang.Default = value
	} else if lowerKey == "author" {
		config.Author = value
	} else if lowerKey == "cpp.lib.location" {
		config.CppLibraryLocation = value
	} else if strings.HasPrefix(key, "custom.template") {
		lang := key[strings.LastIndex(key, ".")+1:]
		config.CustomTemplate[lang] = value
	} else {
		// Unknow key
		return errors.New(fmt.Sprintf("Unknown configuration property %s", key))
	}
	return nil
}

// Sets the configuration property identified by the first argument
// To the value identified by the second argument.
// Note: The configuration keys are not case sensitive, if a configuration key provided
// is not recognized, an error is thrown
func SetAction(context *cli.Context) error {
	argLen := context.Args().Len()
	if argLen != 2 {
		color.Red(fmt.Sprintln("Usage egor config set key value"))
		return errors.New(fmt.Sprintf("Expected 2 parameters, got %d", argLen))
	}
	key, value := context.Args().Get(0), context.Args().Get(1)

	configuration, err := config.LoadDefaultConfiguration()
	if err != nil {
		return err
	}
	if err := UpdateConfiguration(configuration, key, value); err != nil {
		return err
	}

	return config.SaveConfiguration(configuration)
}

// Gets and prints the current configuration associted to the first argument,
// or Prints all if no argument is specified
func GetAction(context *cli.Context) error {
	argLen := context.Args().Len()
	if argLen > 1 {
		color.Red(fmt.Sprintln("Usage egor config get <key?>"))
		return errors.New(fmt.Sprintf("Expected at most 1 parameter, got %d", argLen))
	}
	configuration, err := config.LoadDefaultConfiguration()
	if err != nil {
		return err
	}
	if argLen == 0 {
		fmt.Println("Current configuration: ")
		color.Green("server.port     \t\t %d\n", configuration.Server.Port)
		color.Green("lang.default    \t\t %s\n", configuration.Lang.Default)
		color.Green("author          \t\t %s\n", configuration.Author)
		color.Green("cpp.lib.location\t\t %s\n", configuration.CppLibraryLocation)
		color.Green("config.templates: \n\n")
		for k, v := range configuration.CustomTemplate {
			color.Green("\t%s\t\t %s\n", k, v)
		}
		return nil
	} else {
		value, err := config.GetConfigurationValue(configuration, context.Args().First())
		if err != nil {
			return err
		}
		color.Green("%s\t\t %s\n", context.Args().First(), value)
		return nil
	}
}

var ConfigCommand = cli.Command{
	Name:      "config",
	Aliases:   []string{"c"},
	Usage:     "Read/Change global configuration parameters",
	UsageText: "egor config set <config.name> <config.value> | egor config get <config.name?>",
	Subcommands: []*cli.Command{
		{
			Name:      "set",
			Usage:     "set config parameter",
			UsageText: "Sets a configuration parameter",
			Action:    SetAction,
		},
		{
			Name:      "get",
			Usage:     "get config parameter",
			UsageText: "Gets a configuration parameter",
			Action:    GetAction,
		},
	},
}
