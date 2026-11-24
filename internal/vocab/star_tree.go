package vocab

import (
	"fmt"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
)

// splitPos splits a "file:line" position; both are zero when pos is empty.
func splitPos(pos string) (string, int) {
	i := strings.LastIndex(pos, ":")
	if i < 0 {
		return "", 0
	}
	n, _ := strconv.Atoi(pos[i+1:])
	return pos[:i], n
}

// declare stores v at path p and remembers where it was declared. A second
// declaration of the same path is an error naming both places.
func (l *loader) declare(p []string, v any, at string) error {
	if first, dup := l.pos[pathKey(p)]; dup {
		return fmt.Errorf("%s: declared twice (%s and %s)", showPath(p), first, at)
	}
	l.pos[pathKey(p)] = at
	l.putKey(l.ensureTable(p[:len(p)-1]), p[:len(p)-1], p[len(p)-1], v)
	return nil
}

// ensureTable returns the table at p, creating it and its parents in order.
func (l *loader) ensureTable(p []string) map[string]any {
	cur := l.raw
	for i, k := range p {
		next, ok := cur[k].(map[string]any)
		if !ok {
			next = map[string]any{}
			l.putKey(cur, p[:i], k, next)
		}
		cur = next
	}
	return cur
}

// putKey sets t[k] where t is the table at parent, recording k in key order.
func (l *loader) putKey(t map[string]any, parent []string, k string, v any) {
	if _, ok := t[k]; !ok {
		l.kids[pathKey(parent)] = append(l.kids[pathKey(parent)], k)
	}
	t[k] = v
}

// fill copies the dict d into table t at path p, keeping the dict's key order.
func (l *loader) fill(p []string, t map[string]any, d *starlark.Dict) error {
	for _, kv := range d.Items() {
		key, ok := kv[0].(starlark.String)
		if !ok {
			return fmt.Errorf("%s: keys must be strings, got %s", showPath(p), kv[0].Type())
		}
		v, err := l.cell(extend(p, string(key)), kv[1])
		if err != nil {
			return err
		}
		l.putKey(t, p, string(key), v)
	}
	return nil
}

// cell converts a Starlark value to what the TOML reader would have produced.
// A value of any other type becomes nil, which the load steps reject as "must be a string".
func (l *loader) cell(p []string, v starlark.Value) (any, error) {
	switch t := v.(type) {
	case starlark.String:
		return string(t), nil
	case starlark.Int:
		n, ok := t.Int64()
		if !ok {
			return nil, fmt.Errorf("%s: number out of range", showPath(p))
		}
		return n, nil
	case starlark.Float:
		return float64(t), nil
	case *starlark.Dict:
		sub := map[string]any{}
		return sub, l.fill(p, sub, t)
	}
	return nil, nil
}

// newTable makes a table for path p from d, with use first when it is set.
func (l *loader) newTable(p []string, use string, d *starlark.Dict) (map[string]any, error) {
	t := map[string]any{}
	if use != "" {
		l.putKey(t, p, "use", use)
	}
	if d == nil {
		return t, nil
	}
	return t, l.fill(p, t, d)
}
