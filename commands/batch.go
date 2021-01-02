package commands

import (
	"os"
	"path"
	template2 "text/template"

	"github.com/chermehdi/egor/config"
	"github.com/chermehdi/egor/templates"
	"github.com/fatih/color"
	"io/ioutil"

	"github.com/urfave/cli/v2"
)

type Batcher struct {
	judge Judge
}

func (b *Batcher) Setup() error {
	b.judge.Setup()
	return nil
}

func (b *Batcher) Run() error {
	return nil
}

func (b *Batcher) Cleanup() error {
	b.judge.Cleanup()
	return nil
}

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

	genPath := path.Join(cwd, "gen.cpp")
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
	randH := path.Join(cwd, "rand.h")
	// Create the rand.h file
	if err = ioutil.WriteFile(randH, []byte(templates.RandH), 0755); err != nil {
		return err
	}
	// Update the metadata
	fileName := path.Join(cwd, configuration.ConfigFileName)
	if err = egorMeta.SaveToFile(fileName); err != nil {
		return err
	}
	color.Green("Batch creation completed successfully")
	return nil
}

func RunBatch(context *cli.Context) error {
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

	judge, err := NewJudgeFor(egorMeta, configuration, getChecker(context.String("checker")))
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
		caseDescription := NewCaseDescription(input, output, egorMeta.TimeLimit)
		status := judge.RunTestCase(*caseDescription)
		report.Add(status, *caseDescription)
	}
	report.Display()
	return nil
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
		},
	},
}
