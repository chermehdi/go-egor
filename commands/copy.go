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

func getGenFile(metaData *config.EgorMeta, conf *config.Config) string {
	if (metaData.TaskLang == "cpp" || metaData.TaskLang == "c") && conf.HasCppLibrary() {
		return "main_gen.cpp"
	}
	return metaData.TaskFile
}

func copyToClipboard(taskFile string) error {
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

	taskFile := getGenFile(metaData, conf)
	return copyToClipboard(taskFile)
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
