package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

func TestStarterLoadsAndRenders(t *testing.T) {
	p := filepath.Join(t.TempDir(), "index.tg")
	if err := os.WriteFile(p, starter, 0o644); err != nil {
		t.Fatal(err)
	}
	v, errs := vocab.Load(p)
	if len(errs) > 0 {
		t.Fatalf("starter does not load: %v", errs)
	}
	if out := v.Render(); out == "" {
		t.Fatal("starter renders empty")
	}
}

func TestInitWritesNamedFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "some", "name.tg")
	if code := runInit([]string{p}); code != 0 {
		t.Fatalf("exit %d, want 0", code)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, starter) {
		t.Fatal("written file differs from starter")
	}
}

func TestInitDefaultsToIndexTg(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if code := runInit(nil); code != 0 {
		t.Fatalf("exit %d, want 0", code)
	}
	got, err := os.ReadFile(filepath.Join(dir, "index.tg"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, starter) {
		t.Fatal("index.tg differs from starter")
	}
}

func TestInitRefusesOverwrite(t *testing.T) {
	p := filepath.Join(t.TempDir(), "index.tg")
	keep := []byte("# mine\n")
	if err := os.WriteFile(p, keep, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := runInit([]string{p}); code != 2 {
		t.Fatalf("exit %d, want 2", code)
	}
	got, _ := os.ReadFile(p)
	if !bytes.Equal(got, keep) {
		t.Fatal("existing file was modified")
	}
}

func TestInitThenBuildAndCheck(t *testing.T) {
	t.Chdir(t.TempDir())
	if code := runInit(nil); code != 0 {
		t.Fatalf("init exit %d, want 0", code)
	}
	if code := Run([]string{"build"}); code != 0 {
		t.Fatalf("build exit %d, want 0", code)
	}
	if code := Run([]string{"check"}); code != 0 {
		t.Fatalf("check exit %d, want 0", code)
	}
}

// uncommentRoles strips the comment markers of the role examples: each block
// starts at a "# role(" line and runs to the next blank line.
func uncommentRoles(src []byte) []byte {
	var out []string
	inRole := false
	for _, line := range strings.Split(string(src), "\n") {
		switch {
		case strings.HasPrefix(line, "# role("):
			inRole = true
		case line == "":
			inRole = false
		}
		if inRole {
			line = strings.TrimPrefix(line, "# ")
		}
		out = append(out, line)
	}
	return []byte(strings.Join(out, "\n"))
}

func TestStarterRoleExamplesBuild(t *testing.T) {
	t.Chdir(t.TempDir())
	src := uncommentRoles(starter)
	if !bytes.Contains(src, []byte("role(\"action\", \"quiet\",")) {
		t.Fatal("no role examples found in starter")
	}
	if err := os.WriteFile("index.tg", src, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := Run([]string{"build"}); code != 0 {
		t.Fatalf("build with role examples uncommented: exit %d, want 0", code)
	}
}
