// Package plan compares the vocabulary's stylesheet with the one on disk and
// lists what a vocabulary change would add, remove or change. It writes nothing.
package plan

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// Run is the plan command: it prints the difference and returns the exit code.
func Run(args []string) int {
	lines, code := report(args)
	for _, l := range lines {
		fmt.Println(l)
	}
	return code
}

// report is Run without the printing. A missing stylesheet counts as empty, so everything is added.
func report(args []string) ([]string, int) {
	v, code := vocab.LoadFlag("plan", args)
	if v == nil {
		return nil, code
	}
	disk, err := os.ReadFile(v.Project.Out)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintln(os.Stderr, "tidgrid plan:", err)
		return nil, 2
	}
	lines := Diff(string(disk), v.Render())
	if len(lines) == 0 {
		return []string{"no difference"}, 0
	}
	return lines, 0
}
