package cli

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sneikki/tidgrid/internal/vocab"
)

//go:embed starter.tg
var starter []byte

// runInit writes the starter vocabulary to the file named by args, default index.tg.
func runInit(args []string) int {
	fset := flag.NewFlagSet("init", flag.ContinueOnError)
	if err := fset.Parse(args); err != nil {
		return 2
	}
	if fset.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "usage: tidgrid init [file]")
		return 2
	}
	path := vocab.DefaultEntry
	if fset.NArg() == 1 {
		path = fset.Arg(0)
	}
	return writeStarter(path)
}

// writeStarter never overwrites: O_EXCL makes the existence check and the create one step.
func writeStarter(path string) int {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "tidgrid init: %v\n", err)
		return 2
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		fmt.Fprintf(os.Stderr, "tidgrid init: %s exists; not overwritten\n", path)
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
	fmt.Fprintf(os.Stdout, "tidgrid init: wrote %s\n", path)
	return 0
}
