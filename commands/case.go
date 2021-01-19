package commands

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// Read from stdin till ctrl D or Command D
func readFromStdin() []string {
	scn := bufio.NewScanner(os.Stdin)
	var lines []string
	for scn.Scan() {
		line := scn.Text()
		if len(line) == 1 {
			if line[0] == '\x1D' {
				break
			}
		}
		lines = append(lines, line)
	}
	return lines
}

// Write given lines to given filename
func writeLinesToFile(filename string, lines []string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	for _, line := range lines {
		_, err := fmt.Fprintln(f, line)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddNewCaseInput Create and save user specified custom case input, and update the given egor meta data
func AddNewCaseInput(inputLines []string,
	caseName string,
	metaData config.EgorMeta) (config.EgorMeta, error) {

	inputFileName := caseName + ".in"
	err := writeLinesToFile(path.Join("inputs", inputFileName), inputLines)
	if err != nil {
		return metaData, err
	}
	inputFile := config.NewIoFile(caseName, path.Join("inputs", inputFileName), true)
	metaData.Inputs = append(metaData.Inputs, inputFile)

	return metaData, nil
}

// AddNewCaseOutput Create and save user specified custom csae output, and update the given egor meta data
func AddNewCaseOutput(outputLines []string,
	caseName string,
	metaData config.EgorMeta) (config.EgorMeta, error) {

	outputFileName := caseName + ".ans"
	err := writeLinesToFile(path.Join("outputs", outputFileName), outputLines)
	if err != nil {
		return metaData, err
	}
	outputFile := config.NewIoFile(caseName, path.Join("outputs", outputFileName), true)
	metaData.Outputs = append(metaData.Outputs, outputFile)

	return metaData, nil
}

// CustomCaseAction Create a user custom test case
func CustomCaseAction(context *cli.Context) error {
	color.Green("Creating Custom Test Case...")

	// Load meta data
	cwd, err := os.Getwd()
	if err != nil {
		color.Red("Failed to Generate Custom Case")
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

	caseName := "test-" + strconv.Itoa(len(metaData.Inputs))
	color.Green("Provide your input:")
	inputLines := readFromStdin()
	metaData, err = AddNewCaseInput(inputLines, caseName, metaData)

	if err != nil {
		color.Red("Failed to add new case input")
		return err
	}

	outputLines := []string{}
	if !context.Bool("no-output") {
		color.Green("Provide your output:")
		outputLines = readFromStdin()
	}

	metaData, err = AddNewCaseOutput(outputLines, caseName, metaData)

	if err != nil {
		color.Red("Failed to add new case output")
		return err
	}

	err = metaData.SaveToFile(path.Join(cwd, configFileName))

	if err != nil {
		color.Red("Failed to save to MetaData")
		return err
	}

	color.Green("Created Custom Test Case")
	return nil
}

// CaseCommand Command to add costum test cases to the current task.
// Running this command will ask the user to provide their input and output, then
// saves the new test case data into appropriate files and add their meta data into
// the current task egor meta data.
// The user can add a flag --no-output to specify that this test case have no output
// associated with it. The user will not be asked to provide output in this case.
var CaseCommand = cli.Command{
	Name:      "case",
	Aliases:   []string{"tc", "testcase"},
	UsageText: "egor case <--no-output?>",
	Usage:     "Add custom test cases to egor task.",
	Action:    CustomCaseAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-output",
			Usage: "This test case doesnt have output",
			Value: false,
		},
	},
}
