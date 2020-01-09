package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var ParseCommand cli.Command = cli.Command{
	Name:      "parse",
	Aliases:   []string{"p"},
	Usage:     "parse task from navigator",
	UsageText: "parse task from navigator",
	Action: func(context *cli.Context) error {
		fmt.Println("Hello from parse")
		return nil
	},
}
