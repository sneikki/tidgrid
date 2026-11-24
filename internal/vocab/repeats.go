package vocab

import (
	"slices"
	"strconv"

	"go.starlark.net/starlark"
)

// defaultRepeats are the (lines, roles) pairs the rule fires on unless project(repeats=...) says otherwise.
var defaultRepeats = [][2]int{{3, 2}, {2, 4}}

// repeatGroup is a set of lines (ids, ascending) and every role that has them all.
type repeatGroup struct {
	lines []int
	roles []int // indexes into mix.roles, ascending
}

// setRepeats reads project(repeats=...): a list of (lines, roles) pairs, lines at least 1, roles at least 2.
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
		l.errf([]string{"project", "repeats"}, "must be a list of (lines, roles) pairs, lines at least 1 and roles at least 2, such as [(3, 2), (2, 4)]; [] turns the rule off")
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
		least := int64(1 + i) // lines at least 1, roles at least 2
		if !isInt || !fits || k < least || k > 1<<20 {
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

// checkRepeats reports the groups of own lines that fire a pair: most useful first, at most maxRepeatReports.
func (l *loader) checkRepeats() {
	pairs := l.repeatPairs()
	if len(pairs) == 0 {
		return
	}
	texts, sets := l.ownLineSets()
	groups := pairGroups(sets, pairs)
	rankGroups(groups)
	shown := min(len(groups), maxRepeatReports)
	for _, g := range groups[:shown] {
		l.reportRepeat(g, texts)
	}
	if len(groups) > 0 {
		l.repeatClosing(len(groups), shown)
	}
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

func key(s []int) string {
	var b []byte
	for _, x := range s {
		b = strconv.AppendInt(b, int64(x), 10)
		b = append(b, ',')
	}
	return string(b)
}
