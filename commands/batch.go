package commands

import (
	"errors"
	"fmt"
	"github.com/chermehdi/egor/config"
	"github.com/chermehdi/egor/templates"
	"github.com/chermehdi/egor/utils"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	template2 "text/template"

	"github.com/urfave/cli/v2"
)

const (
	RandName  = "rand.h"
	BruteName = "main_brute.cpp"
	GenName   = "gen.cpp"
)

// CreateBatch will create the batch directory structure with the files needed
func CreateBatch(context *cli.Context) error {
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

	if egorMeta.HasBatch() {
		color.Red("The task already has a batch file, Aborting...")
		return nil
	}

	genPath := path.Join(cwd, GenName)
	egorMeta.BatchFile = genPath
	// Create the generator file

	genTemp := template2.New("Solution template")
	compiledTemplate, err := genTemp.Parse(templates.GeneratorTemplate)
	if err != nil {
		return err
	}
	genFile, err := config.CreateFile(genPath)
	if err != nil {
		return nil
	}
	// TODO(chermehdi): Move this outside, it smells of code repetition
	if err = compiledTemplate.Execute(genFile, configuration); err != nil {
		return err
	}
	if context.Bool("has-brute") {
		randH := path.Join(cwd, RandName)
		// Create the rand.h file
		if err = ioutil.WriteFile(randH, []byte(templates.RandH), 0755); err != nil {
			return err
		}

		bruteH := path.Join(cwd, BruteName)
		if err = ioutil.WriteFile(bruteH, []byte(templates.BruteH), 0755); err != nil {
			return err
		}
	}
	// Update the metadata
	fileName := path.Join(cwd, configuration.ConfigFileName)
	if err = egorMeta.SaveToFile(fileName); err != nil {
		return err
	}
	color.Green("Batch creation completed successfully")
	return nil
}

func runBatchInternal(context *cli.Context) error {
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
	if !egorMeta.HasBatch() {
		color.Red("The current task does not have a batch setting, did you forget to run egor batch create?")
		return nil
	}
	n := context.Int("tests")
	runner, found := utils.CreateRunner(egorMeta.TaskLang)
	// This used to run the generator
	if !found {
		color.Red("No task runner found for the task's default language, make sure to not change the egor metafile unless manually unless you know what you are doing!")
		return errors.New("Task runner not found!")
	}
	cppr, _ := utils.CreateRunner("cpp")

	// Compile the generator
	c := &utils.ExecutionContext{
		FileName:   GenName,
		BinaryName: "gen",
		Input:      nil,
		Dir:        cwd,
	}

	log.Println("Compiling the generator...")

	res, err := cppr.Compile(c)
	if err != nil {
		if res.Stderr.String() != "" {
			color.Red(res.Stderr.String())
		}
		return err
	}

	log.Println("Finished compiling the generator")

	log.Println("Compiling the solution ...")
	binName := getBinaryName(egorMeta.TaskLang)
	// Compile the solution
	c1 := &utils.ExecutionContext{
		FileName:   egorMeta.TaskFile,
		BinaryName: binName,
		Input:      nil,
		Dir:        cwd,
	}

	res, err = runner.Compile(c1)
	if err != nil {
		if res.Stderr.String() != "" {
			color.Red(res.Stderr.String())
		}
		return err
	}

	log.Println("Finished compiling the solution")

	log.Println("Finished compiling the brute force solution")
	// Compile the brute force solution
	c2 := &utils.ExecutionContext{
		FileName:   BruteName,
		BinaryName: "brute",
		Input:      nil,
		Dir:        cwd,
	}

	res, err = cppr.Compile(c2)
	if err != nil {
		if res.Stderr.String() != "" {
			color.Red(res.Stderr.String())
		}
		return err
	}

	log.Println("Finished compiling the brute force solution")

	checker := &TokenChecker{}

	for i := 0; i < n; i = i + 1 {
		color.Green(fmt.Sprintf("Running test %d", i+1))

		log.Println("Running the generator.. ")
		// Args should contain the random seed to give generator as argument
		c.Args = fmt.Sprintf("%d", rand.Intn(int(1<<30)))

		// Get the random input
		res, err = cppr.Run(c)
		if err != nil {
			color.Red("Could not run the generator ")
			if res.Stderr.String() != "" {
				color.Red(res.Stderr.String())
			}
			return err
		}
		log.Println("Finished running the generator")
		log.Printf("Generator output: \n%s\n", res.Stdout.String())

		log.Println("Running the main solution..")

		// making two input copies:
		//   - one for the brute solution consumer
		//   - second only used in case of a found diff to print the output to the
		//   input to the user
		inC1 := res.Stdout
		inC2 := res.Stdout

		// Run the main solution
		// Feed the output of the generator as the main solution's input.
		c1.Input = &res.Stdout
		main, err := runner.Run(c1)
		if err != nil {
			color.Red("Could not run the solution")
			if main.Stderr.String() != "" {
				color.Red(main.Stderr.String())
			}
			return err
		}

		log.Println("Finished running main the solution")

		// Run the main solution
		// Feed the output of the generator as the brute solution's input.
		log.Println("Running the brute solution ...")

		c2.Input = &inC1
		brute, err := cppr.Run(c2)
		if err != nil {
			color.Red("Could not run the brute solution")
			if brute.Stderr.String() != "" {
				color.Red(brute.Stderr.String())
			}
			return err
		}

		if err := checker.Check(brute.Stdout.String(), main.Stdout.String()); err != nil {
			color.Red("Found diff: ")
			color.Red("Input: ")
			fmt.Println(inC2.String())
			color.Red("Expected: ")
			fmt.Println(brute.Stdout.String())
			color.Red("Got: ")
			fmt.Println(main.Stdout.String())
			color.Red(err.Error())
			return err
		}
	}

	color.Green("All green, not diff found - Your solution is probably correct!")
	return nil
}

// Compute the binary name according to the language.
func getBinaryName(lang string) string {
	switch lang {
	case "cpp":
		return "sol"
	case "java":
		return "Main"
	case "python":
		return "main.py"
	}
	// should not get here.
	return ""
}

// RunBatch will delegate to the internal function to do the actual batch
// running and it will handle clean up afterwards (delete residual binary
// files).
func RunBatch(context *cli.Context) error {
	err := runBatchInternal(context)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Remove(path.Join(cwd, "sol"))
	os.Remove(path.Join(cwd, "brute"))
	os.Remove(path.Join(cwd, "gen"))
	return err
}

var BatchCommand = cli.Command{
	Name:      "batch",
	Aliases:   []string{"b"},
	Usage:     "Create and Run bach tests",
	UsageText: "egor batch (create | run)",
	Subcommands: []*cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create the template for the batch test",
			Action:  CreateBatch,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "has-brute",
					Usage:   "Does this solution have a brute force solution",
					Aliases: []string{"b"},
					Value:   true,
				},
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run the batch test",
			Action:  RunBatch,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "tests",
					Usage:   "Number of tests to run in the batch round",
					Aliases: []string{"t"},
					Value:   100,
				},
			},
		},
	},
}
