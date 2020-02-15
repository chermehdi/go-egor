package commands 

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"fmt"
    "os"
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
	
	// TODO add check for errors
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
	
	color.Green("Provide your input:")
	input_lines, _ := readFromStdin()

	color.Green("Provide your output:")
	output_lines, _ := readFromStdin()

	writeLinesToFile("sample.in", input_lines)
	writeLinesToFile("sample.out", output_lines)
	
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
