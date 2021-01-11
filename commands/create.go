package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chermehdi/egor/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func ReadLine(reader io.Reader) string {
	result := ""
	r := bufio.NewReader(reader)
	for true {
		readBytes, isPrefix, _ := r.ReadLine()
		result += string(readBytes)
		if !isPrefix {
			break
		}
	}
	return result
}

func askInt(reader io.Reader, question string, defaultValue float64) float64 {
	color.Green("%s (%d by default): ", question, int(defaultValue))
	input := ReadLine(reader)
	var value = -1
	n, err := fmt.Sscanf(input, "%d", &value)
	if err != nil || n == 0 || value == -1 {
		return defaultValue
	}
	return float64(value)
}

func askBool(reader io.Reader, question string, defaultValue bool) bool {
	value := "N"
	if defaultValue {
		value = "Y"
	}
	response := ask(reader, question, value)
	return response == "Y" || response == "y"
}

func ask(reader io.Reader, question, defaultValue string) string {
	color.Green("%s (%s by default): ", question, defaultValue)
	value := ReadLine(reader)
	if value == "" {
		return defaultValue
	}
	return value
}

func askEndOfFile(question string) string {
	color.Green(question)
	scanner := bufio.NewScanner(os.Stdin)
	result := ""
	for scanner.Scan() {
		result += scanner.Text()
		result += "\n"
	}
	return result
}

// Fill the given task reference according to input coming from the reader object.
func fillTaskFromQuestions(task *config.Task, reader io.Reader) error {
	task.Name = ask(reader, "Task Name?", task.Name)
	task.TimeLimit = askInt(reader, "Time limit in ms?", task.TimeLimit)
	testCases := askInt(reader, "Number of test cases?", 0)
	for i := 0; i < int(testCases); i++ {
		input := askEndOfFile(fmt.Sprintf("Enter testcase #%d input:\n", i))
		output := askEndOfFile(fmt.Sprintf("Enter testcase #%d expected output:\n", i))
		task.Tests = append(task.Tests, config.TestCase{
			Input:  input,
			Output: output,
		})
	}
	// TODO(chermehdi): This information is irrelevant now, because we don't do anything special with the languages
	// sections, make sure to remove When this is updated.
	isJava := askBool(reader, "Is your language Java (you will be asked other questions if it's a yes)? (Y/N)?",
		false)
	if isJava {
		className := ask(reader, "\nMain class name?", "Main")
		// Avoid class names with strings
		className = strings.Join(strings.Split(className, " "), "")
		task.Languages["java"] = config.LanguageDescription{
			MainClass: "Main",
			TaskClass: className,
		}
	}
	return nil
}

// CreateTaskAction Creates a task directory structure, either in a default manner (a task with a generic name, and default settings).
// Or by specifying the values in an interactive manner, if `-i` flag is specified.
func CreateTaskAction(context *cli.Context) error {
	// Default task containing default values.
	task := config.Task{
		Name:        "Random Task",
		TimeLimit:   10000,
		MemoryLimit: 256,
		Tests:       nil,
		TestType:    "single",
		Input: config.IOType{
			Type: "stdin",
		},
		Output: config.IOType{
			Type: "stdout",
		},
		Languages: map[string]config.LanguageDescription{},
	}
	configuration, err := config.LoadDefaultConfiguration()
	if err != nil {
		return err
	}
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if context.Bool("i") {
		// interactive mode
		if err = fillTaskFromQuestions(&task, os.Stdin); err != nil {
			return err
		}
	}
	location, err := CreateDirectoryStructure(task, *configuration, curDir, nil)
	if err != nil {
		return err
	}
	color.Green("Task created at %s\n", location)
	return nil
}

// CreateTaskCommand Command will ask the user for a bunch of questions about the task name
// Input and outputs and other metadata required to create the task
var CreateTaskCommand = cli.Command{
	Name:      "create",
	Aliases:   []string{"c"},
	Usage:     "Create a new task directory",
	UsageText: "egor create or egor create -i",
	Action:    CreateTaskAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "interactive",
			Aliases: []string{"i"},
			Usage:   "Create the task in an interactive form (answer questions, egor will do the rest)",
			Value:   false,
		},
	},
}
