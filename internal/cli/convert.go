package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// runConvert writes index.tg, tokens.tg and roles.tg beside a TOML vocabulary.
// It refuses when any of them exists.
func runConvert(args []string) int {
	src := vocab.DefaultConfig
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "usage: tidgrid convert [tidgrid.toml]")
		return 2
	}
	if len(args) == 1 {
		src = args[0]
	}
	files, err := vocab.Convert(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tidgrid convert: %s: %v\n", src, err)
		return 2
	}
	dir := filepath.Dir(src)
	names := sortedNames(files)
	for _, n := range names {
		if _, err := os.Stat(filepath.Join(dir, n)); err == nil {
			fmt.Fprintf(os.Stderr, "tidgrid convert: %s exists; nothing written\n", filepath.Join(dir, n))
			return 2
		}
	}
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), []byte(files[n]), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "tidgrid convert:", err)
			return 2
		}
		fmt.Println("wrote", filepath.Join(dir, n))
	}
	fmt.Println("note: TOML comments are not carried over")
	return 0
}

func sortedNames(files map[string]string) []string {
	var names []string
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
