package vocab

import (
	"testing"
)

// setSrc is a project with one mixin set on line 2 and two roles that include it.
func setSrc(set string) string {
	return projectLine + `m = mixin("m", use = "u", sets = {"s": ` + set + `})` + "\n" +
		`role("text", "a", use = "A.", include = [m.s])` + "\n" +
		`role("text", "b", use = "B.", include = [m.s])` + "\n"
}

func TestMixinSetNeedsTwoDeclarations(t *testing.T) {
	bad := []string{`{"margin": "0"}`, `{"&:hover": {"opacity": "1"}}`}
	for _, set := range bad {
		got := loadErr(t, map[string]string{"index.tg": setSrc(set)})
		wantIn(t, got, "index.tg:2: mixin.m.s: a set needs at least two lines; a single shared value belongs in a token")
	}
}

func TestMixinSetWithTwoDeclarationsBuilds(t *testing.T) {
	good := []string{
		`{"margin": "0", "&:hover": {"opacity": "1"}}`,
		`{"&:focus-visible": {"outline": "none", "outline-offset": "0"}}`,
		`{"margin": "0", "padding": "1"}`,
	}
	for _, set := range good {
		mixinLoad(t, map[string]string{"index.tg": setSrc(set)})
	}
}

func TestMixinOneTopLineSetHandWrittenBuilds(t *testing.T) {
	// one top-level line plus a nested rule: the one line by hand is not a written-out set
	mixinLoad(t, map[string]string{"index.tg": projectLine + `m = mixin("m", use = "u", sets = {"s": {"font-weight": "600", "&:hover": {"opacity": "1"}}})
role("text", "a", use = "A.", include = [m.s])
role("text", "b", use = "B.", include = [m.s])
role("text", "c", use = "C.", css = {"font-weight": "600"})
`})
}
