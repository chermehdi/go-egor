package commands

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/chermehdi/egor/config"
	"github.com/chermehdi/egor/utils"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"strconv"
)

func GetTestCase(egorMeta config.EgorMeta, id int) *TestCaseIO {
	var testCase *TestCaseIO
	for _, input := range egorMeta.Inputs {
		if input.GetId() == id {
			testCase = &TestCaseIO{
				Id:         input.GetId(),
				Name:       input.Name,
				InputPath:  input.Path,
				OutputPath: "",
				Custom:     input.Custom,
			}

			break
		}
	}

	if testCase == nil {
		return nil
	}

	for _, output := range egorMeta.Outputs {
		if output.Name == testCase.Name {
			testCase.OutputPath = output.Path
		}
	}

	if testCase == nil {
		return nil
	}
	return testCase
}

func PrintTestCaseInput(testCase *TestCaseIO) {

	file, err := config.OpenFileFromPath(testCase.InputPath)
	if err != nil {
		color.Red("Failed to read test case input")
	} else {
		color.Green(utils.GetHeader("Input", "========", 10))
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}

}

func PrintTestCaseOutput(testCase *TestCaseIO) {
	file, err := config.OpenFileFromPath(testCase.OutputPath)
	if err != nil {
		color.Red("Failed to read test case input")
	} else {
		color.Green(utils.GetHeader("Output", "========", 10))
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}
}

func PrintCaseAction(context *cli.Context) error {
	if context.Bool("input-only") && context.Bool("output-only") {
		color.Red("only --input-only or --output-only should be set.")
		return errors.New("Invalid commands arguments")
	}

	if context.NArg() == 0 {
		color.Red("Test id required argument missing!")
		return errors.New("Missing required argument 'test_id'")
	}

	id, err := strconv.Atoi(context.Args().Get(0))

	if err != nil {
		color.Red(fmt.Sprintf("Cannot parse test id = '%s', a number required!", context.Args().Get(0)))
		return errors.New(fmt.Sprintf("Failed to parse test id = %s", context.Args().Get(0)))
	}

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

	testCase := GetTestCase(metaData, id)
	if testCase == nil {
		color.Red(fmt.Sprintf("Could not find test case with id = %d", id))
		return errors.New(fmt.Sprintf("Unknown test case with id %d", id))
	}

	if !context.Bool("output-only") {
		PrintTestCaseInput(testCase)
	}

	if !context.Bool("input-only") {
		PrintTestCaseOutput(testCase)
	}

	return nil
}

// Command to print a test case. this command can be used to print inputs and/or outputs
// to the consol. The user can choose to print the input only or the output only. The
// user should provide a valid test id.
// Running this command will fetch egor meta data, get the test case with the given id,
// and then print the content of the input and/or of the output files.
var PrintCaseCommand = cli.Command{
	Name:      "printcase",
	Aliases:   []string{"pc"},
	Usage:     "print input and/or output of a given test case",
	UsageText: "print input and/or output of a given test case",
	Action:    PrintCaseAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "output-only",
			Usage: "Print the output only of the test case",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "input-only",
			Usage: "Print the input only of the test case",
			Value: false,
		},
	},
}
