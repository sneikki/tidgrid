// Package vocabtest is the example vocabulary the tests of other packages load.
package vocabtest

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

//go:embed example.toml
var exampleTOML []byte

// Example loads the example vocabulary. The file is written to a temporary
// directory first, because vocab.Load reads from a path.
func Example(t testing.TB) *vocab.Vocabulary {
	t.Helper()
	path := filepath.Join(t.TempDir(), vocab.DefaultConfig)
	if err := os.WriteFile(path, exampleTOML, 0o644); err != nil {
		t.Fatal(err)
	}
	v, errs := vocab.Load(path)
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	return v
}
