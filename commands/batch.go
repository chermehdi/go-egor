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

	"github.com/urfave/cli/v2"
)

const (
	RandName  = "rand.h"
	BruteName = "main_brute.cpp"
	GenName   = "gen.cpp"
)

// This is usefull to mock the implementation easily in tests.
type runnerProvider func(string) (utils.CodeRunner, bool)

func writeRandLib() error {
	// Create the rand.h file
	if err := ioutil.WriteFile(RandName, []byte(templates.RandH), 0755); err != nil {
		return err
	}
	return nil
}

func writeBruteTpl() error {
	if err := ioutil.WriteFile(BruteName, []byte(templates.BruteH), 0755); err != nil {
		return err
	}
	return nil
}

// CreateBatch will create the batch directory structure with the files needed
func CreateBatch(context *cli.Context) error {
	conf, err := config.LoadDefaultConfiguration()
	if err != nil {
		return err
	}

	egorMeta, err := config.GetMetadata()
	if err != nil {
		return err
	}

	if egorMeta.HasBatch() {
		color.Red("The task already has a batch file, Aborting...")
		return nil
	}

	egorMeta.BatchFile = GenName
	genFile, err := config.CreateFile(GenName)

	if err != nil {
		return nil
	}
	if err := templates.Compile(templates.GeneratorTemplate, genFile, conf); err != nil {
		return err
	}

	if err := writeRandLib(); err != nil {
		return nil
	}

	if err := writeBruteTpl(); err != nil {
		return nil
	}

	if err = egorMeta.SaveToFile(conf.ConfigFileName); err != nil {
		return err
	}
	color.Green("Batch creation completed successfully")
	return nil
}

func runBatchInternal(n int, egorMeta *config.EgorMeta, provider runnerProvider) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if !egorMeta.HasBatch() {
		color.Red("The current task does not have a batch setting, did you forget to run `egor batch create`?")
		return nil
	}
	runner, found := provider(egorMeta.TaskLang)
	// This used to run the generator
	if !found {
		color.Red("No task runner found for the task's default language, make sure to not change the egor metafile unless manually unless you know what you are doing!")
		return errors.New("Task runner not found!")
	}

	// This is safe because a cpp runner will always be found.
	cppr, _ := provider("cpp")

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
	panic(fmt.Sprintf("Unknown language '%s'", lang))
}

// RunBatch will delegate to the internal function to do the actual batch
// running and it will handle clean up afterwards (delete residual binary
// files).
func RunBatch(context *cli.Context) error {
	provider := func(name string) (utils.CodeRunner, bool) {
		return utils.CreateRunner(name)
	}
	egorMeta, err := config.GetMetadata()
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Remove(path.Join(cwd, "sol"))
	defer os.Remove(path.Join(cwd, "brute"))
	defer os.Remove(path.Join(cwd, "gen"))
	n := context.Int("tests")
	return runBatchInternal(n, egorMeta, provider)
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
