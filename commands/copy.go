package commands

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
)

// Load the content of a given file
func GetFileContent(filePath string) (string, error) {
	filebytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	filecontent := string(filebytes)
	return filecontent, nil
}

func CopyAction(*cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to list test cases : %s", err.Error()))
		return err
	}

	configuration, err := config.LoadDefaultConfiguration()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to load egor configuration: %s", err.Error()))
		return err
	}

	configFileName := configuration.ConfigFileName
	metaData, err := config.LoadMetaFromPath(path.Join(cwd, configFileName))
	if err != nil {
		color.Red(fmt.Sprintf("Failed to load egor MetaData : %s", err.Error()))
		return err
	}
	var taskFile string
	//TODO(chermehdi): make the test on weather we have library location set in the configuration object.
	// TODO(chermehdi): the name of the generated file should be in a unique location
	if (metaData.TaskLang == "cpp" || metaData.TaskLang == "c") && configuration.LibraryLocation != "" {
		taskFile = "main_gen.cpp"
	} else {
		taskFile = metaData.TaskFile
	}
	taskContent, err := GetFileContent(taskFile)
	if err != nil {
		color.Red(fmt.Sprintf("Failed to load task file content : %s", err.Error()))
		return err
	}

	err = clipboard.WriteAll(taskContent)
	if err != nil {
		color.Red(fmt.Sprintf("Failed to copy task content to clipboard : %s", err.Error()))
		return err
	}

	color.Green("Task copied to clipboard successfully")
	return nil
}

// Command to copy task source file to clipboard for easy submit.
// Running this command will fetch egor meta data, get the content of the task source
// and then copy the content to the clipboard.
var CopyCommand = cli.Command{
	Name:      "copy",
	Aliases:   []string{"cp"},
	Usage:     "Copy task file into clipboard",
	UsageText: "Copy task file into clipboard",
	Action:    CopyAction,
}
