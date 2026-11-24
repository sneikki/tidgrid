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

func loadTidy(t *testing.T, call string) *Vocabulary {
	t.Helper()
	dir := writeTree(t, map[string]string{"index.tg": tidyBase + call})
	v, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	return v
}

func renderTidy(t *testing.T, call string) string {
	t.Helper()
	return loadTidy(t, call).Render()
}

func TestTidinessUnsetOutputIsTheSameWithTheFeatureRemoved(t *testing.T) {
	with := loadTidy(t, `tidiness("full")`)
	with.Tidiness, with.tidyRules = "", nil
	if got, want := with.Render(), renderTidy(t, ""); got != want {
		t.Fatalf("output differs:\n%s\nwant\n%s", got, want)
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
	if strings.Count(out, "/*") > 1 {
		t.Fatal("comments written")
	}
	if strings.Contains(out, "\n\n\n") {
		t.Fatal("blank run written")
	}
	if !strings.Contains(out, "\n    color: color-mix(in oklab, currentcolor 50%, transparent);\n") {
		t.Fatal("nested @supports rule not indented like other rules")
	}
}

// Every declaration of each embedded file reaches the output, so an edit that
// the parser would swallow fails here.
func TestTidinessEveryDeclarationOfTheFilesIsRendered(t *testing.T) {
	dir := writeTree(t, map[string]string{"index.tg": "project(out = \"o.css\")\ntidiness(\"full\")\n"})
	v, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	out := v.Render()
	out = out[:strings.Index(out, "\n.tg-")] // the reset only, not the atoms after it
	total := 0
	for _, lvl := range TidinessLevels {
		src, _ := tidinessFS.ReadFile("tidiness/" + lvl + ".css")
		for _, d := range strings.Split(cssComment.ReplaceAllString(string(src), ""), ";") {
			d = d[strings.LastIndexAny(d, "{}")+1:] // the text after the block opener or closer
			prop, value, ok := strings.Cut(d, ":")
			if !ok {
				continue
			}
			total++
			want := strings.TrimSpace(prop) + ": " + strings.Join(strings.Fields(value), " ") + ";"
			if !strings.Contains(out, want) {
				t.Errorf("%s.css: %q not rendered", lvl, want)
			}
		}
	}
	if got := strings.Count(out, ";\n"); got != total {
		t.Errorf("rendered %d declarations, files have %d", got, total)
	}
}

func TestTidinessParserRefusesWhatItCannotRead(t *testing.T) {
	for name, src := range map[string]string{
		"stray close brace":     "a { color: red; }\n}\nb { color: blue; }",
		"unclosed block":        "a { color: red;",
		"no colon":              "a { color red; }",
		"empty value":           "a { color: ; }",
		"missing semicolon":     "a { color: red\n  margin: 0; }",
		"top-level statement":   "@import url(x.css);\na { color: red; }",
		"unclosed quote":        `a { content: "x; }`,
		"unclosed bracket":      "a[title='x' { color: red; }",
		"stray close bracket":   "a] { color: red; }",
		"unterminated comment":  "/* x\na { color: red; }",
		"text after last block": "a { color: red; } b",
	} {
		if _, err := parseRules(src); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

func TestTidinessParserKeepsQuotedAndBracketedDelimiters(t *testing.T) {
	rules, err := parseRules(`a[title=';'] { content: "}{;"; }`)
	if err != nil || len(rules) != 1 || rules[0].Key != `a[title=';']` || rules[0].Nodes[0].Value != `"}{;"` {
		t.Fatalf("got %v, %v", rules, err)
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
		if name != "twice" && !(strings.Contains(msg, `"minimal"`) && strings.Contains(msg, `"full"`)) {
			t.Errorf("%s: levels not named in %q", name, msg)
		}
	}
}
