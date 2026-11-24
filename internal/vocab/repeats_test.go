package vocab

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// rr writes a role of kind text named name with the given own lines, "prop=value" each.
func rr(name string, lines ...string) string {
	var css []string
	for _, ln := range lines {
		k, v, _ := strings.Cut(ln, "=")
		css = append(css, fmt.Sprintf("%q: %q", k, v))
	}
	return fmt.Sprintf("role(\"text\", %q, use = \"U.\", css = {%s})\n", name, strings.Join(css, ", "))
}

func repeatsLoad(t *testing.T, proj string, body ...string) (string, bool) {
	t.Helper()
	src := proj + strings.Join(body, "")
	_, errs := Load(filepath.Join(writeTree(t, map[string]string{"index.tg": src}), "index.tg"))
	var msgs []string
	for _, e := range errs {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "\n"), len(errs) > 0
}

func repeatsFail(t *testing.T, body ...string) string {
	t.Helper()
	got, failed := repeatsLoad(t, projectLine, body...)
	if !failed {
		t.Fatal("want an error")
	}
	return got
}

func repeatsPass(t *testing.T, proj string, body ...string) {
	t.Helper()
	if got, failed := repeatsLoad(t, proj, body...); failed {
		t.Fatalf("want no error, got %s", got)
	}
}

const l1, l2, l3 = "margin=0", "padding=1", "display=flex"

func TestRepeatsThreeLinesTwoRolesFails(t *testing.T) {
	got := repeatsFail(t, rr("a", l1, l2, l3), rr("b", l1, l2, l3))
	wantIn(t, got, "index.tg:2: role.text.a: own lines also in role text.b: margin: 0; padding: 1; display: flex; make them a mixin set and include it")
}

func TestRepeatsTwoLinesInTwoOrThreeRolesPasses(t *testing.T) {
	repeatsPass(t, projectLine, rr("a", l1, l2), rr("b", l1, l2), rr("c", l1, l2))
}

func TestRepeatsTwoLinesFourRolesFails(t *testing.T) {
	got := repeatsFail(t, rr("a", l1, l2), rr("b", l1, l2), rr("c", l1, l2), rr("d", l1, l2))
	wantIn(t, got, "role.text.a: own lines also in role text.b, text.c, text.d: margin: 0; padding: 1")
}

func TestRepeatsDifferentValueDoesNotCount(t *testing.T) {
	repeatsPass(t, projectLine, rr("a", l1, l2, "display=flex"), rr("b", l1, l2, "display=grid"))
}

func TestRepeatsSetLinesDoNotCount(t *testing.T) {
	set := `s = mixin("s", use = "u", sets = {"x": {"margin": "0", "padding": "1", "display": "flex"}})` + "\n"
	inc := func(n string) string {
		return fmt.Sprintf("role(\"text\", %q, use = \"U.\", include = [s.x])\n", n)
	}
	repeatsPass(t, projectLine, set, inc("a"), inc("b"))
}

func TestRepeatsLargerGroupSwallowsSubGroups(t *testing.T) {
	// a,b,c share three lines; the pair of lines inside them has the same roles, so one report
	got := repeatsFail(t, rr("a", l1, l2, l3), rr("b", l1, l2, l3), rr("c", l1, l2, l3))
	if n := strings.Count(got, "own lines also in"); n != 1 {
		t.Fatalf("want one report, got %d: %s", n, got)
	}
	wantIn(t, got, "text.b, text.c: margin: 0; padding: 1; display: flex")
}

func TestRepeatsGroupNoPairYieldsAlone(t *testing.T) {
	// a and d share only {m,p}; the pair yields it with all four roles holding it, so (2,4) fires
	got := repeatsFail(t,
		rr("a", l1, l2, "x=1", "y=1"), rr("b", l1, l2, "x=1", "z=1"),
		rr("c", l1, l2, "y=1", "z=1"), rr("d", l1, l2, "w=1"))
	wantIn(t, got, "role.text.a: own lines also in role text.b, text.c, text.d: margin: 0; padding: 1")
}

