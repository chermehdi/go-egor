package commands

import (
	"bufio"
	"fmt"
	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"strconv"
)

func readFromStdin() []string {
	scn := bufio.NewScanner(os.Stdin)
	var lines []string
	for scn.Scan() {
		line := scn.Text()
		if len(line) == 1 {
			// Group Separator (GS ^]): ctrl-]
			if line[0] == '\x1D' {
				break
			}
		}
		lines = append(lines, line)
	}
	return lines
}

func writeLinesToFile(filename string, lines []string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, line := range lines {
		fmt.Fprintln(f, line)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func AddNewCaseInput(inputLines []string,
	caseName string,
	metaData config.EgorMeta,
	noTimeOut bool) (config.EgorMeta, error) {

	inputFileName := caseName + ".in"
	writeLinesToFile(path.Join("inputs", inputFileName), inputLines)
	inputFile := config.NewIoFile(inputFileName, path.Join("inputs", inputFileName), true, noTimeOut)
	metaData.Inputs = append(metaData.Inputs, inputFile)

	return metaData, nil
}

func AddNewCaseOutput(outputLines []string,
	caseName string,
	metaData config.EgorMeta,
	noTimeOut bool) (config.EgorMeta, error) {

	outputFileName := caseName + ".ans"
	writeLinesToFile(path.Join("outputs", outputFileName), outputLines)
	outputFile := config.NewIoFile(outputFileName, path.Join("outputs", outputFileName), true, noTimeOut)
	metaData.Outputs = append(metaData.Outputs, outputFile)

	return metaData, nil
}

// TODO(Eroui): add checks on errors
func CustomCaseAction(context *cli.Context) error {
	color.Green("Creating Custom Test Case...")

	// Load meta data
	cwd, err := os.Getwd()
	if err != nil {
		color.Red("Failed to Generate Custom Case")
		return err
	}

	metaData, err := config.LoadMetaFromPath(path.Join(cwd, "egor-meta.json"))
	if err != nil {
		color.Red("Failed to load egor MetaData ")
		return err
	}

	noTimeOut := context.Bool("no-timeout")

	caseName := "test-" + strconv.Itoa(len(metaData.Inputs))
	color.Green("Provide your input:")
	inputLines := readFromStdin()
	metaData, err = AddNewCaseInput(inputLines, caseName, metaData, noTimeOut)

	if err != nil {
		color.Red("Failed to add new case input")
		return err
	}

	if !context.Bool("no-output") {
		color.Green("Provide your output:")
		outputLines := readFromStdin()
		metaData, err = AddNewCaseOutput(outputLines, caseName, metaData, noTimeOut)

		if err != nil {
			color.Red("Failed to add new case output")
			return err
		}
	}

	metaData.SaveToFile(path.Join(cwd, "egor-meta.json"))

	if err != nil {
		color.Red("Failed to save to MetaData")
		return err
	}

	color.Green("Created Custom Test Case")
	return nil
}

var CaseCommand = cli.Command{
	Name:      "case",
	Aliases:   []string{"c"},
	Usage:     "Create a custom test case.",
	UsageText: "Add custom test cases to egor task.",
	Action:    CustomCaseAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-output",
			Usage: "This test case doesnt have output",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "no-timeout",
			Usage: "This test case should not timeout when passed time limit",
			Value: false,
		},
	},
}
