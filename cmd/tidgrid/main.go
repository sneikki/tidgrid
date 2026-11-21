// tidgrid: a closed CSS vocabulary for coding agents. See README.md.
package main

import (
	"os"

	"github.com/sneikki/tidgrid/internal/cli"
)

func main() { os.Exit(cli.Run(os.Args[1:])) }
