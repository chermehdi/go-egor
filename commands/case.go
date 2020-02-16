package commands 

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"fmt"
	"os"
	"github.com/chermehdi/egor/config"
	"path"
	"strconv"
)

func readFromStdin() ([]string, error) {
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
	
	// TODO(Eroui) add check for errors
	return lines, nil
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

	color.Green("Provide your input:")
	input_lines, _ := readFromStdin()

	color.Green("Provide your output:")
	output_lines, _ := readFromStdin()
	
	case_name := "test-" + strconv.Itoa(len(meta_data.Inputs))

	input_file_name := case_name + ".in"
	output_file_name := case_name + ".out"

	writeLinesToFile("inputs/" + input_file_name, input_lines)
	writeLinesToFile("outputs/" + output_file_name, output_lines)

	input_file := config.NewIoFile(input_file_name, "inputs/" + input_file_name, true)
	output_file := config.NewIoFile(input_file_name, "outputs/" + output_file_name, true)

	meta_data.Inputs = append(meta_data.Inputs, input_file)
	meta_data.Outputs = append(meta_data.Outputs, output_file)

	meta_data.SaveToFile(path.Join(cwd, "egor-meta.json"))
	
	if err != nil {
		fmt.Println(err)
	}
	
	color.Green("Created Custom Test Case...")
	return nil
}

var CaseCommand = cli.Command{
	Name:		"case",
	Aliases:	[]string{"c"},
	Usage:		"Create a custom test case.",
	UsageText:	"Add custom test cases to egor task.",
	Action: 	CustomCaseAction,
	// TODO(Eroui): add necessary flags
}
