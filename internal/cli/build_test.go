package cli

import (
	"os"
	"path/filepath"
	"testing"
)

const base = `
[project]
out = "o.css"
[breakpoints]
md = "64em"
[tokens.color]
ink = "#000"
[tokens.space]
sm = "1rem"
`

func TestBuildCmd(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "t.toml")
	os.WriteFile(cfg, []byte(base), 0o644)
	if Run([]string{"build", cfg}) != 0 {
		t.Fatal("build failed")
	}
	a, _ := os.ReadFile(filepath.Join(dir, "o.css"))
	Run([]string{"build", cfg})
	b, _ := os.ReadFile(filepath.Join(dir, "o.css"))
	if len(a) == 0 || string(a) != string(b) {
		t.Fatal("not deterministic")
	}
	if Run([]string{"build", filepath.Join(dir, "missing.toml")}) != 2 {
		t.Fatal("missing file must exit 2")
	}
	if Run([]string{"check"}) != 2 {
		t.Fatal("check without a vocabulary in the directory must exit 2")
	}
}
