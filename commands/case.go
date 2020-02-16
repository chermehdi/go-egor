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

func AddNewCaseInput(input_lines []string,
	case_name string,
	meta_data config.EgorMeta) (config.EgorMeta, error) {

	input_file_name := case_name + ".in"
	writeLinesToFile("inputs/"+input_file_name, input_lines)
	input_file := config.NewIoFile(input_file_name, "inputs/"+input_file_name, true)
	meta_data.Inputs = append(meta_data.Inputs, input_file)

	return meta_data, nil
}

func AddNewCaseOutput(output_lines []string,
	case_name string,
	meta_data config.EgorMeta) (config.EgorMeta, error) {

	output_file_name := case_name + ".ans"
	writeLinesToFile("outputs/"+output_file_name, output_lines)
	output_file := config.NewIoFile(output_file_name, "outputs/"+output_file_name, true)
	meta_data.Outputs = append(meta_data.Outputs, output_file)

	return meta_data, nil
}

// TODO(Eroui): add checks on errors
func CustomCaseAction(context *cli.Context) error {
	color.Green("Creating Custom Test Case...")

	// Load meta data
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	meta_data, err := config.LoadMetaFromPath(path.Join(cwd, "egor-meta.json"))
	if err != nil {
		return err
	}

	case_name := "test-" + strconv.Itoa(len(meta_data.Inputs))
	color.Green("Provide your input:")
	input_lines := readFromStdin()
	meta_data, err = AddNewCaseInput(input_lines, case_name, meta_data)

	if !context.Bool("no-output") {
		color.Green("Provide your output:")
		output_lines := readFromStdin()
		meta_data, err = AddNewCaseOutput(output_lines, case_name, meta_data)
	}

	meta_data.SaveToFile(path.Join(cwd, "egor-meta.json"))

	if err != nil {
		color.Red("Failed to Generate Custom Case")
		return err
	}

	color.Green("Created Custom Test Case...")
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
	},
}
