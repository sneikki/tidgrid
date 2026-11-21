package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSheetExample(t *testing.T) {
	v := loadExampleVocab(t)
	out := sheetHTML(v)
	if again := sheetHTML(v); again != out {
		t.Fatal("sheetHTML is not deterministic")
	}
	if !strings.Contains(out, v.render()) {
		t.Error("rendered stylesheet not inlined")
	}
	states := 0
	for _, r := range v.Roles {
		if r.Kind == "state" {
			states++
			if !strings.Contains(out, `data-state="`+r.Class+`"`) {
				t.Errorf("state %s has no checkbox", r.Class)
			}
			if !strings.Contains(out, `<body class="`+r.Class+`"`) {
				t.Errorf("state %s is not on by default", r.Class)
			}
			continue
		}
		if n := strings.Count(out, `class="`+r.Class+`"`); n != 1 {
			t.Errorf("role %s: %d samples, want 1", r.Class, n)
		}
	}
	if n := strings.Count(out, `type="checkbox"`); n != states {
		t.Errorf("%d checkboxes, want %d", n, states)
	}
	if !strings.Contains(out, "src/styles/fonts.css") || !strings.Contains(out, `tgs-missing">missing`) {
		t.Error("missing allow.css not listed as missing")
	}
}

func TestSheetEscapesUseAndInlinesFonts(t *testing.T) {
	dir := t.TempDir()
	vocab := `[project]
out = "out.css"

[tokens.color]
ink = "#111111"

[role.text.plain]
use = "Use <b>this</b> & that."
color = "{color.ink}"

[allow.css]
"fonts.css" = "self-hosted"
`
	if err := os.WriteFile(filepath.Join(dir, "tidgrid.toml"), []byte(vocab), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "fonts.css"), []byte("/*FONT-MARK*/</STYLE><script>x</script>"), 0o644); err != nil {
		t.Fatal(err)
	}
	v, errs := load(filepath.Join(dir, "tidgrid.toml"))
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	out := sheetHTML(v)
	if strings.Contains(out, "<b>this</b>") {
		t.Error("use text not escaped")
	}
	if !strings.Contains(out, "Use &lt;b&gt;this&lt;/b&gt; &amp; that.") {
		t.Error("escaped use text missing")
	}
	if !strings.Contains(out, "/*FONT-MARK*/") {
		t.Error("existing allow.css not inlined")
	}
	if strings.Contains(out, "</STYLE><script>") {
		t.Error("inlined CSS can close the style element")
	}
	if strings.Contains(out, `tgs-missing">missing`) {
		t.Error("existing allow.css listed as missing")
	}
}
