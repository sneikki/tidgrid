package vocab

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

// mixinState is what the mixin builtin and role includes leave for the checks.
type mixinState struct {
	defs  []*mixinDef
	roles []roleUse
}

type roleRef struct{ kind, name string }

// roleUse is one role that was written, for the hand-written-set check.
type roleUse struct {
	p    []string
	incl []*setDef
	own  map[string]bool // "prop\x00value" of the role's own top-level lines
}

// mixinDef is the value mixin() returns: `type.heading` and `type["heading"]` give a set.
type mixinDef struct {
	name, use string
	sets      []*setDef
	byName    map[string]*setDef
}

// setDef is a tagged set; only include accepts it.
type setDef struct {
	mixin *mixinDef
	name  string
	css   *starlark.Dict
	lines []string
	tops  [][2]string // top-level property and value as written
	users []roleRef
}

func (s *setDef) label() string { return s.mixin.name + "." + s.name }

func (s *setDef) String() string        { return "mixin set " + s.label() }
func (s *setDef) Type() string          { return "mixin set" }
func (s *setDef) Freeze()               {}
func (s *setDef) Truth() starlark.Bool  { return true }
func (s *setDef) Hash() (uint32, error) { return 0, fmt.Errorf("unhashable type: mixin set") }

func (m *mixinDef) set(n string) (starlark.Value, error) {
	s, ok := m.byName[n]
	if !ok {
		return nil, fmt.Errorf("mixin %s has no set %q", m.name, n)
	}
	return s, nil
}

func (m *mixinDef) Attr(n string) (starlark.Value, error) { return m.set(n) }

func (m *mixinDef) AttrNames() []string {
	out := make([]string, len(m.sets))
	for i, s := range m.sets {
		out[i] = s.name
	}
	return out
}

func (m *mixinDef) Get(k starlark.Value) (starlark.Value, bool, error) {
	s, ok := k.(starlark.String)
	if !ok {
		return nil, false, fmt.Errorf("set names are strings, got %s", k.Type())
	}
	v, err := m.set(string(s))
	return v, err == nil, err
}

func (m *mixinDef) String() string        { return "mixin(" + m.name + ")" }
func (m *mixinDef) Type() string          { return "mixin" }
func (m *mixinDef) Freeze()               {}
func (m *mixinDef) Truth() starlark.Bool  { return true }
func (m *mixinDef) Hash() (uint32, error) { return 0, fmt.Errorf("unhashable type: mixin") }

func (e *evaluator) mixin(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var name, use string
	var sets *starlark.Dict
	if err := starlark.UnpackArgs("mixin", args, kw, "name", &name, "use?", &use, "sets?", &sets); err != nil {
		return nil, err
	}
	l := e.l
	p := []string{"mixin", name}
	if first, dup := l.pos[pathKey(p)]; dup {
		return nil, fmt.Errorf("%s: declared twice (%s and %s)", showPath(p), first, at(th))
	}
	l.pos[pathKey(p)] = at(th)
	m := &mixinDef{name: name, use: use, byName: map[string]*setDef{}}
	if strings.TrimSpace(use) == "" {
		l.errf(extend(p, "use"), "required: when the mixin applies")
	}
	if sets == nil {
		l.errf(extend(p, "sets"), "required: a dict of named css sets")
	} else {
		l.addSets(p, m, sets)
	}
	l.mix.defs = append(l.mix.defs, m)
	return m, nil
}

func (l *loader) addSets(p []string, m *mixinDef, sets *starlark.Dict) {
	for _, kv := range sets.Items() {
		key, ok := kv[0].(starlark.String)
		d, isDict := kv[1].(*starlark.Dict)
		if !ok || !isDict {
			l.errf(p, "sets: want names mapped to css dicts")
			continue
		}
		s := &setDef{mixin: m, name: string(key), css: d, lines: lineTexts(d), tops: topLines(d)}
		m.sets = append(m.sets, s)
		m.byName[s.name] = s
	}
}

// plain is a css value as written.
func plain(v starlark.Value) string {
	if s, ok := v.(starlark.String); ok {
		return string(s)
	}
	return v.String()
}

// topLines are the property and value of each non-nested line of d, as written.
func topLines(d *starlark.Dict) [][2]string {
	var out [][2]string
	for _, kv := range d.Items() {
		if _, nested := kv[1].(*starlark.Dict); !nested {
			out = append(out, [2]string{plain(kv[0]), plain(kv[1])})
		}
	}
	return out
}

// lineTexts renders d for `tidgrid list`: `prop: value`, nested rules as `sel { a: b; c: d }`.
func lineTexts(d *starlark.Dict) []string {
	var out []string
	for _, kv := range d.Items() {
		if sub, ok := kv[1].(*starlark.Dict); ok {
			out = append(out, plain(kv[0])+" { "+strings.Join(lineTexts(sub), "; ")+" }")
		} else {
			out = append(out, plain(kv[0])+": "+plain(kv[1]))
		}
	}
	return out
}
