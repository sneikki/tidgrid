package vocab

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	// DefaultConfig is the TOML vocabulary file name; init writes it.
	DefaultConfig = "tidgrid.toml"
	// DefaultEntry is the vocabulary file the commands read when none is given.
	DefaultEntry = "index.tg"
)

// LoadFile loads the vocabulary at path. On failure it prints each error to
// stderr and returns a nil vocabulary with exit code 2; else the vocabulary and 0.
func LoadFile(path string) (*Vocabulary, int) {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "%s: no such file; run `tidgrid init` to create a vocabulary, or give the file as an argument\n", path)
		return nil, 2
	}
	v, errs := Load(path)
	if len(errs) > 0 {
		for _, e := range errs {
			if filepath.Ext(path) == ".tg" {
				fmt.Fprintln(os.Stderr, e) // already names its file and line
			} else {
				fmt.Fprintf(os.Stderr, "%s: %v\n", path, e)
			}
		}
		return nil, 2
	}
	return v, 0
}

// LoadFlag parses args as the command's flags plus an optional vocabulary file and loads it.
func LoadFlag(name string, args []string) (*Vocabulary, int) {
	return LoadFlagSet(flag.NewFlagSet(name, flag.ContinueOnError), args)
}

// LoadFlagSet is LoadFlag for a command with flags of its own.
func LoadFlagSet(set *flag.FlagSet, args []string) (*Vocabulary, int) {
	if err := set.Parse(args); err != nil {
		return nil, 2
	}
	if set.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "tidgrid %s: want at most one vocabulary file, got %d arguments\n", set.Name(), set.NArg())
		return nil, 2
	}
	path := DefaultEntry
	if set.NArg() == 1 {
		path = set.Arg(0)
	} else if msg := convertHint(); msg != "" {
		fmt.Fprintln(os.Stderr, msg)
		return nil, 2
	}
	return LoadFile(path)
}

// convertHint is the message for a directory that has only the old tidgrid.toml, else "".
func convertHint() string {
	if _, err := os.Stat(DefaultEntry); !errors.Is(err, fs.ErrNotExist) {
		return ""
	}
	if _, err := os.Stat(DefaultConfig); err != nil {
		return ""
	}
	return DefaultEntry + " not found, but " + DefaultConfig + " is here: run `tidgrid convert` to write " + DefaultEntry
}
