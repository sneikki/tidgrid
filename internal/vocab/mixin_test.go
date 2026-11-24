package vocab

import (
	"path/filepath"
	"strings"
	"testing"
)

const mixinHead = projectLine + `c = tokens("color", ink = "#111", sand = "#eee")
type = mixin("type", use = "A text style.", sets = {
    "heading": {"font-size": "700", "color": c.ink, "&:hover": {"color": c.sand}},
    "label": {"font-weight": "600"},
})
`

func mixinErr(t *testing.T, body string) string {
	t.Helper()
	return loadErr(t, map[string]string{"index.tg": mixinHead + body})
}

func wantIn(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

const twoUsers = `role("text", "a", use = "A.", include = [type.heading, type.label], css = {"margin": "0"})
role("text", "b", use = "B.", include = [type.heading, type.label])
`

func TestMixinRenderOrder(t *testing.T) {
	dir := writeTree(t, map[string]string{"index.tg": mixinHead + twoUsers})
	v, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	css := v.Render()
	wantIn(t, css, ".tg-text\\(a\\) {\n  font-size: 700;\n  color: var(--tg-color-ink);\n  &:hover {\n    color: var(--tg-color-sand);\n  }\n  font-weight: 600;\n  margin: 0;\n}\n")
	if strings.Contains(css, "heading") || strings.Contains(css, "mixin") {
		t.Error("a mixin reached the stylesheet")
	}
	r := v.Roles[0]
	if len(r.Includes) != 2 || r.Includes[0] != "type.heading" {
		t.Fatalf("includes %v", r.Includes)
	}
	own := r.Rule.Nodes[len(r.Rule.Nodes)-1]
	if own.Prop != "margin" || own.From != "" || r.Rule.Nodes[0].From != "type.heading" {
		t.Fatalf("origins not kept: %+v", r.Rule.Nodes)
	}
}

func TestMixinOwnClashesWithSet(t *testing.T) {
	got := mixinErr(t, `role("text", "a", use = "A.", include = [type.heading], css = {"font-size": "1"})`+"\n"+twoUsers[:0])
	wantIn(t, got, `index.tg:7: role.text.a.font-size: given twice: by set type.heading and by the role's own css`)
}

func TestMixinTwoSetsClash(t *testing.T) {
	got := mixinErr(t, `x = mixin("x", use = "u", sets = {"s": {"color": "red"}})`+"\n"+
		`role("text", "a", use = "A.", include = [type.heading, x.s])`)
	wantIn(t, got, "role.text.a.color: given twice: by set type.heading and by set x.s")
}

func TestMixinNeedsTwoUsers(t *testing.T) {
	got := mixinErr(t, `role("text", "a", use = "A.", include = [type.heading])`)
	wantIn(t, got, "mixin.type.heading: included only by role text.a; a set needs at least two")
	wantIn(t, got, "mixin.type.label: included by no role; a set needs at least two")
}

func TestMixinHandWrittenSet(t *testing.T) {
	got := mixinErr(t, twoUsers+`role("text", "c", use = "C.", css = {"font-weight": "600", "color": "red"})`)
	wantIn(t, got, "role.text.c: own lines repeat set type.label; write include = [type.label]")
	if strings.Contains(got, "role.text.a") || strings.Contains(got, "role.text.b") {
		t.Errorf("fired on a role that includes the set: %s", got)
	}
}

func TestMixinIncludeNeedsSet(t *testing.T) {
	got := mixinErr(t, twoUsers+`role("text", "c", use = "C.", include = ["heading", c.ink])`)
	wantIn(t, got, `index.tg:9: role.text.c.include: "heading" is not a mixin set`)
}

func TestMixinValueIsNotCSS(t *testing.T) {
	got := mixinErr(t, twoUsers+`role("text", "c", use = "C.", css = {"color": type, "x": type.heading})`)
	wantIn(t, got, "role.text.c.color: value must be a string")
}

func TestMixinNeedsUse(t *testing.T) {
	got := loadErr(t, map[string]string{"index.tg": projectLine + `mixin("m", sets = {"s": {"color": "red"}})`})
	wantIn(t, got, "index.tg:2: mixin.m.use: required")
}

func TestMixinDeclaredTwice(t *testing.T) {
	got := mixinErr(t, `mixin("type", use = "again", sets = {})`)
	wantIn(t, got, "mixin.type: declared twice (index.tg:3 and index.tg:7)")
}
