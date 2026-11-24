package vocab

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// quote is a Starlark string literal; Go's escapes are a subset of Starlark's.
func quote(s string) string { return strconv.Quote(s) }

// call writes fn(lead, use = "...", css = {...}) over several lines, one argument per line after lead.
func (l *loader) call(fn string, p []string, t map[string]any, lead string) string {
	use, hasUse := t["use"].(string)
	css := l.dictExcept(p, t, 1, "use")
	args := []string{lead}
	if hasUse {
		args = append(args, "use = "+quote(use))
	}
	if css != "{}" {
		args = append(args, "css = "+css)
	}
	return fn + "(\n    " + strings.Join(args, ",\n    ") + ",\n)\n\n"
}

// pyReserved are the names a token group variable must not take: Starlark's
// keywords and reserved words, its constants, and the vocabulary builtins.
var pyReserved = strings.Fields(`and as assert async await break class continue def del elif else except
	finally for from global if import in is lambda load nonlocal not or pass raise return try while with yield
	True False None project breakpoints tokens element keyframes role allow_class allow_css source`)

var nonIdent = regexp.MustCompile(`[^A-Za-z0-9_]`)

// safeName turns a group name into a variable name that is a valid identifier,
// not reserved and not in taken; the chosen name is added to taken.
func safeName(group string, taken map[string]bool) string {
	n := nonIdent.ReplaceAllString(group, "_")
	if n == "" || (n[0] >= '0' && n[0] <= '9') {
		n = "_" + n
	}
	for taken[n] || slices.Contains(pyReserved, n) {
		n += "_"
	}
	taken[n] = true
	return n
}

// args writes the keys of t as keyword arguments, in file order.
func (l *loader) args(p []string, t map[string]any, skip string) string {
	var out []string
	for _, k := range l.children(p...) {
		if k != skip {
			out = append(out, strings.ReplaceAll(k, "-", "_")+" = "+l.value(extend(p, k), t[k], 0))
		}
	}
	return strings.Join(out, ", ")
}

func (l *loader) dict(p []string, t map[string]any, depth int) string {
	return l.dictExcept(p, t, depth, "")
}

// dictExcept is the dict literal of t in file order, leaving out the key skip at this level.
func (l *loader) dictExcept(p []string, t map[string]any, depth int, skip string) string {
	pad := strings.Repeat("    ", depth+1)
	var b strings.Builder
	for _, k := range l.children(p...) {
		if k != skip {
			b.WriteString(pad + quote(k) + ": " + l.value(extend(p, k), t[k], depth+1) + ",\n")
		}
	}
	if b.Len() == 0 {
		return "{}"
	}
	return "{\n" + b.String() + strings.Repeat("    ", depth) + "}"
}

// value writes a table as a dict, a list as a list and any scalar as the string written in the file.
func (l *loader) value(p []string, v any, depth int) string {
	switch t := v.(type) {
	case map[string]any:
		return l.dict(p, t, depth)
	case []any:
		items := make([]string, len(t))
		for i, x := range t {
			items[i] = l.value(p, x, depth)
		}
		return "[" + strings.Join(items, ", ") + "]"
	}
	if s, ok := scalar(v); ok {
		return quote(s)
	}
	return "None"
}
