package plan

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

const planBase = `
[project]
prefix = "tg"
out = "o.css"
[breakpoints]
md = "64em"
[tokens.color]
ink = "#000"
accent = "#111"
[tokens.space]
md = "1rem"
lg = "2rem"
[role.surface.card]
use = "A panel."
padding = "{space.md}"
[role.surface.card."&:hover"]
color = "{color.ink}"
`

// writeVocab writes a vocabulary and returns its path.
func writeVocab(t *testing.T, dir, text string) string {
	t.Helper()
	cfg := filepath.Join(dir, "t.toml")
	if err := os.WriteFile(cfg, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	return cfg
}

// build writes the stylesheet the vocabulary renders to, as `tidgrid build` does.
func build(t *testing.T, cfg string) {
	t.Helper()
	v, code := vocab.LoadFile(cfg)
	if v == nil {
		t.Fatalf("vocabulary does not load: exit %d", code)
	}
	if err := os.WriteFile(v.Project.Out, []byte(v.Render()), 0o644); err != nil {
		t.Fatal(err)
	}
}

// planAfter runs plan for the vocabulary text on disk after the built one.
func planAfter(t *testing.T, built, now string) ([]string, int) {
	t.Helper()
	dir := t.TempDir()
	cfg := writeVocab(t, dir, built)
	build(t, cfg)
	cfg = writeVocab(t, dir, now)
	return report([]string{cfg})
}

func wantLines(t *testing.T, got []string, code int, want ...string) {
	t.Helper()
	if code != 0 {
		t.Fatalf("exit %d, want 0", code)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestNoStylesheetListsEverythingAdded(t *testing.T) {
	dir := t.TempDir()
	lines, code := report([]string{writeVocab(t, dir, planBase)})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	for _, want := range []string{"added class tg-surface(card)", "added token color.ink: #000", "added class tg-surface(card) in &:hover"} {
		if !slices.Contains(lines, want) {
			t.Errorf("missing %q in\n%s", want, strings.Join(lines, "\n"))
		}
	}
}

func TestIdenticalSaysNoDifference(t *testing.T) {
	dir := t.TempDir()
	cfg := writeVocab(t, dir, planBase)
	build(t, cfg)
	lines, code := report([]string{cfg})
	wantLines(t, lines, code, "no difference")
}

func TestAddedRole(t *testing.T) {
	lines, code := planAfter(t, planBase, planBase+"\n[role.text.title]\nuse = \"Title.\"\nfont-weight = \"bold\"\n")
	wantLines(t, lines, code, "added class tg-text(title)")
}

func TestRemovedRole(t *testing.T) {
	lines, code := planAfter(t, planBase+"\n[role.text.title]\nuse = \"Title.\"\nfont-weight = \"bold\"\n", planBase)
	wantLines(t, lines, code, "removed class tg-text(title)")
}

func TestChangedDeclaration(t *testing.T) {
	now := strings.Replace(planBase, `padding = "{space.md}"`, `padding = "{space.lg}"`, 1)
	lines, code := planAfter(t, planBase, now)
	if code != 0 || len(lines) != 1 || !strings.HasPrefix(lines[0], "changed class tg-surface(card): padding: ") {
		t.Fatalf("exit %d, lines %q", code, lines)
	}
}

func TestChangedNestedRule(t *testing.T) {
	now := strings.Replace(planBase, `color = "{color.ink}"`, `color = "{color.accent}"`, 1)
	lines, code := planAfter(t, planBase, now)
	if code != 0 || len(lines) != 1 || !strings.HasPrefix(lines[0], "changed class tg-surface(card) in &:hover: color: ") {
		t.Fatalf("exit %d, lines %q", code, lines)
	}
}

func TestChangedToken(t *testing.T) {
	now := strings.Replace(planBase, `ink = "#000"`, `ink = "#222"`, 1)
	lines, code := planAfter(t, planBase, now)
	wantLines(t, lines, code, "changed token color.ink: #000 -> #222", "  reaches 1 classes: tg-surface(card)")
}

func TestVocabularyThatDoesNotLoadExits2(t *testing.T) {
	dir := t.TempDir()
	lines, code := report([]string{writeVocab(t, dir, "[project\nout = ")})
	if code != 2 || lines != nil {
		t.Fatalf("exit %d, lines %q; want exit 2 and no lines", code, lines)
	}
}

func TestBadStylesheetOnDiskDoesNotCrash(t *testing.T) {
	dir := t.TempDir()
	cfg := writeVocab(t, dir, planBase)
	build(t, cfg)
	garbage := []string{
		"}}}{{{",
		".a { color: red",
		"/* unclosed",
		`.a { content: "unterminated } }`,
		strings.Repeat("@media x{", 500) + ".a{b:c}",
		".a\\",
		": ; { } ;",
	}
	for _, css := range garbage {
		if err := os.WriteFile(filepath.Join(dir, "o.css"), []byte(css), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, code := report([]string{cfg}); code != 0 {
			t.Fatalf("garbage %q: exit %d, want 0", css, code)
		}
	}
}
