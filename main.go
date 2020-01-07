package main

import (
	"fmt"
	"os"

	"github.com/chermehdi/egor/parse"
)

func main() {
	parser := parse.NewParser(os.Args)
	options := parser.Parse()
	fmt.Printf("%v\n", options)
}
