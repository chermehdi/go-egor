package main

import (
	"bytes"
	"github.com/chermehdi/egor/commands"
	"github.com/chermehdi/egor/config"
	"github.com/urfave/cli/v2"
	"html/template"
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

var EgorTemplate, _ = template.New("egor").Parse(Egor)

func main() {
	var egor bytes.Buffer
	configuration, err := config.LoadDefaultConfiguration()

	if err != nil {
		log.Fatal(err)
	}

	err = EgorTemplate.Execute(&egor, configuration)
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:        egor.String(),
		Description: "Competitive programming helper CLI",
		UsageText:   "Run egor <subcommand> [--flags]*",
		Commands: []*cli.Command{
			&commands.ParseCommand,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
