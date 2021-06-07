package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/chermehdi/egor/config"
)

type ExecutionResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

type ExecutionContext struct {
	// FileName is the name of the source file.
	FileName string

	// BinaryName is the name of the binary that is going to be produced after the
	// compilation step.
	BinaryName string

	// Dir is the CWD (current working directory) of the execution.
	Dir string

	// Input is the Stdin that will be fed to the executing command.
	// If the command running does not expect an stdin, this can be nil.
	Input *bytes.Buffer

	// Additional arguments passed to the command
	Args string
}

type CodeRunner interface {
	Compile(*ExecutionContext) (*ExecutionResult, error)
	Run(*ExecutionContext) (*ExecutionResult, error)
}

type CppRunner struct {
}

func (r *CppRunner) Compile(context *ExecutionContext) (*ExecutionResult, error) {
	return execute(context, "g++", "--std=c++17", context.FileName, "-o", context.BinaryName, "-Wall", "-Wextra",
		"-Wshadow", "-D_GLIBCXX_DEBUG", "-D_GLIBCXX_DEBUG_PEDANTIC")
}

func (r *CppRunner) Run(context *ExecutionContext) (*ExecutionResult, error) {
	return execute(context, fmt.Sprintf("./%s", context.BinaryName), context.Args)
}

type JavaRunner struct {
}

func (r *JavaRunner) Compile(context *ExecutionContext) (*ExecutionResult, error) {
	return execute(context, "javac", context.FileName)
}

func (r *JavaRunner) Run(context *ExecutionContext) (*ExecutionResult, error) {
	// The binary name in this case would be the name of the class containing the
	// main method. i.e if your solution file is Main.java, this should be Main
	return execute(context, "java", context.BinaryName)
}

type PythonRunner struct {
}

func (r *PythonRunner) Compile(context *ExecutionContext) (*ExecutionResult, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	// This is pyton no compilation needed
	return &ExecutionResult{stderr, stdout}, nil
}

func (r *PythonRunner) Run(context *ExecutionContext) (*ExecutionResult, error) {
	return execute(context, "python3", context.FileName)
}

type RustRunner struct {
}

func (r *RustRunner) Compile(context *ExecutionContext) (*ExecutionResult, error) {
	return execute(context, "rustc", "--edition=2018", "-O", "-o", context.BinaryName, context.FileName)
}

func (r *RustRunner) Run(context *ExecutionContext) (*ExecutionResult, error) {
	return execute(context, fmt.Sprintf("./%s", context.BinaryName), context.Args)
}

func CreateRunner(Lang string) (CodeRunner, bool) {
	switch Lang {
	case "cpp":
		return &CppRunner{}, true
	case "c":
		return &CppRunner{}, true
	case "java":
		return &JavaRunner{}, true
	case "python":
		return &PythonRunner{}, true
	case "rust":
		return &RustRunner{}, true
	}
	return nil, false
}

// Utility function to execute a given command and insure to stop it after a timeOut (in miliseconds).
// The function returns the status of the execution, the duration of the exeuction, and an error (if any).
func ExecuteWithTimeout(cmd *exec.Cmd, timeOut float64) (int8, time.Duration, error) {
	cmd.Start()
	start := time.Now()
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	timeout := time.After(time.Duration(timeOut) * time.Millisecond)
	select {
	case <-timeout:
		elapsed := time.Since(start)
		cmd.Process.Kill()
		return config.TO, elapsed, nil
	case err := <-done:
		elapsed := time.Since(start)
		return config.OK, elapsed, err
	}
}

func execute(context *ExecutionContext, command string, args ...string) (*ExecutionResult, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	log.Printf("Running the command: %v with args: %v\n", command, args)
	cmd := exec.Command(command, args...)
	cmd.Dir = context.Dir
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if context.Input != nil {
		cmd.Stdin = context.Input
	}

	if err := cmd.Run(); err != nil {
		return &ExecutionResult{
			Stdout: stdout,
			Stderr: stderr,
		}, err
	}
	return &ExecutionResult{
		Stdout: stdout,
		Stderr: stderr,
	}, nil
}
