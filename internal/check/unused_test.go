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
	dir := t.TempDir()
	vb, err := os.ReadFile("testdata/check/tidgrid.toml")
	if err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(dir, "tidgrid.toml")
	if err := os.WriteFile(cfg, vb, 0o644); err != nil {
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

func TestRoleOnlyInHelperCallUsed(t *testing.T) {
	code, out := runSources(t, map[string]string{
		"src/index.html": `<body class="tg-state(ready)"><p class="tg-surface(card)">Hi</p></body>`,
		"src/menu.tsx":   `export const M = () => <p className={cn("tg-text(title)")} />;`,
	})
	if code != 0 {
		t.Fatalf("exit %d, want 0\n%s", code, out)
	}
}
