package main

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden files")

func loadStr(t *testing.T, s string) (*Vocabulary, []error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "tidgrid.toml")
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	return load(p)
}

func TestGolden(t *testing.T) {
	v, errs := load("testdata/example.toml")
	if errs != nil {
		t.Fatal(errs)
	}
	got := v.render()
	if *update {
		os.WriteFile("testdata/example.css", []byte(got), 0o644)
	}
	want, _ := os.ReadFile("testdata/example.css")
	if got != string(want) {
		t.Fatal("output differs from testdata/example.css (go test -update)")
	}
	if v.render() != got {
		t.Fatal("render not deterministic")
	}
	for _, s := range []string{`.tg-state\(ready\) & {`, "@keyframes tg-rise", "  &:hover {", "@media (min-width: 64em) {\n    border-width"} {
		if !strings.Contains(got, s) {
			t.Errorf("missing %q", s)
		}
	}
}

const base = `
[project]
out = "o.css"
[breakpoints]
md = "64em"
[tokens.color]
ink = "#000"
[tokens.space]
sm = "1rem"
`

func TestValidationErrors(t *testing.T) {
	cases := []struct{ name, toml, want string }{
		{"no out", "[project]\nprefix = \"tg\"\n", "project.out"},
		{"no use", base + "[role.surface.card]\ncolor = \"{color.ink}\"\n", "role.surface.card.use"},
		{"unknown ref", base + "[role.surface.card]\nuse = \"x\"\ncolor = \"{color.nope}\"\n", "role.surface.card.color: unknown reference {color.nope}"},
		{"raw hex", base + "[role.surface.card]\nuse = \"x\"\ncolor = \"#fff\"\n", "role.surface.card.color: raw colour"},
		{"raw fn", base + "[role.surface.card]\nuse = \"x\"\ncolor = \"rgb(1 2 3)\"\n", "raw colour"},
		{"raw length", base + "[role.surface.card]\nuse = \"x\"\npadding = \"calc({space.sm} + 2px)\"\n", "role.surface.card.padding: raw length"},
		{"raw time", base + "[role.surface.card]\nuse = \"x\"\ntransition = \"all 200ms\"\n", "raw length or time 200ms"},
		{"kind is atom", base + "[role.gap.big]\nuse = \"x\"\n", "role.gap: kind collides"},
		{"no &", base + "[role.surface.card]\nuse = \"x\"\n[role.surface.card.\":hover\"]\ncolor = \"red\"\n", "role.surface.card.\":hover\": selector must contain &"},
		{"unknown class", base + "[role.surface.card]\nuse = \"x\"\n[role.surface.card.\".tg-gap(zz) &\"]\ncolor = \"red\"\n", "unknown class tg-gap(zz); tg-gap takes: sm"},
		{"empty reason", base + "[allow.class]\n\"x\" = \"\"\n", "allow.class.x: reason required"},
		{"css outside", base + "[allow.css]\n\"../../secret.css\" = \"x\"\n", "path must stay inside"},
		{"bad bp", base + "[role.surface.card]\nuse = \"x\"\n[role.surface.card.\"@lg\"]\ncolor = \"red\"\n", "unknown breakpoint"},
		{"colour name", base + "[role.surface.card]\nuse = \"x\"\ncolor = \"red\"\n", "role.surface.card.color: raw colour red"},
		{"colour name in shorthand", base + "[role.surface.card]\nuse = \"x\"\nborder = \"{space.sm} solid Tan\"\n", "raw colour Tan"},
		{"unit dvmin", base + "[role.surface.card]\nuse = \"x\"\nwidth = \"10dvmin\"\n", "raw length or time 10dvmin"},
		{"unit exponent", base + "[role.surface.card]\nuse = \"x\"\nwidth = \"1.5e3px\"\n", "raw length or time 1.5e3px"},
		{"at-rule raw", base + "[role.surface.card]\nuse = \"x\"\n[role.surface.card.\"@media (min-width: 800px)\"]\ncolor = \"{color.ink}\"\n", "role.surface.card.\"@media (min-width: 800px)\": raw length or time 800px"},
		{"at-rule unknown ref", base + "[role.surface.card]\nuse = \"x\"\n[role.surface.card.\"@media (min-width: {size.nope})\"]\ncolor = \"{color.ink}\"\n", "unknown reference {size.nope}"},
		{"token name slash", "[project]\nout = \"o.css\"\n[tokens.size]\n\"1/2\" = \"50%\"\n", `tokens.size."1/2": token name`},
		{"token group slash", "[project]\nout = \"o.css\"\n[tokens.\"a b\"]\nx = \"1\"\n", `tokens."a b": group name`},
		{"use in sub-table", base + "[role.surface.card]\nuse = \"x\"\n[role.surface.card.\"&:hover\"]\nuse = \"y\"\n", `role.surface.card."&:hover".use: use belongs`},
		{"scan string", "[project]\nout = \"o.css\"\nscan = \"src\"\n", "project.scan: must be a list"},
		{"helpers string", "[project]\nout = \"o.css\"\nhelpers = \"cn\"\n", "project.helpers: must be a list"},
		{"token loop", base + "[tokens.t]\na = \"{t.b}\"\nb = \"{t.a}\"\n", "reference loop"},
	}
	for _, c := range cases {
		_, errs := loadStr(t, c.toml)
		var all []string
		for _, e := range errs {
			all = append(all, e.Error())
		}
		if !strings.Contains(strings.Join(all, "\n"), c.want) {
			t.Errorf("%s: want %q in %v", c.name, c.want, all)
		}
	}
}

