package vocab

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTree writes files (slash path -> text) under a new directory and returns it.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, src := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func loadErr(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := writeTree(t, files)
	_, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) == 0 {
		t.Fatal("want an error")
	}
	var msgs []string
	for _, e := range errs {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "\n")
}

func TestStarlarkMatchesTOML(t *testing.T) {
	want, _ := os.ReadFile("vocabtest/example.css")
	v, errs := Load("vocabtest/example.tg")
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if v.Render() != string(want) {
		t.Fatal("example.tg renders differently from example.css")
	}
}

const projectLine = `project(out = "o.css")` + "\n"

func TestErrorNamesDeclaringCall(t *testing.T) {
	got := loadErr(t, map[string]string{
		"index.tg": projectLine + `source("roles.tg")`,
		"roles.tg": "\n\nrole(\"surface\", \"badge\", use = \"x\", css = {\"padding\": \"12px\"})\n",
	})
	want := "roles.tg:3: role.surface.badge.padding: raw length or time 12px (move it to tokens)"
	if !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestLoadRules(t *testing.T) {
	cases := map[string]map[string]string{
		"leaves the directory": {"index.tg": projectLine + `source("../x.tg")`},
		"must be relative":     {"index.tg": projectLine + `source("/etc/x.tg")`},
		"load cycle: index.tg -> a.tg -> index.tg": {
			"index.tg": projectLine + `source("a.tg")`,
			"a.tg":     `source("index.tg")`,
		},
		"declared twice (index.tg:2 and roles.tg:1)": {
			"index.tg": projectLine + `role("text", "a", use = "x")` + "\n" + `source("roles.tg")`,
			"roles.tg": `role("text", "a", use = "y")`,
		},
		"project: declared twice": {"index.tg": projectLine + projectLine},
		"has no token \"nope\"":   {"index.tg": projectLine + "c = tokens(\"color\", ink = \"#000\")\nx = c.nope"},
	}
	for want, files := range cases {
		if got := loadErr(t, files); !strings.Contains(got, want) {
			t.Errorf("got %q, want it to contain %q", got, want)
		}
	}
}

func TestSymlinkOutOfTreeRefused(t *testing.T) {
	outside := writeTree(t, map[string]string{"evil.tg": `role("text", "evil", use = "x")`})
	dir := writeTree(t, map[string]string{"index.tg": projectLine + `source("link.tg")`})
	if err := os.Symlink(filepath.Join(outside, "evil.tg"), filepath.Join(dir, "link.tg")); err != nil {
		t.Skip(err)
	}
	_, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), "cannot read inside the entry file's directory") {
		t.Fatalf("want a refusal, got %v", errs)
	}
}

func TestLoadOnceAndRelative(t *testing.T) {
	dir := writeTree(t, map[string]string{
		"index.tg": projectLine + "source(\"sub/a.tg\")\nsource(\"sub/b.tg\")\n",
		"sub/a.tg": "load(\"c.tg\", \"x\")\nrole(\"text\", \"a\", use = x)\n",
		"sub/b.tg": "load(\"c.tg\", \"x\")\n",
		"sub/c.tg": "x = \"shared\"\nrole(\"text\", \"c\", use = \"once\")\n",
	})
	v, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if len(v.Roles) != 2 {
		t.Fatalf("want 2 roles (c.tg once), got %d", len(v.Roles))
	}
}
