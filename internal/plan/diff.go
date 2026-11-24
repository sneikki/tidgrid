package plan

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// Diff lists what changing the stylesheet from old to cur does: classes, then
// tokens, then element rules, one fact per line. It is empty when they match.
// v names tokens and breakpoints as the vocabulary writes them; nil keeps raw names.
func Diff(old, cur string, v *vocab.Vocabulary) []string {
	nm := newNames(v)
	o := buildIndex(parseSheet(old), nm)
	n := buildIndex(parseSheet(cur), nm)
	reach := func(prop string, removed bool) string {
		if removed {
			return o.reachLine(prop)
		}
		return n.reachLine(prop)
	}
	var out []string
	out = append(out, diffTable(o.classes, n.classes, nil)...)
	out = append(out, diffTable(o.tokens, n.tokens, reach)...)
	out = append(out, diffTable(o.elements, n.elements, nil)...)
	for i, l := range out {
		out[i] = nm.rename(l)
	}
	return out
}

// diffTable lists the items added or changed in cur, in its order, then the removed ones in old's order.
// note, when set, adds a line under a changed or removed item; removed says which sheet it reads.
func diffTable(old, cur *table, note func(key string, removed bool) string) []string {
	var out []string
	for _, k := range cur.keys {
		it := cur.items[k]
		prev, ok := old.items[k]
		if !ok {
			out = append(out, "added "+present(it))
			continue
		}
		for _, ch := range diffDecls(prev.decls, it.decls) {
			out = append(out, "changed "+it.label+": "+ch)
			out = appendNote(out, note, k, false)
		}
	}
	for _, k := range old.keys {
		if _, ok := cur.items[k]; !ok {
			out = append(out, "removed "+present(old.items[k]))
			out = appendNote(out, note, k, true)
		}
	}
	return out
}

func appendNote(out []string, note func(string, bool) string, key string, removed bool) []string {
	if note == nil {
		return out
	}
	if s := note(key, removed); s != "" {
		out = append(out, s)
	}
	return out
}

// present is an item's label, with the value for a token.
func present(it *item) string {
	if it.token && len(it.decls) > 0 {
		return it.label + ": " + it.decls[0].value
	}
	return it.label
}

// propValues is a stylesheet block's values grouped by property, in order of first appearance.
type propValues struct {
	order []string
	vals  map[string][]string
}

func groupProps(ds []decl) propValues {
	g := propValues{vals: map[string][]string{}}
	for _, d := range ds {
		if _, ok := g.vals[d.prop]; !ok {
			g.order = append(g.order, d.prop)
		}
		g.vals[d.prop] = append(g.vals[d.prop], d.value)
	}
	return g
}

// diffDecls lists the declaration changes from old to cur. A property with one
// value on each side is a change; otherwise values are added or removed.
func diffDecls(old, cur []decl) []string {
	before, after := groupProps(old), groupProps(cur)
	var out []string
	for _, p := range after.order {
		out = append(out, diffValues(p, before.vals[p], after.vals[p])...)
	}
	for _, p := range before.order {
		if _, ok := after.vals[p]; !ok {
			out = append(out, diffValues(p, before.vals[p], nil)...)
		}
	}
	return out
}

// diffValues compares a property's values in order. Values that only moved
// are a change of the whole list, since the last one wins.
func diffValues(prop string, old, cur []string) []string {
	if slices.Equal(old, cur) {
		return nil
	}
	if len(old) == 1 && len(cur) == 1 {
		return []string{change(prop, old[0], cur[0])}
	}
	var out []string
	for _, v := range subtract(cur, old) {
		out = append(out, "added "+show(prop, v))
	}
	for _, v := range subtract(old, cur) {
		out = append(out, "removed "+show(prop, v))
	}
	if len(out) == 0 {
		return []string{change(prop, strings.Join(old, "; "), strings.Join(cur, "; "))}
	}
	return out
}

// subtract returns the values of a not matched by one of b, counting repeats.
func subtract(a, b []string) []string {
	rest := slices.Clone(b)
	var out []string
	for _, v := range a {
		if i := slices.Index(rest, v); i >= 0 {
			rest = slices.Delete(rest, i, i+1)
		} else {
			out = append(out, v)
		}
	}
	return out
}

func change(prop, from, to string) string {
	if prop == "" {
		return from + " -> " + to
	}
	return fmt.Sprintf("%s: %s -> %s", prop, from, to)
}

func show(prop, value string) string {
	if prop == "" {
		return value
	}
	return prop + ": " + value
}