func TestNotRaw(t *testing.T) {
	ok := base + `[tokens.hue]
red = "#f00"
[tokens.size]
"2xl" = "9rem"
[role.surface.card]
use = "x"
margin = "0 {space.sm}"
width = "50%"
opacity = "0.5"
flex = "1 1 0"
transform = "translate3d(0, 0, 0) rotate(45deg)"
height = "calc({size.2xl} * 2)"
background = "color-mix(in srgb, {color.ink} 18%, transparent)"
content = "\"10px red\""
background-image = "url(img-10px.png), url(#face)"
font-family = "\"Tan Pearl\", serif"
grid-area = "gold"
transition = "tg-rise predicate {color.ink} transparent currentColor"
border-color = "{hue.red}"
rotate = "0.5turn"
[role.surface.card."@media (max-height: {size.2xl})"]
color = "{color.ink}"
[role.surface.card."@below-md"]
color = "{color.ink}"
[role.surface.card."&:hover"."@md"."& h1"]
color = "{color.ink}"
[role.surface.card."@media (prefers-reduced-motion: reduce)"]
animation = "none"
`
	if _, errs := loadStr(t, ok); errs != nil {
		t.Fatal(errs)
	}
}

// A token in an at-rule key becomes its written value, nested references resolved, never var().
func TestAtRuleKeyLiteral(t *testing.T) {
	v, errs := loadStr(t, base+"[tokens.size]\nshort = \"{space.sm}\"\n[role.surface.card]\nuse = \"x\"\n[role.surface.card.\"@media (max-height: {size.short})\"]\ncolor = \"{color.ink}\"\n")
	if errs != nil {
		t.Fatal(errs)
	}
	css := v.render()
	if !strings.Contains(css, "@media (max-height: 1rem) {") {
		t.Errorf("at-rule key not literal:\n%s", css)
	}
}

func TestEscapeAndOrder(t *testing.T) {
	v, errs := loadStr(t, base+"[role.surface.card]\nuse = \"x\"\n[role.surface.card.\".md:tg-gap(sm) &\"]\ncolor = \"{color.ink}\"\n")
	if errs != nil {
		t.Fatal(errs)
	}
	if !strings.Contains(v.render(), `.md\:tg-gap\(sm\) & {`) {
		t.Fatal("selector class not escaped")
	}
	if !v.Valid("md:tg-gap(sm)") || v.Valid("tg-gap(md)") {
		t.Fatal("class set wrong")
	}
	if tk, _ := v.Takes("md:tg-gap"); strings.Join(tk, " ") != "sm" {
		t.Fatal(tk)
	}
}

