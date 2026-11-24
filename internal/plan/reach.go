package plan

import (
	"fmt"
	"slices"
	"strings"
)

// refersTo reports whether a declaration value reads custom property prop with var(prop), with or without a fallback.
func refersTo(value, prop string) bool {
	return strings.Contains(value, "var("+prop+")") || strings.Contains(value, "var("+prop+",")
}

// referrers lists the classes, in stylesheet order and once each, that have a declaration reading prop.
func (ix index) referrers(prop string) []string {
	var out []string
	for _, k := range ix.classes.keys {
		it := ix.classes.items[k]
		if slices.Contains(out, it.name) {
			continue
		}
		for _, d := range it.decls {
			if refersTo(d.value, prop) {
				out = append(out, it.name)
				break
			}
		}
	}
	return out
}

// reachLine is the line under a token: how many classes read it and which. Empty when none do.
func (ix index) reachLine(prop string) string {
	refs := ix.referrers(prop)
	if len(refs) == 0 {
		return ""
	}
	return fmt.Sprintf("  reaches %d classes: %s", len(refs), strings.Join(refs, " "))
}
