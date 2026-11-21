// tidgrid: a closed CSS vocabulary for coding agents. See README.md.
package main

import (
	"fmt"
	"os"
)

const defaultConfig = "tidgrid.toml"

func main() { os.Exit(run(os.Args[1:])) }

func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: tidgrid init|build|check|list|sheet [-c file]")
		return 2
	}
	switch args[0] {
	case "init":
		return runInit(args[1:])
	case "build":
		return runBuild(args[1:])
	case "check":
		return runCheck(args[1:])
	case "list":
		return runList(args[1:])
	case "sheet":
		return runSheet(args[1:])
	}
	fmt.Fprintf(os.Stderr, "tidgrid: unknown command %q\n", args[0])
	return 2
}
