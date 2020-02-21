package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/chermehdi/egor/config"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	AC int8 = 0
	SK int8 = 1
	RE int8 = 2
	WA int8 = 3
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
	// Compare the trimmed output from both input and output
	if strings.Trim(got, " \t\n\r") != strings.Trim(expected, " \t\n\r") {
		return errors.New(fmt.Sprintf("Checker failed, expected %s, found %s", got, expected))
	}
	return nil
}

// Case description contains minimum information required to run one test case.
type CaseDescription struct {
	InputFile  string
	OutputFile string
	WorkFile   string
	CustomCase bool
}

func getWorkFile(fileName string) string {
	fileNameParts := strings.Split(fileName, ".")
	fileNameNoExtension := fileNameParts[0]
	return fmt.Sprintf("%s-ex.out", fileNameNoExtension)
}

// Creates a new CaseDescription from a pair of input and output IoFiles
func NewCaseDescription(input, output config.IoFile) *CaseDescription {
	base, file := filepath.Split(input.Path)
	workFilePath := path.Join(base, getWorkFile(file))
	return &CaseDescription{
		InputFile:  input.Path,
		OutputFile: output.Path,
		WorkFile:   workFilePath,
		CustomCase: input.Custom,
	}
}

// Report the execution status for a given testcase.
// Type stores also checker response
type CaseStatus struct {
	Status       int8
	CheckerError error
	Stderr       string
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

// This represents the result of running the testcases of a given task
type JudgeReport interface {
	Add(status CaseStatus)

	Display() string
}

type ConsoleJudgeReport struct {
	Stats []CaseStatus
}

func (c *ConsoleJudgeReport) Add(status CaseStatus) {
	c.Stats = append(c.Stats, status)
}

func (c *ConsoleJudgeReport) Display() string {
	panic("implement me")
}

func newJudgeReport() JudgeReport {
	return &ConsoleJudgeReport{Stats: []CaseStatus{}}
}

//
// Java judge
//
type JavaJudge struct {
	Meta           config.EgorMeta
	CurrentWorkDir string
	Checker        Checker
}

func (judge *JavaJudge) Setup() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	workDirPath := path.Join(currentDir, "work")
	if err = os.Mkdir(workDirPath, 777); err != nil {
		return err
	}
	//TODO(chermehdi): make the executables path configurable #14

	// Compilation for Java
	cmd := exec.Command("javac", judge.Meta.TaskFile, "-d", workDirPath)
	if err = cmd.Run(); err != nil {
		return err
	}
	judge.CurrentWorkDir = workDirPath
	return nil
}

func (judge *JavaJudge) RunTestCase(desc CaseDescription) CaseStatus {
	// We suppose that all java executables will be called Main
	execFilePath := path.Join(judge.CurrentWorkDir, "Main")
	cmd := exec.Command("java", execFilePath)
	cmd.Dir = judge.CurrentWorkDir
	inputFile, err := os.Open(desc.InputFile)
	if err != nil {
		return CaseStatus{
			Status:       RE,
			CheckerError: nil,
		}
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(desc.WorkFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return CaseStatus{
			Status:       RE,
			CheckerError: nil,
		}
	}
	defer outputFile.Close()

	var stderrBuffer bytes.Buffer
	cmd.Stdin = inputFile
	cmd.Stdout = bufio.NewWriter(outputFile)
	cmd.Stderr = &stderrBuffer
	if err = cmd.Run(); err != nil {
		return CaseStatus{
			Status:       RE,
			CheckerError: nil,
			Stderr:       stderrBuffer.String(),
		}

	}

	expectedOutput, err := ioutil.ReadFile(desc.OutputFile)
	if err != nil {
		return CaseStatus{
			Status:       RE,
			CheckerError: nil,
		}
	}
	output, err := ioutil.ReadFile(desc.WorkFile)
	if err != nil {
		return CaseStatus{
			Status:       RE,
			CheckerError: nil,
		}
	}
	err = judge.Checker.Check(string(output), string(expectedOutput))
	if err != nil {
		return CaseStatus{
			Status:       WA,
			CheckerError: err,
			Stderr:       stderrBuffer.String(),
		}
	}
	return CaseStatus{
		Status:       AC,
		CheckerError: err,
		Stderr:       stderrBuffer.String(),
	}
}

func (judge *JavaJudge) Cleanup() error {
	return os.RemoveAll(judge.CurrentWorkDir)
}

//
// C / Cpp Judge
//
type CppJudge struct {
	Meta config.EgorMeta
}

func (judge *CppJudge) Setup() error {
	return nil
}

func (judge *CppJudge) RunTestCase(description CaseDescription) CaseStatus {
	return CaseStatus{}
}
func (judge *CppJudge) Cleanup() error {
	return nil
}

//
// Python Judge
//
type PythonJudge struct {
	Meta config.EgorMeta
}

func (judge *PythonJudge) Setup() error {
	panic("implement me")
}

func (judge *PythonJudge) RunTestCase(CaseDescription) CaseStatus {
	panic("implement me")
}

func (judge *PythonJudge) Cleanup() error {
	panic("implement me")
}

// Creates and returns a Judge implementation corresponding to the given language
func NewJudgeFor(meta config.EgorMeta) (Judge, error) {
	switch meta.TaskLang {
	case "java":
		return &JavaJudge{Meta: meta}, nil
	case "cpp":
	case "c":
		return &CppJudge{Meta: meta}, nil
	case "python":
		return &PythonJudge{Meta: meta}, nil
	}
	return nil, errors.New(fmt.Sprintf("Cannot find judge for the given lang %s", meta.TaskLang))
}

func RunAction(_ *cli.Context) error {
	configuration, err := config.LoadDefaultConfiguration()
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	egorMetaFile := path.Join(cwd, configuration.ConfigFileName)
	egorMeta, err := config.LoadMetaFromPath(egorMetaFile)
	if err != nil {
		return err
	}

	judge, err := NewJudgeFor(egorMeta)
	if err != nil {
		return err
	}
	if err := judge.Setup(); err != nil {
		return err
	}
	defer judge.Cleanup()
	report := newJudgeReport()
	for i, input := range egorMeta.Inputs {
		output := egorMeta.Outputs[i]
		caseDescription := NewCaseDescription(input, output)
		status := judge.RunTestCase(*caseDescription)
		report.Add(status)
	}
}

var TestCommand = cli.Command{
	Name:      "test",
	Aliases:   []string{"r"},
	Usage:     "Run the current task testcases",
	UsageText: "Run the current task testcases",
	Action:    RunAction,
}
