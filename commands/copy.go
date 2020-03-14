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

func CopyAction(context *cli.Context) error {
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

	taskFile := metaData.TaskFile
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
	Usage:     "copy task file into clipboad",
	UsageText: "list meta data about of the tests cases in the current task",
	Action:    CopyAction,
}
