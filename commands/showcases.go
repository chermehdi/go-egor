package commands

import (
	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

type TestCaseIO struct {
	Id         int
	Name       string
	InputPath  string
	OutputPath string
	Custom     bool
}

// parse input and output from egor meta to TestCase
func GetTestCases(egorMeta config.EgorMeta) []*TestCaseIO {
	testCasesMap := make(map[string]*TestCaseIO)
	testCases := make([]*TestCaseIO, 0)
	for _, input := range egorMeta.Inputs {
		testCase := &TestCaseIO{
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

// print test cases table
func PrintTestCasesTable(testCases []*TestCaseIO) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Test Name", "Input Path", "Output Path", "Custon"})
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

// print task test cases
func PrintTestCases(egorMeta config.EgorMeta) {
	PrintTestCasesTable(GetTestCases(egorMeta))
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

	PrintTestCases(metaData)

	return nil
}

// Command to print the list of test cases input and outputs into the consol.
// Running this command will fetch egor meta data and load all inputs and outputs meta data
// and prints it as an array into the consol.
var ShowCasesCommand = cli.Command{
	Name:      "showcases",
	Aliases:   []string{"sc"},
	Usage:     "list information about test cases",
	UsageText: "list meta data about of the tests cases in the current task",
	Action:    ShowCasesAction,
}
