package vocab

import (
	"strconv"
	"strings"
)

// quote is a Starlark string literal; Go's escapes are a subset of Starlark's.
func quote(s string) string { return strconv.Quote(s) }

// call writes `fn(lead, use = "...", css = {...})` for the table t at p.
func (l *loader) call(fn string, p []string, t map[string]any, lead string) string {
	use, hasUse := t["use"].(string)
	css := l.dictExcept(p, t, 0, "use")
	var args []string
	args = append(args, lead)
	if hasUse {
		args = append(args, "use = "+quote(use))
	}
	if css != "{}" {
		args = append(args, "css = "+css)
	}
	return fn + "(" + strings.Join(args, ", ") + ")\n\n"
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
