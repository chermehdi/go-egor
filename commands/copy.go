package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/atotto/clipboard"
	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// GetFileContent Load the content of a given file
func GetFileContent(filePath string) (string, error) {
	filebytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	filecontent := string(filebytes)
	return filecontent, nil
}

func CopyAction(*cli.Context) error {
	conf, err := config.LoadDefaultConfiguration()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to load default configuration: %s", err.Error()))
		return err
	}
	metaData, err := config.GetMetadata()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to load egor MetaData : %s", err.Error()))
		return err
	}

	var taskFile string

	// TODO(chermehdi): the name of the generated file should be in a unique location
	if (metaData.TaskLang == "cpp" || metaData.TaskLang == "c") && conf.HasCppLibrary() {
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

// CopyCommand Command to copy task source file to clipboard for easy submit.
// Running this command will fetch egor meta data, get the content of the task source
// and then copy the content to the clipboard.
var CopyCommand = cli.Command{
	Name:      "copy",
	Aliases:   []string{"cp"},
	Usage:     "Copy task file into clipboard",
	UsageText: "Copy task file into clipboard",
	Action:    CopyAction,
}
