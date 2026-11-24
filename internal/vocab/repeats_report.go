package vocab

import (
	"cmp"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

// maxRepeatReports caps the group lines of one build; the closing line gives the total.
const maxRepeatReports = 20

// rankGroups orders the groups most useful first: more roles, then more lines, then file order of the roles.
func rankGroups(groups []repeatGroup) {
	slices.SortFunc(groups, func(a, b repeatGroup) int {
		return cmp.Or(
			cmp.Compare(len(b.roles), len(a.roles)),
			cmp.Compare(len(b.lines), len(a.lines)),
			slices.Compare(a.roles, b.roles),
			slices.Compare(a.lines, b.lines),
		)
	})
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

// repeatClosing is the last line of a repeat report: the total, how many are shown, and the way out.
func (l *loader) repeatClosing(total, shown int) {
	count := fmt.Sprintf("%d groups of repeated lines", total)
	if total > shown {
		count += fmt.Sprintf(", %d shown", shown)
	}
	l.errs = append(l.errs, fmt.Errorf("%s: %s; fix the first and build again: one mixin set often clears several",
		filepath.Base(l.v.Path), count))
}
