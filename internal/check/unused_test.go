package check

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// runSources writes the shared vocabulary, its built stylesheet and the given
// source files into a fresh project, then runs the check on it.
func runSources(t *testing.T, files map[string]string) (int, string) {
	t.Helper()
	vb, err := os.ReadFile("testdata/check/tidgrid.toml")
	if err != nil {
		t.Fatal(err)
	}
	return runVocabSources(t, string(vb), files)
}

// runVocabSources is runSources with the vocabulary text given.
func runVocabSources(t *testing.T, vocabText string, files map[string]string) (int, string) {
	t.Helper()
	dir := t.TempDir()
	cfg := filepath.Join(dir, "tidgrid.toml")
	if err := os.WriteFile(cfg, []byte(vocabText), 0o644); err != nil {
		t.Fatal(err)
	}
	v, errs := vocab.Load(cfg)
	if errs != nil {
		t.Fatal(errs)
	}
	if err := os.MkdirAll(filepath.Dir(v.Project.Out), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(v.Project.Out, []byte(v.Render()), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		p := filepath.Join(dir, name)
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var buf bytes.Buffer
	code := reportProject(v, &buf)
	return code, buf.String()
}

const allRolesHTML = `<body class="tg-state(ready)"><p class="tg-surface(card) tg-text(title)">Hi</p></body>`

func TestUnusedRoleExit1(t *testing.T) {
	code, out := runSources(t, map[string]string{
		"src/index.html": `<body class="tg-state(ready)"><p class="tg-surface(card)">Hi</p></body>`,
	})
	if code != 1 {
		t.Fatalf("exit %d, want 1\n%s", code, out)
	}
	want := "tidgrid.toml:1: role tg-text(title) is not used in the scanned sources: use it, remove it, " +
		"or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()"
	if !strings.Contains(out, want+"\n") {
		t.Errorf("message lacks the role and its ways out\n%s", out)
	}
}

func TestAllRolesUsedExit0(t *testing.T) {
	code, out := runSources(t, map[string]string{"src/index.html": allRolesHTML})
	if code != 0 {
		t.Fatalf("exit %d, want 0\n%s", code, out)
	}
}

// A role named only in another role's selector is not a use of it.
func TestRoleOnlyInSelectorUnused(t *testing.T) {
	const vocabText = `[project]
out = "out/tidgrid.css"
scan = ["src"]

[tokens.color]
ink = "#111111"

[role.state.ready]
use = "Set by a script."

[role.surface.card]
use = "A panel."

[role.surface.card.".tg-state(ready) &"]
color = "{color.ink}"
`
	code, out := runVocabSources(t, vocabText, map[string]string{
		"src/index.html": `<body><p class="tg-surface(card)">Hi</p></body>`,
	})
	if code != 1 || !strings.Contains(out, "role tg-state(ready) is not used in the scanned sources") {
		t.Fatalf("exit %d, want 1 with tg-state(ready) unused\n%s", code, out)
	}
}

// A role used only between tidgrid-ignore-start and tidgrid-ignore-end is unused.
func TestRoleOnlyInIgnoredBlockUnused(t *testing.T) {
	code, out := runSources(t, map[string]string{
		"src/index.html": `<body class="tg-state(ready)"><p class="tg-surface(card)">Hi</p>` +
			`<!-- tidgrid-ignore-start --><p class="tg-text(title)">x</p><!-- tidgrid-ignore-end --></body>`,
	})
	if code != 1 || !strings.Contains(out, "role tg-text(title) is not used in the scanned sources") {
		t.Fatalf("exit %d, want 1 with tg-text(title) unused\n%s", code, out)
	}
}

func TestRoleOnlyInHelperCallUsed(t *testing.T) {
	code, out := runSources(t, map[string]string{
		"src/index.html": `<body class="tg-state(ready)"><p class="tg-surface(card)">Hi</p></body>`,
		"src/menu.tsx":   `export const M = () => <p className={cn("tg-text(title)")} />;`,
	})
	if code != 0 {
		t.Fatalf("exit %d, want 0\n%s", code, out)
	}
}

// A .tg role that nothing uses is reported at its own call, not at line 1 of the entry file.
func TestUnusedRoleStarlarkPosition(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"index.tg": "project(out = \"o.css\")\nsource(\"roles.tg\")\n",
		"roles.tg": "\n\nrole(\"text\", \"title\", use = \"x\")\n",
	}
	for n, s := range files {
		if err := os.WriteFile(filepath.Join(dir, n), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	v, errs := vocab.Load(filepath.Join(dir, "index.tg"))
	if errs != nil {
		t.Fatal(errs)
	}
	os.WriteFile(v.Project.Out, []byte(v.Render()), 0o644)
	var buf bytes.Buffer
	if code := reportProject(v, &buf); code != 1 || !strings.Contains(buf.String(), "roles.tg:3") {
		t.Fatalf("exit %d, want 1 at roles.tg:3\n%s", code, buf.String())
	}
}
