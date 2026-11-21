package cli

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/sneikki/tidgrid/internal/vocab"
)

//go:embed starter.toml
var starter []byte

// runInit writes the starter vocabulary. It never overwrites: O_EXCL makes the
// existence check and the create one step.
func runInit(args []string) int {
	fset := flag.NewFlagSet("init", flag.ContinueOnError)
	cfg := fset.String("c", vocab.DefaultConfig, "vocabulary file to create")
	if err := fset.Parse(args); err != nil {
		return 2
	}
	f, err := os.OpenFile(*cfg, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		fmt.Fprintf(os.Stderr, "tidgrid init: %s exists; not overwritten\n", *cfg)
		return 2
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "tidgrid init: %v\n", err)
		return 2
	}
	if _, err := f.Write(starter); err != nil {
		f.Close()
		fmt.Fprintf(os.Stderr, "tidgrid init: %v\n", err)
		return 2
	}
	if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "tidgrid init: %v\n", err)
		return 2
	}
	fmt.Fprintf(os.Stdout, "tidgrid init: wrote %s\n", *cfg)
	return 0
}
