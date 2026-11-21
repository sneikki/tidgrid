package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// loadFlag parses -c and loads the vocabulary; on failure it prints the
// errors and returns exit code 2.
func loadFlag(name string, args []string) (*Vocabulary, int) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	cfg := fs.String("c", defaultConfig, "vocabulary file")
	if err := fs.Parse(args); err != nil {
		return nil, 2
	}
	v, errs := load(*cfg)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "%s: %v\n", *cfg, e)
		}
		return nil, 2
	}
	return v, 0
}

func runBuild(args []string) int {
	v, code := loadFlag("build", args)
	if v == nil {
		return code
	}
	out := v.Project.Out
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "tidgrid build:", err)
		return 2
	}
	if err := os.WriteFile(out, []byte(v.render()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "tidgrid build:", err)
		return 2
	}
	fmt.Println("wrote", out)
	return 0
}
