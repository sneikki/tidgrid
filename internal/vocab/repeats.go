package vocab

import (
	"slices"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
)

// defaultRepeats are the (lines, roles) pairs the rule fires on unless project(repeats=...) says otherwise.
var defaultRepeats = [][2]int{{3, 2}, {2, 4}}

// repeatGroup is a set of lines (ids, ascending) and every role that has them all.
type repeatGroup struct {
	lines []int
	roles []int // indexes into mix.roles, ascending
}

// setRepeats reads project(repeats=...): a list of (lines, roles) pairs of positive integers.
func (l *loader) setRepeats(v starlark.Value) {
	var out [][2]int
	seq, ok := v.(starlark.Sequence)
	if _, isStr := v.(starlark.String); ok && !isStr {
		for it, x := seq.Iterate(), starlark.Value(nil); it.Next(&x); {
			pair, good := repeatPair(x)
			if !good {
				ok = false
				break
			}
			out = append(out, pair)
		}
	} else {
		ok = false
	}
	if !ok {
		l.errf([]string{"project", "repeats"}, "must be a list of (lines, roles) pairs of positive integers, such as [(3, 2), (2, 4)]; [] turns the rule off")
		return
	}
	l.mix.repeats, l.mix.repeatsSet = out, true
}

func repeatPair(x starlark.Value) ([2]int, bool) {
	t, ok := x.(starlark.Tuple)
	if !ok || len(t) != 2 {
		return [2]int{}, false
	}
	var pair [2]int
	for i, n := range t {
		v, isInt := n.(starlark.Int)
		k, fits := v.Int64()
		if !isInt || !fits || k < 1 || k > 1<<20 {
			return pair, false
		}
		pair[i] = int(k)
	}
	return pair, true
}

func (l *loader) repeatPairs() [][2]int {
	if l.mix.repeatsSet {
		return l.mix.repeats
	}
	return defaultRepeats
}

// checkRepeats reports each group of own lines shared by enough roles once, at its first role.
func (l *loader) checkRepeats() {
	pairs := l.repeatPairs()
	if len(pairs) == 0 {
		return
	}
	texts, sets := l.ownLineSets()
	groups := closedGroups(sets, minLines(pairs))
	// file order of the first role, then of the roles, so output never depends on discovery order
	slices.SortFunc(groups, func(a, b repeatGroup) int { return slices.Compare(a.roles, b.roles) })
	for _, g := range groups {
		if len(g.roles) >= 2 && fires(pairs, len(g.lines), len(g.roles)) {
			l.reportRepeat(g, texts)
		}
	}
}

func minLines(pairs [][2]int) int {
	m := pairs[0][0]
	for _, p := range pairs {
		m = min(m, p[0])
	}
	return m
}

func fires(pairs [][2]int, lines, roles int) bool {
	for _, p := range pairs {
		if lines >= p[0] && roles >= p[1] {
			return true
		}
	}
	return false
}

// ownLineSets numbers every distinct own line in first-seen order and gives each role its ids, ascending.
func (l *loader) ownLineSets() ([]string, [][]int) {
	ids := map[string]int{}
	var texts []string
	sets := make([][]int, len(l.mix.roles))
	for i, r := range l.mix.roles {
		for _, ln := range r.ownList {
			k := ln[0] + "\x00" + ln[1]
			id, ok := ids[k]
			if !ok {
				id = len(texts)
				ids[k] = id
				texts = append(texts, ln[0]+": "+ln[1])
			}
			sets[i] = append(sets[i], id)
		}
		slices.Sort(sets[i])
	}
	return texts, sets
}

// closedGroups finds every line set of at least min lines that is the whole
// overlap of the roles holding it: the intersections of role sets, closed
// under intersecting with one more role. Pruning at min is safe because
// intersections only shrink.
func closedGroups(sets [][]int, min int) []repeatGroup {
	seen := map[string]bool{}
	var queue [][]int
	add := func(s []int) {
		k := key(s)
		if len(s) >= min && !seen[k] {
			seen[k] = true
			queue = append(queue, s)
		}
	}
	for _, s := range sets {
		add(s)
	}
	for i := 0; i < len(queue); i++ {
		for _, s := range sets {
			add(intersect(queue[i], s))
		}
	}
	var out []repeatGroup
	for _, lines := range queue {
		out = append(out, repeatGroup{lines: lines, roles: holders(sets, lines)})
	}
	return out
}

func key(s []int) string {
	var b strings.Builder
	for _, x := range s {
		b.WriteString(strconv.Itoa(x) + ",")
	}
	return b.String()
}

func intersect(a, b []int) []int {
	var out []int
	for i, j := 0, 0; i < len(a) && j < len(b); {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			out = append(out, a[i])
			i++
			j++
		}
	}
	return out
}

func holders(sets [][]int, lines []int) []int {
	var out []int
	for i, s := range sets {
		if len(intersect(lines, s)) == len(lines) {
			out = append(out, i)
		}
	}
	return out
}

// reportRepeat writes the error at the group's first role, lines in that role's order.
func (l *loader) reportRepeat(g repeatGroup, texts []string) {
	first := l.mix.roles[g.roles[0]]
	var others, lines []string
	for _, i := range g.roles[1:] {
		p := l.mix.roles[i].p
		others = append(others, p[1]+"."+p[2])
	}
	shared := map[string]bool{}
	for _, id := range g.lines {
		shared[texts[id]] = true
	}
	for _, ln := range first.ownList {
		if t := ln[0] + ": " + ln[1]; shared[t] {
			lines = append(lines, t)
		}
	}
	l.errf(first.p, "own lines also in role %s: %s; make them a mixin set and include it",
		strings.Join(others, ", "), strings.Join(lines, "; "))
}
