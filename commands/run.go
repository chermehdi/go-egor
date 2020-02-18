package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

// Checks the output of a given testcase against it's expected output
type Checker interface {
	// Execute the check (got, expected) and returns
	// nil if the output match, otherwise an error with a description message.
	Check(string, string) error
}

// Default implementation of the Checker interface.
type DiffChecker struct {
}

func (c *DiffChecker) Check(got, expected string) error {
	// Compre the trimmed output from both input and output
	if strings.Trim(got, " \t\n\r") != strings.Trim(expected, " \t\n\r") {
		return errors.New(fmt.Sprintf("Checker failed, expected %s, found %s", got, expected))
	}
	return nil
}

// Case description contains minimum information required to run one test case.
type CaseDescription struct {
	InputFile  string
	OutputFile string
	CustomCase bool
}

// Report the execution status for a given testcase.
// Type also stores the compiler output, and the checker response
type CaseStatus struct {
}

// Implementation must be able to prepare the working environement to compile and execute testscases,
// And run each testcase and report the status back to the invoker, and perform any necessary cleanup (binaries created, directories created ...)
type Judge interface {
	// setup the working directory and perform any necessary compilation of the task
	// if the setup returned an error, the Judge should abort the operation and report the error back.
	Setup() error

	// Run on every testcase, and the status is reported back to the invoker.
	// The implementation is free to Run all testcases at once, or report every testcase execution status once it finishes.
	// If it's needed, running independent cases can be done on different go routines.
	RunTestCase(CaseDescription) CaseStatus

	// Cleanup the working directory, if an error occured, implementation must report it to the caller.
	Cleanup() error
}

func RunAction(context *cli.Context) error {
	return nil
}

var RunCommand = cli.Command{
	Name:      "run",
	Aliases:   []string{"r"},
	Usage:     "Run the current task testcases",
	UsageText: "Run the current task testcases",
	Action:    RunAction,
}
