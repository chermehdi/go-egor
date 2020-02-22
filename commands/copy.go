package commands

import (
	"github.com/urfave/cli/v2"
	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/atotto/clipboard"
	"os"
	"path"
	"io/ioutil"
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
		color.Red("Failed to list test cases!")
		return err
	}

	configuration, err := config.LoadDefaultConfiguration()
	if err != nil {
		color.Red("Failed to load egor configuration")
		return err
	}

	configFileName := configuration.ConfigFileName
	metaData, err := config.LoadMetaFromPath(path.Join(cwd, configFileName))
	if err != nil {
		color.Red("Failed to load egor MetaData ")
		return err
	}

	taskFile := metaData.TaskFile
	taskContent, err := GetFileContent(taskFile)
	if err != nil {
		color.Red("Failed to load task file content")
		return err
	}

	err = clipboard.WriteAll(taskContent)
	if err != nil {
		color.Red("Failed to copy task content to clipboard")
		return err
	}

	color.Green("Done!")
	return nil
}

// Command to copy task source file to clipboard for easy submit.
// Running this command will fetch egor meta data, get the content of the task source
// and then copy the content to the clipboard.
var CopyCommand = cli.Command{
	Name:      "showcases",
	Aliases:   []string{"cp"},
	Usage:     "copy task file into clipboad",
	UsageText: "list meta data about of the tests cases in the current task",
	Action:    CopyAction,
}
