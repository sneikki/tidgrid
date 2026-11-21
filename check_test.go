package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func copyTree(t *testing.T, from, to string) {
	t.Helper()
	err := filepath.WalkDir(from, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		r, _ := filepath.Rel(from, p)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		dst := filepath.Join(to, r)
		os.MkdirAll(filepath.Dir(dst), 0o755)
		return os.WriteFile(dst, b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// runCase copies a case next to the shared vocabulary, writes the built
// stylesheet unless the case brings its own (or says no-out), and checks.
func runCase(t *testing.T, name string) (code int, violations []string, out string) {
	t.Helper()
	dir := t.TempDir()
	copyTree(t, "testdata/check/cases/"+name, dir)
	vb, _ := os.ReadFile("testdata/check/tidgrid.toml")
	cfg := filepath.Join(dir, "tidgrid.toml")
	os.WriteFile(cfg, vb, 0o644)
	v, errs := load(cfg)
	if errs != nil {
		t.Fatal(errs)
	}
	_, ownOut := os.Stat(filepath.Join(dir, "out/tidgrid.css"))
	_, noOut := os.Stat(filepath.Join(dir, "no-out"))
	if ownOut != nil && noOut != nil {
		os.MkdirAll(filepath.Join(dir, "out"), 0o755)
		os.WriteFile(v.Project.Out, []byte(v.render()), 0o644)
	}
	var buf bytes.Buffer
	code = doCheck(v, &buf)
	out = buf.String()
	head, _, _ := strings.Cut(out, "\n\n")
	if code == 1 {
		violations = strings.Split(head, "\n")
	}
	return
}

func TestCheckCases(t *testing.T) {
	ents, err := os.ReadDir("testdata/check/cases")
	if err != nil || len(ents) == 0 {
		t.Fatalf("no check cases: %v", err)
	}
	for _, e := range ents {
		name := e.Name()
		t.Run(name, func(t *testing.T) {
			code, got, out := runCase(t, name)
			wb, err := os.ReadFile("testdata/check/cases/" + name + "/want")
			if err != nil {
				t.Fatal(err)
			}
			var want []string
			if s := strings.TrimSpace(string(wb)); s != "" {
				want = strings.Split(s, "\n")
			}
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Errorf("violations differ\n got:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}
			if wantCode := map[bool]int{true: 1, false: 0}[len(want) > 0]; code != wantCode {
				t.Errorf("exit %d, want %d\n%s", code, wantCode, out)
			}
		})
	}
}

func TestCheckSummary(t *testing.T) {
	code, _, out := runCase(t, "clean")
	if code != 0 {
		t.Fatalf("exit %d\n%s", code, out)
	}
	for _, s := range []string{
		"roles used:\n",
		"  src/b.astro: tg-surface(card)\n",
		"  src/index.html: tg-state(ready) tg-surface(card) tg-text(title)\n",
		"roles unused:\n  none\n",
		"  class has-js: set by an inline script before first paint\n",
		"  css src/fonts.css: self-hosted @font-face only\n",
	} {
		if !strings.Contains(out, s) {
			t.Errorf("summary lacks %q\n%s", s, out)
		}
	}
	_, _, out = runCase(t, "unknown-class")
	if !strings.Contains(out, "roles unused:\n  tg-surface(card)\n  tg-text(title)\n  tg-state(ready)\n") {
		t.Errorf("unused roles not listed\n%s", out)
	}
}

func TestScriptUseCountsRole(t *testing.T) {
	_, _, out := runCase(t, "script-use")
	if !strings.Contains(out, "  src/menu.js: tg-state(ready)\n") {
		t.Errorf("script use not counted\n%s", out)
	}
}

func TestRunCheckExit2(t *testing.T) {
	if c := runCheck([]string{"-c", "testdata/check/nope.toml"}); c != 2 {
		t.Errorf("missing vocabulary: exit %d, want 2", c)
	}
}
