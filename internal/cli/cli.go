// Package cli dispatches the tidgrid commands.
package cli

import (
	"fmt"
	"os"

	"github.com/sneikki/tidgrid/internal/check"
	"github.com/sneikki/tidgrid/internal/list"
	"github.com/sneikki/tidgrid/internal/plan"
	"github.com/sneikki/tidgrid/internal/sheet"
)

// Run executes the command named by args[0] and returns its exit code.
func Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: tidgrid init|convert|build|check|list|plan|sheet [file]")
		return 2
	}
	switch args[0] {
	case "convert":
		return runConvert(args[1:])
	case "init":
		return runInit(args[1:])
	case "build":
		return runBuild(args[1:])
	case "check":
		return check.Run(args[1:])
	case "list":
		return list.Run(args[1:])
	case "plan":
		return plan.Run(args[1:])
	case "sheet":
		return sheet.Run(args[1:])
	}
	fmt.Fprintf(os.Stderr, "tidgrid: unknown command %q\n", args[0])
	return 2
}
