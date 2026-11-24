package vocab

import (
	"strings"

	"go.starlark.net/starlark"
)

// cssLine is one line of a css dict, copied so later changes to the dict
// do not reach it: a declaration (val, text) or a nested rule (sub).
type cssLine struct {
	key    string
	val    starlark.Value
	text   string
	sub    []cssLine
	nested bool
}

// snapshot copies d, nested dicts too. A key that is not a string is
// reported at p and left out.
func (l *loader) snapshot(p []string, d *starlark.Dict) []cssLine {
	var out []cssLine
	for _, kv := range d.Items() {
		key, ok := kv[0].(starlark.String)
		if !ok {
			l.errf(p, "keys must be strings, got %s", kv[0].Type())
			continue
		}
		ln := cssLine{key: string(key), val: kv[1], text: plain(kv[1])}
		if sub, nested := kv[1].(*starlark.Dict); nested {
			ln.nested, ln.sub = true, l.snapshot(extend(p, ln.key), sub)
		}
		out = append(out, ln)
	}
	return out
}

// plain is a css value as written.
func plain(v starlark.Value) string {
	if s, ok := v.(starlark.String); ok {
		return string(s)
	}
	return v.String()
}

// topLines are the property and value of each non-nested line, as written.
func topLines(ls []cssLine) [][2]string {
	var out [][2]string
	for _, ln := range ls {
		if !ln.nested {
			out = append(out, [2]string{ln.key, ln.text})
		}
	}
	return out
}

// lineTexts renders lines for `tidgrid list`: `prop: value`, nested rules as `sel { a: b; c: d }`.
func lineTexts(ls []cssLine) []string {
	var out []string
	for _, ln := range ls {
		if ln.nested {
			out = append(out, ln.key+" { "+strings.Join(lineTexts(ln.sub), "; ")+" }")
		} else {
			out = append(out, ln.key+": "+ln.text)
		}
	}
	return out
}
