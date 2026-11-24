package cli

import (
	"errors"
	"fmt"
	"io/fs"
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
	files, errs := vocab.Convert(src)
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "%s: %v\n", src, err)
	}
	if len(errs) > 0 {
		return 2
	}
	written, err := writeNew(filepath.Dir(src), files)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tidgrid convert:", err)
		return 2
	}
	for _, p := range written {
		fmt.Println("wrote", p)
	}
	fmt.Println("note: TOML comments are not carried over")
	return 0
}

// writeNew creates each file in dir exclusively, so an existing file is never touched.
// When one fails, the files this call created are removed again.
func writeNew(dir string, files map[string]string) ([]string, error) {
	var written []string
	for _, n := range sortedNames(files) {
		p := filepath.Join(dir, n)
		err := createFile(p, files[n])
		if err != nil {
			for _, w := range written {
				os.Remove(w)
			}
			if errors.Is(err, fs.ErrExist) {
				return nil, fmt.Errorf("%s exists; nothing written", p)
			}
			return nil, err
		}
		written = append(written, p)
	}
	return written, nil
}

func createFile(p, text string) error {
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(text)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(p)
	}
	return err
}

func sortedNames(files map[string]string) []string {
	var names []string
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
