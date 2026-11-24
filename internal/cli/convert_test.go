package cli

import (
	"os"
	"path/filepath"
	"testing"
)

const fixtures = "../vocab/vocabtest/"

func TestConvertThenBuildMatchesGolden(t *testing.T) {
	dir := t.TempDir()
	toml, err := os.ReadFile(fixtures + "example.toml")
	if err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "tidgrid.toml")
	if err := os.WriteFile(src, toml, 0o644); err != nil {
		t.Fatal(err)
	}
	if Run([]string{"convert", src}) != 0 {
		t.Fatal("convert failed")
	}
	if Run([]string{"convert", src}) != 2 {
		t.Fatal("second convert must refuse to overwrite")
	}
	want, _ := os.ReadFile(fixtures + "example.css")
	t.Chdir(dir)
	if Run([]string{"build"}) != 0 { // no argument: index.tg
		t.Fatal("build failed")
	}
	got, _ := os.ReadFile(filepath.Join(dir, "src/styles/tidgrid.css"))
	if string(got) != string(want) {
		t.Fatal("converted vocabulary builds a different stylesheet")
	}
}

func TestNoEntryWithTOMLHintsConvert(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "tidgrid.toml"), []byte(base), 0o644)
	t.Chdir(dir)
	if Run([]string{"build"}) != 2 {
		t.Fatal("build with only tidgrid.toml must exit 2")
	}
}
