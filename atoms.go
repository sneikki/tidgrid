package main

import (
	"strconv"
	"strings"
)

// Atom is one layout class (no breakpoint prefix): Family "gap", Value "sm",
// Class "tg-gap(sm)".
type Atom struct {
	Family string
	Value  string
	Class  string
	Decls  [][2]string
}

type family struct {
	name   string
	static [][2]string // valueless atoms
	props  []string    // properties set by a valued atom
	group  string      // token group supplying values ("space", "size")
	enum   [][2]string // fixed value -> css value
	multi  []multiVal  // fixed value -> declarations
}

type multiVal struct {
	val   string
	decls [][2]string
}

// numbered builds the values of tg-cols and tg-span.
func numbered(list []int, decl func(n string) [][2]string) []multiVal {
	var o []multiVal
	for _, n := range list {
		s := strconv.Itoa(n)
		o = append(o, multiVal{s, decl(s)})
	}
	return o
}

func upTo(n int) []int {
	var o []int
	for i := 1; i <= n; i++ {
		o = append(o, i)
	}
	return o
}

func d(kv ...string) [][2]string {
	var o [][2]string
	for i := 0; i < len(kv); i += 2 {
		o = append(o, [2]string{kv[i], kv[i+1]})
	}
	return o
}

var (
	alignEnum   = d("start", "flex-start", "center", "center", "end", "flex-end", "baseline", "baseline", "stretch", "stretch")
	justifyEnum = d("start", "flex-start", "center", "center", "end", "flex-end", "between", "space-between", "around", "space-around", "evenly", "space-evenly")
	selfEnum    = d("start", "flex-start", "center", "center", "end", "flex-end", "stretch", "stretch")
	textEnum    = d("start", "start", "center", "center", "end", "end")
)

// families is the README atom table, in its order.
var families = []family{
	{name: "row", static: d("display", "flex", "flex-direction", "row")},
	{name: "col", static: d("display", "flex", "flex-direction", "column")},
	{name: "wrap", static: d("flex-wrap", "wrap")},
	{name: "grow", static: d("flex", "1 1 0", "min-width", "0")},
	{name: "fixed", static: d("flex", "none")},
	{name: "center", static: d("margin-inline", "auto")},
	{name: "align", props: []string{"align-items"}, enum: alignEnum},
	{name: "justify", props: []string{"justify-content"}, enum: justifyEnum},
	{name: "self", props: []string{"align-self"}, enum: selfEnum},
	{name: "text-align", props: []string{"text-align"}, enum: textEnum},
	{name: "gap", props: []string{"gap"}, group: "space"},
	{name: "x-gap", props: []string{"column-gap"}, group: "space"},
	{name: "y-gap", props: []string{"row-gap"}, group: "space"},
	{name: "pad", props: []string{"padding"}, group: "space"},
	{name: "pad-x", props: []string{"padding-inline"}, group: "space"},
	{name: "pad-y", props: []string{"padding-block"}, group: "space"},
	{name: "pad-t", props: []string{"padding-top"}, group: "space"},
	{name: "pad-r", props: []string{"padding-right"}, group: "space"},
	{name: "pad-b", props: []string{"padding-bottom"}, group: "space"},
	{name: "pad-l", props: []string{"padding-left"}, group: "space"},
	{name: "width", props: []string{"width"}, group: "size"},
	{name: "min-width", props: []string{"min-width"}, group: "size"},
	{name: "max-width", props: []string{"max-width"}, group: "size"},
	{name: "height", props: []string{"height"}, group: "size"},
	{name: "min-height", props: []string{"min-height"}, group: "size"},
	{name: "cols", multi: numbered(append(upTo(6), 12), func(n string) [][2]string {
		return d("display", "grid", "grid-template-columns", "repeat("+n+", minmax(0, 1fr))")
	})},
	{name: "span", multi: numbered(upTo(12), func(n string) [][2]string { return d("grid-column", "span "+n) })},
	{name: "order", multi: []multiVal{{"first", d("order", "-1")}, {"last", d("order", "1")}}},
	{name: "show", static: d("display", "revert")}, // show then hide: both last, so hide md:show works
	{name: "hide", static: d("display", "none")},
}

func isAtomName(n string) bool {
	for _, f := range families {
		if f.name == n {
			return true
		}
	}
	return false
}

// atoms expands the families for the given token names (file order).
func atoms(prefix string, tokens map[string][]string) []Atom {
	var out []Atom
	for _, f := range families {
		base := prefix + "-" + f.name
		switch {
		case f.static != nil:
			out = append(out, Atom{Family: f.name, Class: base, Decls: f.static})
		case f.multi != nil:
			for _, m := range f.multi {
				out = append(out, Atom{f.name, m.val, base + "(" + m.val + ")", m.decls})
			}
		case f.enum != nil:
			for _, e := range f.enum {
				out = append(out, Atom{f.name, e[0], base + "(" + e[0] + ")", d(f.props[0], e[1])})
			}
		default:
			for _, n := range tokens[f.group] {
				out = append(out, Atom{f.name, n, base + "(" + n + ")",
					d(f.props[0], "var(--"+prefix+"-"+f.group+"-"+n+")")})
			}
		}
	}
	return out
}

// esc escapes a class name for use in a selector.
func esc(c string) string {
	var b strings.Builder
	for i, r := range c {
		switch {
		case r == '(' || r == ')' || r == ':':
			b.WriteByte('\\')
			b.WriteRune(r)
		case i == 0 && r >= '0' && r <= '9':
			b.WriteString("\\3" + string(r) + " ")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
