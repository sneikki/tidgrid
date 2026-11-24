package vocab

import (
	"path/filepath"
	"strings"
	"testing"
)

const tidyBase = `project(out = "o.css")
tokens("color", {"ink": "#111"})
element("body", css = {"color": "{color.ink}"})
`

func renderTidy(t *testing.T, call string) string {
	t.Helper()
	dir := writeTree(t, map[string]string{"index.tg": tidyBase + call})
	v, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	return v.Render()
}

func TestTidinessUnsetRendersNoReset(t *testing.T) {
	out := renderTidy(t, "")
	if strings.Contains(out, "box-sizing") || strings.Contains(out, "*,") {
		t.Fatalf("reset written without tidiness():\n%s", out)
	}
}

func TestTidinessMinimalSitsBetweenRootAndElements(t *testing.T) {
	out := renderTidy(t, `tidiness("minimal")`)
	root, reset, body := strings.Index(out, ":root"), strings.Index(out, "box-sizing: border-box"), strings.Index(out, "\nbody {\n  color")
	if !(root >= 0 && root < reset && reset < body) {
		t.Fatalf("order wrong: root %d reset %d body %d", root, reset, body)
	}
	if strings.Contains(out, "list-style") {
		t.Fatal("minimal wrote full's rules")
	}
}

func TestTidinessFullWritesMinimalThenFull(t *testing.T) {
	out := renderTidy(t, `tidiness("full")`)
	min, full, body := strings.Index(out, "-webkit-text-size-adjust"), strings.Index(out, "list-style"), strings.Index(out, "\nbody {\n  color")
	if !(min >= 0 && min < full && full < body) {
		t.Fatalf("order wrong: minimal %d full %d body %d", min, full, body)
	}
}

func TestTidinessWritesNoCommentsAndFormatsNested(t *testing.T) {
	out := renderTidy(t, `tidiness("full")`)
	if strings.Contains(out, "/*") && strings.Count(out, "/*") > 1 {
		t.Fatal("comments written")
	}
	if strings.Contains(out, "\n\n\n") {
		t.Fatal("blank run written")
	}
	if !strings.Contains(out, "\n    color: color-mix(in oklab, currentcolor 50%, transparent);\n") {
		t.Fatal("nested @supports rule not indented like other rules")
	}
}

func TestTidinessBadCallsAreLoadErrors(t *testing.T) {
	for name, call := range map[string]string{
		"unknown level": `tidiness("huge")`,
		"not a string":  `tidiness(1)`,
		"twice":         "tidiness(\"full\")\ntidiness(\"minimal\")",
	} {
		msg := loadErr(t, map[string]string{"index.tg": tidyBase + call})
		if !strings.Contains(msg, "index.tg:") {
			t.Errorf("%s: no file:line in %q", name, msg)
		}
		if name != "twice" && !(strings.Contains(msg, "minimal") && strings.Contains(msg, "full")) {
			t.Errorf("%s: levels not named in %q", name, msg)
		}
	}
}
