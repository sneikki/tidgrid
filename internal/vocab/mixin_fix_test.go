package vocab

import (
	"path/filepath"
	"strings"
	"testing"
)

func mixinLoad(t *testing.T, files map[string]string) *Vocabulary {
	t.Helper()
	dir := writeTree(t, files)
	v, errs := Load(filepath.Join(dir, "index.tg"))
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	return v
}

func TestMixinSetFixedAtCall(t *testing.T) {
	v := mixinLoad(t, map[string]string{"index.tg": projectLine + `d = {"font-weight": "600", "&:hover": {"font-weight": "500"}}
m = mixin("m", use = "u", sets = {"s": d})
d["font-weight"] = "700"
d["margin"] = "0"
d["&:hover"]["font-weight"] = "900"
role("text", "a", use = "A.", include = [m.s])
role("text", "b", use = "B.", include = [m.s])
role("text", "c", use = "C.", css = {"font-weight": "700", "margin": "0"})
`})
	css := v.Render()
	wantIn(t, css, ".tg-text\\(a\\) {\n  font-weight: 600;\n  &:hover {\n    font-weight: 500;\n  }\n}\n")
	wantIn(t, css, ".tg-text\\(b\\) {\n  font-weight: 600;\n  &:hover {\n    font-weight: 500;\n  }\n}\n")
	if got := strings.Join(v.Mixins[0].Sets[0].Lines, "|"); got != "font-weight: 600|&:hover { font-weight: 500 }" {
		t.Errorf("list lines %q", got)
	}
}

func TestMixinHandWrittenUsesSetAsCalled(t *testing.T) {
	got := loadErr(t, map[string]string{"index.tg": projectLine + `d = {"font-weight": "600"}
m = mixin("m", use = "u", sets = {"s": d})
d["font-weight"] = "700"
role("text", "a", use = "A.", include = [m.s])
role("text", "b", use = "B.", include = [m.s])
role("text", "c", use = "C.", css = {"font-weight": "600"})
role("text", "d", use = "D.", css = {"font-weight": "700"})
`})
	wantIn(t, got, "role.text.c: own lines repeat set m.s")
	if strings.Contains(got, "role.text.d") {
		t.Errorf("check followed the mutated set: %s", got)
	}
}

func TestMixinSetErrorReportedOnce(t *testing.T) {
	got := loadErr(t, map[string]string{"index.tg": projectLine + `m = mixin("m", use = "u", sets = {
    "s": {"color": "#fff"},
})
role("text", "a", use = "A.", include = [m.s])
role("text", "b", use = "B.", include = [m.s])
role("text", "c", use = "C.", include = [m.s])
`})
	if n := strings.Count(got, "raw colour"); n != 1 {
		t.Fatalf("want 1 error, got %d:\n%s", n, got)
	}
	wantIn(t, got, "index.tg:2: mixin.m.s.color: raw colour #fff")
}

const nestedHead = projectLine + `m = mixin("m", use = "u", sets = {
    "s": {"font-weight": "600", "&:hover": {"font-weight": "500"}},
    "t": {"&:hover": {"font-weight": "400"}},
})
`

func TestMixinNestedRulesMerge(t *testing.T) {
	v := mixinLoad(t, map[string]string{"index.tg": nestedHead + `role("text", "a", use = "A.", include = [m.s], css = {"&:hover": {"margin": "0"}})
role("text", "b", use = "B.", include = [m.s])
role("text", "c", use = "C.", include = [m.t])
role("text", "d", use = "D.", include = [m.t])
`})
	want := ".tg-text\\(a\\) {\n  font-weight: 600;\n  &:hover {\n    font-weight: 500;\n    margin: 0;\n  }\n}\n"
	wantIn(t, v.Render(), want)
	hover := v.Roles[0].Rule.Nodes[1]
	if hover.From != "m.s" || hover.Sub.Nodes[0].From != "m.s" || hover.Sub.Nodes[1].From != "" {
		t.Fatalf("origins after merge: %+v", hover)
	}
}

func TestMixinNestedClash(t *testing.T) {
	got := loadErr(t, map[string]string{"index.tg": nestedHead + `role("text", "a", use = "A.", include = [m.s], css = {"&:hover": {"font-weight": "800"}})
role("text", "b", use = "B.", include = [m.s, m.t])
role("text", "c", use = "C.", include = [m.t])
role("text", "d", use = "D.", include = [m.s])
`})
	wantIn(t, got, `role.text.a."&:hover".font-weight: given twice: by set m.s and by the role's own css`)
	wantIn(t, got, `role.text.b."&:hover".font-weight: given twice: by set m.s and by set m.t`)
}

func TestMixinIncludedTwice(t *testing.T) {
	got := mixinErr(t, `role("text", "a", use = "A.", include = [type.label, type.label])`)
	wantIn(t, got, "role.text.a.include: set type.label is included twice")
	if strings.Contains(got, "given twice") {
		t.Errorf("clash message for a repeated set: %s", got)
	}
}

func TestMixinHandWrittenNeedsSameValue(t *testing.T) {
	mixinLoad(t, map[string]string{"index.tg": mixinHead + twoUsers + `role("text", "c", use = "C.", css = {"font-weight": "700"})` + "\n"})
}

func TestMixinHandWrittenIgnoresEmptyAndNestedOnlySets(t *testing.T) {
	got := loadErr(t, map[string]string{"index.tg": projectLine + `m = mixin("m", use = "u", sets = {
    "empty": {},
    "nest": {"&:hover": {"color": "red"}},
})
role("text", "a", use = "A.", css = {"margin": "0", "&:hover": {"color": "red"}})
`})
	if strings.Contains(got, "role.text.a:") {
		t.Errorf("hand-written rule fired: %s", got)
	}
}

func TestMixinUsersCountedAcrossFiles(t *testing.T) {
	v := mixinLoad(t, map[string]string{
		"index.tg": projectLine + "load(\"sets.tg\", \"m\")\nload(\"b.tg\", \"mk\")\nrole(\"text\", \"a\", use = \"A.\", include = [m.s])\nmk(m)\n",
		"sets.tg":  `m = mixin("m", use = "u", sets = {"s": {"font-weight": "600"}})` + "\n",
		"b.tg":     "def mk(m):\n    role(\"text\", \"b\", use = \"B.\", include = [m.s])\n",
	})
	if got := v.Mixins[0].Sets[0].Roles; len(got) != 2 {
		t.Fatalf("roles %v", got)
	}
}
