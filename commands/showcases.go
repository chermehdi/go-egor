package commands

import (
	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

// print test cases table
func PrintTestCasesTable(inputFiles, outputFiles map[string]config.IoFile) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Test Name", "Input Path", "Output Path", "Custon"})
	for key, inputFile := range inputFiles {
		id := inputFile.GetId()
		testName := inputFile.Name
		inputPath := inputFile.Path
		outputPath := ""
		custom := inputFile.Custom

		outputFile, ok := outputFiles[key]
		if ok {
			outputPath = outputFile.Path
		}

		t.AppendRow([]interface{}{id, testName, inputPath, outputPath, custom})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

// construct inputs and outputs maps from an egor meta data where keys are test names.
func GetIoFilesMaps(egorMeta config.EgorMeta) (map[string]config.IoFile, map[string]config.IoFile) {
	inputFiles := make(map[string]config.IoFile)
	for _, inputFile := range egorMeta.Inputs {
		inputFiles[inputFile.Name] = inputFile
	}

	outputFiles := make(map[string]config.IoFile)
	for _, outputFile := range egorMeta.Outputs {
		outputFiles[outputFile.Name] = outputFile
	}

	return inputFiles, outputFiles
}

// print task test cases
func PrintTestCases(egorMeta config.EgorMeta) error {
	inputFiles, outputFiles := GetIoFilesMaps(egorMeta)
	PrintTestCasesTable(inputFiles, outputFiles)
	return nil
}

// list test cases information command action
// TODO(Eroui): [Refactoring] Duplicate code while loading meta data, consider refactoring...
func ShowCasesAction(context *cli.Context) error {
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

	color.Green("Listing %d testcase(s)...", metaData.CountTestCases())
	color.Green("")
	
	err = PrintTestCases(metaData)
	if err != nil {
		color.Red("Error while printing test cases")
		return err
	}
	return nil
}

// Command to print the list of test cases input and outputs into the consol.
// Running this command will fetch egor meta data and load all inputs and outputs meta data
// and prints it as an array into the consol.
var ShowCasesCommand = cli.Command{
	Name:      "showcases",
	Aliases:   []string{"listcases"},
	Usage:     "list information about test cases",
	UsageText: "list meta data about of the tests cases in the current task",
	Action:    ShowCasesAction,
}
