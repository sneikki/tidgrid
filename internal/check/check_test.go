package check

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
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
	v, errs := vocab.Load(cfg)
	if errs != nil {
		t.Fatal(errs)
	}
	_, ownOut := os.Stat(filepath.Join(dir, "out/tidgrid.css"))
	_, noOut := os.Stat(filepath.Join(dir, "no-out"))
	if ownOut != nil && noOut != nil {
		os.MkdirAll(filepath.Join(dir, "out"), 0o755)
		os.WriteFile(v.Project.Out, []byte(v.Render()), 0o644)
	}
	var buf bytes.Buffer
	code = reportProject(v, &buf)
	out = buf.String()
	head, _, _ := strings.Cut(out, "\n\n")
	if code == 1 {
		violations = withoutUnusedRoles(strings.Split(head, "\n"))
	}
	if len(violations) == 0 {
		code = 0
	}
	return
}

// withoutUnusedRoles drops the unused-role lines: the table cases test file
// rules against one shared vocabulary, and TestUnusedRole* test the role rule.
func withoutUnusedRoles(lines []string) []string {
	var kept []string
	for _, l := range lines {
		if !strings.Contains(l, " is not used in the scanned sources: ") {
			kept = append(kept, l)
		}
	}
	return kept
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
		"  class has-js: set by an inline script before first paint\n",
		"  css src/fonts.css: self-hosted @font-face only\n",
	} {
		if !strings.Contains(out, s) {
			t.Errorf("summary lacks %q\n%s", s, out)
		}
	}
	if strings.Contains(out, "roles unused:") {
		t.Errorf("summary repeats the unused roles\n%s", out)
	}
	_, _, out = runCase(t, "unknown-class")
	for _, r := range []string{"tg-surface(card)", "tg-text(title)", "tg-state(ready)"} {
		if n := strings.Count(out, "role "+r+" is not used"); n != 1 {
			t.Errorf("unused role %s reported %d times, want 1\n%s", r, n, out)
		}
	}
}

// A symlinked directory is not walked: check must say so, unless its target
// is scanned through another root.
func TestSymlinkedDirReported(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "site"), 0o755)
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "a.html"), []byte(`<p class="nope">`), 0o644)
	if err := os.Symlink("../docs", filepath.Join(dir, "site", "docs")); err != nil {
		t.Skip(err)
	}
	run := func(scan string) string {
		cfg := filepath.Join(dir, "site", "tidgrid.toml")
		os.WriteFile(cfg, []byte("[project]\nout = \"o.css\"\nscan = "+scan+"\n"), 0o644)
		v, errs := vocab.Load(cfg)
		if errs != nil {
			t.Fatal(errs)
		}
		os.WriteFile(v.Project.Out, []byte(v.Render()), 0o644)
		var buf bytes.Buffer
		reportProject(v, &buf)
		return buf.String()
	}
	if out := run(`["."]`); !strings.Contains(out, "docs:1: symlinked directory is not scanned") {
		t.Errorf("symlinked directory passed in silence\n%s", out)
	}
	out := run(`[".", "../docs"]`)
	if strings.Contains(out, "symlinked directory") || !strings.Contains(out, "unknown class nope") {
		t.Errorf("target under scan: want its page checked and no symlink report\n%s", out)
	}
}

func TestScriptUseCountsRole(t *testing.T) {
	_, _, out := runCase(t, "script-use")
	if !strings.Contains(out, "  src/menu.js: tg-state(ready)\n") {
		t.Errorf("script use not counted\n%s", out)
	}
}

func TestRunCheckExit2(t *testing.T) {
	if c := Run([]string{"-c", "testdata/check/nope.toml"}); c != 2 {
		t.Errorf("missing vocabulary: exit %d, want 2", c)
	}
}
