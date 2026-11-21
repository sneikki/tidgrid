package main

import (
	"strings"
	"testing"
)

func loadExampleVocab(t *testing.T) *Vocabulary {
	t.Helper()
	v, errs := load("testdata/example.toml")
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	return v
}

func TestListExample(t *testing.T) {
	v := loadExampleVocab(t)
	out := listText(v)
	if again := listText(v); again != out {
		t.Fatal("listText is not deterministic")
	}
	for _, r := range v.Roles {
		if !strings.Contains(out, r.Class+"  "+oneLine(r.Use)) {
			t.Errorf("role %s: class or use text missing", r.Class)
		}
	}
	if !strings.Contains(out, "## set by scripts") {
		t.Error("state kind has no 'set by scripts' heading")
	}
	for _, a := range v.Atoms {
		base, _, _ := strings.Cut(a.Class, "(")
		if !strings.Contains(out, base) {
			t.Errorf("atom family %s missing", base)
		}
	}
	for _, a := range v.AllowClass {
		if !strings.Contains(out, "class "+a.Name+"  "+a.Reason) {
			t.Errorf("allow class %s missing", a.Name)
		}
	}
	for _, a := range v.AllowCSS {
		if !strings.Contains(out, "css "+a.Name+"  "+a.Reason) {
			t.Errorf("allow css %s missing", a.Name)
		}
	}
	for _, g := range v.Tokens {
		for _, tk := range g.Tokens {
			if !strings.Contains(out, tk.Name) {
				t.Errorf("token %s.%s name missing", g.Name, tk.Name)
			}
			if strings.Contains(out, tk.Value) {
				t.Errorf("token %s.%s value %q leaked", g.Name, tk.Name, tk.Value)
			}
		}
	}
	if strings.Contains(out, "var(") {
		t.Error("a resolved var() reached the list")
	}
	for _, s := range []string{"S = a name in space, Z = a name in size, B = a breakpoint name", "tg-gap(S)\n", "tg-width(Z)\n",
		"tg-cols(1-6|12)\n", "tg-span(1-12)\n", "tg-order(first|last)\n", "tg-hide\n", "tg-hide-below(B)\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("list lacks %q", s)
		}
	}
	if strings.Contains(out, "(space)") {
		t.Error("scale repeated on atom lines")
	}
	if !strings.Contains(out, listRules) {
		t.Error("agent rules missing")
	}
}