// README atom table, row by row.
func TestAtomTable(t *testing.T) {
	v, errs := loadStr(t, base+"[tokens.size]\nfull = \"100%\"\n")
	if errs != nil {
		t.Fatal(errs)
	}
	got := map[string]string{}
	for _, a := range v.Atoms {
		var p []string
		for _, kv := range a.Decls {
			p = append(p, kv[0]+": "+kv[1])
		}
		got[a.Class] = strings.Join(p, "; ")
	}
	sp := "var(--tg-space-sm)"
	want := map[string]string{
		"tg-row": "display: flex; flex-direction: row", "tg-col": "display: flex; flex-direction: column",
		"tg-wrap": "flex-wrap: wrap", "tg-grow": "flex: 1 1 0; min-width: 0", "tg-fixed": "flex: none",
		"tg-center": "margin-inline: auto", "tg-hide": "display: none", "tg-show": "display: revert",
		"tg-order(first)": "order: -1", "tg-order(last)": "order: 1",
		"tg-gap(sm)": "gap: " + sp, "tg-x-gap(sm)": "column-gap: " + sp, "tg-y-gap(sm)": "row-gap: " + sp,
		"tg-pad(sm)": "padding: " + sp, "tg-pad-x(sm)": "padding-inline: " + sp, "tg-pad-y(sm)": "padding-block: " + sp,
		"tg-pad-t(sm)": "padding-top: " + sp, "tg-pad-r(sm)": "padding-right: " + sp,
		"tg-pad-b(sm)": "padding-bottom: " + sp, "tg-pad-l(sm)": "padding-left: " + sp,
		"tg-width(full)": "width: var(--tg-size-full)", "tg-min-width(full)": "min-width: var(--tg-size-full)",
		"tg-max-width(full)": "max-width: var(--tg-size-full)", "tg-height(full)": "height: var(--tg-size-full)",
		"tg-min-height(full)": "min-height: var(--tg-size-full)",
	}
	for k, w := range map[string]string{"start": "flex-start", "center": "center", "end": "flex-end", "baseline": "baseline", "stretch": "stretch"} {
		want["tg-align("+k+")"] = "align-items: " + w
	}
	for k, w := range map[string]string{"start": "flex-start", "center": "center", "end": "flex-end", "between": "space-between", "around": "space-around", "evenly": "space-evenly"} {
		want["tg-justify("+k+")"] = "justify-content: " + w
	}
	for k, w := range map[string]string{"start": "flex-start", "center": "center", "end": "flex-end", "stretch": "stretch"} {
		want["tg-self("+k+")"] = "align-self: " + w
	}
	for _, k := range []string{"start", "center", "end"} {
		want["tg-text-align("+k+")"] = "text-align: " + k
	}
	for _, n := range []int{1, 2, 3, 4, 5, 6, 12} {
		s := strconv.Itoa(n)
		want["tg-cols("+s+")"] = "display: grid; grid-template-columns: repeat(" + s + ", minmax(0, 1fr))"
	}
	for n := 1; n <= 12; n++ {
		s := strconv.Itoa(n)
		want["tg-span("+s+")"] = "grid-column: span " + s
	}
	for c, w := range want {
		if got[c] != w {
			t.Errorf("%s: got %q want %q", c, got[c], w)
		}
		if !v.Valid("md:" + c) {
			t.Errorf("no breakpoint variant of %s", c)
		}
	}
	if len(got) != len(want) {
		t.Errorf("atoms: %d, table rows: %d", len(got), len(want))
	}
	n := len(v.Atoms)
	if v.Atoms[n-2].Class != "tg-show" || v.Atoms[n-1].Class != "tg-hide" {
		t.Error("tg-show then tg-hide must be last")
	}
	if !strings.Contains(v.render(), ".md\\:tg-show {\n    display: revert;\n  }\n  .md\\:tg-hide {\n    display: none;\n  }\n}") {
		t.Error("breakpoint block must end with show, hide")
	}
}

func TestBuildCmd(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "t.toml")
	os.WriteFile(cfg, []byte(base), 0o644)
	if run([]string{"build", "-c", cfg}) != 0 {
		t.Fatal("build failed")
	}
	a, _ := os.ReadFile(filepath.Join(dir, "o.css"))
	run([]string{"build", "-c", cfg})
	b, _ := os.ReadFile(filepath.Join(dir, "o.css"))
	if len(a) == 0 || string(a) != string(b) {
		t.Fatal("not deterministic")
	}
	if run([]string{"build", "-c", filepath.Join(dir, "missing.toml")}) != 2 {
		t.Fatal("missing file must exit 2")
	}
	if run([]string{"check"}) != 2 {
		t.Fatal("stub must exit 2")
	}
}