func TestRepeatsOutputStable(t *testing.T) {
	body := []string{rr("a", l1, l2, l3), rr("b", l1, l2, l3), rr("c", "x=1", "y=1", "z=1"), rr("d", "x=1", "y=1", "z=1")}
	first := repeatsFail(t, body...)
	for i := 0; i < 20; i++ {
		if got := repeatsFail(t, body...); got != first {
			t.Fatalf("output changed:\n%s\n--\n%s", first, got)
		}
	}
	if strings.Index(first, "role.text.a") > strings.Index(first, "role.text.c") {
		t.Fatalf("not in file order: %s", first)
	}
}

func TestRepeatsEmptyListTurnsOff(t *testing.T) {
	repeatsPass(t, `project(out = "o.css", repeats = [])`+"\n", rr("a", l1, l2, l3), rr("b", l1, l2, l3))
}

func TestRepeatsSettingChangesPairs(t *testing.T) {
	proj := `project(out = "o.css", repeats = [(2, 2)])` + "\n"
	got, failed := repeatsLoad(t, proj, rr("a", l1, l2), rr("b", l1, l2))
	if !failed {
		t.Fatal("want (2, 2) to fire")
	}
	wantIn(t, got, "role.text.a: own lines also in role text.b")
}

func TestRepeatsMalformedSettingErrors(t *testing.T) {
	for _, bad := range []string{`"x"`, `[3]`, `[(3, 0)]`, `[(3, 1)]`, `[(0, 2)]`, `[(3, 2, 1)]`, `[(3, "a")]`, `[(-1, 2)]`, `3`, `[[3, 2]]`} {
		got, failed := repeatsLoad(t, `project(out = "o.css", repeats = `+bad+")\n")
		if !failed {
			t.Errorf("repeats = %s: want an error", bad)
		}
		wantIn(t, got, "project.repeats: must be a list of (lines, roles) pairs")
	}
}

// assertRepeatsBounded loads body and checks it finishes fast and prints at most 21 lines (20 groups and the closing line).
func assertRepeatsBounded(t *testing.T, body []string) {
	t.Helper()
	start := time.Now()
	got, failed := repeatsLoad(t, projectLine, body...)
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("took %s, want well under a second", elapsed)
	}
	if !failed {
		t.Fatal("want repeat groups reported")
	}
	if n := strings.Count(got, "\n") + 1; n > 21 {
		t.Fatalf("printed %d lines, want at most 21", n)
	}
}

func TestRepeatsDenseCaseIsBounded(t *testing.T) {
	assertRepeatsBounded(t, denseBody())
}

func TestRepeatsHeavyOverlapIsBounded(t *testing.T) {
	// 100 roles, 15 lines each, drawn from a pool of 20 so neighbours share most lines
	var body []string
	for i := 0; i < 100; i++ {
		var ls []string
		for j := 0; j < 15; j++ {
			ls = append(ls, fmt.Sprintf("p%d=v", (i+j)%20))
		}
		body = append(body, rr(fmt.Sprintf("r%d", i), ls...))
	}
	assertRepeatsBounded(t, body)
}

func TestRepeatsReportsAtMostTwentyWithClosingLine(t *testing.T) {
	got := repeatsFail(t, denseBody()...)
	if n := strings.Count(got, "own lines also in"); n != 20 {
		t.Fatalf("want 20 group lines, got %d", n)
	}
	wantIn(t, got, "index.tg: 300 groups of repeated lines, 20 shown; fix the first and build again: one mixin set often clears several")
}

func TestRepeatsClosingLineWithFewGroups(t *testing.T) {
	got := repeatsFail(t, rr("a", l1, l2, l3), rr("b", l1, l2, l3))
	wantIn(t, got, "index.tg: 1 group of repeated lines; fix the first and build again: one mixin set often clears several")
}

func TestRepeatsOutputIsStableAcrossRuns(t *testing.T) {
	body := denseBody()
	first := repeatsFail(t, body...)
	for i := 0; i < 5; i++ {
		if got := repeatsFail(t, body...); got != first {
			t.Fatalf("output changed between runs")
		}
	}
}

// denseBody is 25 roles, each holding all of 25 lines but its own.
func denseBody() []string {
	var body []string
	for i := 0; i < 25; i++ {
		var ls []string
		for j := 0; j < 25; j++ {
			if j != i {
				ls = append(ls, fmt.Sprintf("p%d=1", j))
			}
		}
		body = append(body, rr(fmt.Sprintf("r%d", i), ls...))
	}
	return body
}
