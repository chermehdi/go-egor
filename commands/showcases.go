package commands

import (
	"os"

	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
)

type testCaseIO struct {
	Id         int
	Name       string
	InputPath  string
	OutputPath string
	Custom     bool
}

// GetTestCases parse input and output from egor meta to TestCase
func GetTestCases(egorMeta *config.EgorMeta) []*testCaseIO {
	testCasesMap := make(map[string]*testCaseIO)
	testCases := make([]*testCaseIO, 0)
	for _, input := range egorMeta.Inputs {
		testCase := &testCaseIO{
			Id:         input.GetId(),
			Name:       input.Name,
			InputPath:  input.Path,
			OutputPath: "",
			Custom:     input.Custom,
		}
		testCasesMap[input.Name] = testCase
		testCases = append(testCases, testCase)
	}

	for _, output := range egorMeta.Outputs {
		if testCase, ok := testCasesMap[output.Name]; ok {
			testCase.OutputPath = output.Path
		}
	}

	return testCases
}

// PrintTestCasesTable print test cases table
func PrintTestCasesTable(testCases []*testCaseIO) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Test Name", "Input Path", "Output Path", "Custom"})
	for _, testCase := range testCases {
		t.AppendRow([]interface{}{
			testCase.Id,
			testCase.Name,
			testCase.InputPath,
			testCase.OutputPath,
			testCase.Custom,
		})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

// PrintTestCases print task test cases
func PrintTestCases(egorMeta *config.EgorMeta) {
	PrintTestCasesTable(GetTestCases(egorMeta))
}

// ShowCasesAction list test cases information command action
func ShowCasesAction(_ *cli.Context) error {
	metaData, err := config.GetMetadata()
	if err != nil {
		color.Red("Failed to load egor MetaData ")
		return err
	}

	color.Green("Listing %d testcase(s)...", metaData.CountTestCases())
	color.Green("")
	PrintTestCases(metaData)
	return nil
}

// ShowCasesCommand Command to print the list of test cases input and outputs into the consol.
// Running this command will fetch egor meta data and load all inputs and outputs meta data
// and prints it as an array into the consol.
var ShowCasesCommand = cli.Command{
	Name:      "showcases",
	Aliases:   []string{"sc"},
	Usage:     "list meta data about tests cases in the current task",
	UsageText: "egor showcases",
	Action:    ShowCasesAction,
}
