package main

import (
	"bytes"
	"github.com/chermehdi/egor/commands"
	"github.com/chermehdi/egor/config"
	"github.com/urfave/cli/v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

const Egor = `
     ______
    / ____/___  _____ _____
   / __/ / __  / __ \/ ___/
  / /___/ /_/ / /_/ / /
 /_____/\__, /\____/_/
       /____/            version: {{ .Version }}
|------------------------------------>>
`

// returns true if `--dev` flag has been supplied
func isDev() bool {
	for _, v := range os.Args {
		if v == "--dev" {
			return true
		}
	}
	return false
}

func main() {
	var egor bytes.Buffer
	configuration, err := config.LoadDefaultConfiguration()

	if err != nil {
		log.Fatal(err)
	}

	egorTemplate, _ := template.New("egor").Parse(Egor)
	err = egorTemplate.Execute(&egor, configuration)

	if err != nil {
		log.Fatal(err)
	}

	// Disable logging if not in trace mode.
	if !isDev() {
		log.SetOutput(ioutil.Discard)
	}

	app := &cli.App{
		Name:        egor.String(),
		Usage:       "Run egor -help to print usage",
		Description: "Competitive programming helper CLI",
		UsageText:   "Run egor <subcommand> [--flags]*",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dev",
				Usage: "Set to true if you want detailed logs, useful in dev mode.",
				Value: false,
			},
		},
		Commands: []*cli.Command{
			&commands.ParseCommand,
			&commands.ConfigCommand,
			&commands.CaseCommand,
			&commands.ShowCasesCommand,
			&commands.CopyCommand,
			&commands.PrintCaseCommand,
			&commands.TestCommand,
			&commands.CreateTaskCommand,
			&commands.BatchCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
