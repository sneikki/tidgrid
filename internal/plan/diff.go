package plan

import (
	"fmt"
	"slices"
)

// Diff lists what changing the stylesheet from old to cur does: classes, then
// tokens, then element rules, one fact per line. It is empty when they match.
func Diff(old, cur string) []string {
	o := buildIndex(parseSheet(old))
	n := buildIndex(parseSheet(cur))
	var out []string
	out = append(out, diffTable(o.classes, n.classes)...)
	out = append(out, diffTable(o.tokens, n.tokens)...)
	out = append(out, diffTable(o.elements, n.elements)...)
	return out
}

// diffTable lists the items added or changed in cur, in its order, then the removed ones in old's order.
func diffTable(old, cur *table) []string {
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
		}
	}
	for _, k := range old.keys {
		if _, ok := cur.items[k]; !ok {
			out = append(out, "removed "+present(old.items[k]))
		}
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
