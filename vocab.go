// Shared API (other commands build on these; do not change signatures):
//
//	load(path) (*Vocabulary, []error)  parse + validate; errors name the table path and key
//	Vocabulary.Project                 Prefix, Out, Scan (resolved against the vocabulary dir), Helpers
//	Vocabulary.Breakpoints             []Breakpoint, file order
//	Vocabulary.Tokens                  []TokenGroup{Name, Tokens []Token}, file order
//	Vocabulary.Roles                   []Role{Kind, Name, Use, Class, Rule}, file order
//	Vocabulary.Elements, Keyframes     file order; Rule/Node hold resolved CSS (Raw = as written)
//	Vocabulary.AllowClass, AllowCSS    []Allow{Name, Reason}, file order (AllowCSS paths as written)
//	Vocabulary.Atoms                   []Atom, README table order (no breakpoint variants)
//	Vocabulary.Classes()               every valid class name, in output order
//	Vocabulary.Valid(class)            is the class in the vocabulary
//	Vocabulary.Takes(family)           allowed values of "tg-gap" / "md:tg-gap" / "tg-surface"
//	Vocabulary.render()                the full stylesheet, byte-deterministic
package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Project struct {
	Prefix  string
	Out     string
	Scan    []string
	Helpers []string
}

type Breakpoint struct{ Name, Value string }
type Token struct{ Name, Value, CSS string } // Value as written, CSS with references resolved
type TokenGroup struct {
	Name   string
	Tokens []Token
}
type Allow struct{ Name, Reason string }

// Node is a declaration (Prop, Value, Raw) or a nested rule (Sub).
type Node struct {
	Prop, Value, Raw string
	Sub              *Rule
}

// Rule: Key is the rendered selector or at-rule prelude, Written the key as in the file.
type Rule struct {
	Key, Written string
	Nodes        []Node
}

type Role struct {
	Kind, Name, Use, Class string
	Rule                   *Rule
}
type Element struct {
	Selector, Use string
	Rule          *Rule
}
type Keyframe struct {
	Name  string
	Stops []*Rule
}

type Vocabulary struct {
	Path        string
	Project     Project
	Breakpoints []Breakpoint
	Tokens      []TokenGroup
	Elements    []Element
	Keyframes   []Keyframe
	Roles       []Role
	AllowClass  []Allow
	AllowCSS    []Allow
	Atoms       []Atom

	classes []string
	valid   map[string]bool
	takes   map[string][]string
	bpValue map[string]string
}

func (v *Vocabulary) Classes() []string   { return v.classes }
func (v *Vocabulary) Valid(c string) bool { return v.valid[c] }

// Takes returns the allowed values of a class family, e.g. "tg-gap" -> sm md.
func (v *Vocabulary) Takes(family string) ([]string, bool) {
	if i := strings.Index(family, ":"); i >= 0 {
		family = family[i+1:]
	}
	t, ok := v.takes[family]
	return t, ok
}

type loader struct {
	v    *Vocabulary
	raw  map[string]any
	kids map[string][]string
	errs []error
	tok  map[string]bool // "group.name"
	kf   map[string]bool
}

func jk(p []string) string { return strings.Join(p, "\x00") }

var bare = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func show(p []string) string {
	q := make([]string, len(p))
	for i, k := range p {
		if bare.MatchString(k) {
			q[i] = k
		} else {
			q[i] = strconv.Quote(k)
		}
	}
	return strings.Join(q, ".")
}

func (l *loader) errf(p []string, f string, a ...any) {
	l.errs = append(l.errs, fmt.Errorf("%s: %s", show(p), fmt.Sprintf(f, a...)))
}

func (l *loader) children(p ...string) []string { return l.kids[jk(p)] }

func (l *loader) table(p ...string) map[string]any {
	var cur any = l.raw
	for _, k := range p {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[k]
	}
	m, _ := cur.(map[string]any)
	return m
}

func scalar(x any) (string, bool) {
	switch t := x.(type) {
	case string:
		return t, true
	case int64:
		return strconv.FormatInt(t, 10), true
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), true
	}
	return "", false
}

var refRe = regexp.MustCompile(`\{([A-Za-z0-9_-]+)\.([A-Za-z0-9_-]+)\}`)

// resolve substitutes {group.name}, {keyframes.x}, {breakpoints.x}.
func (l *loader) resolve(p []string, s string) string {
	pre := l.v.Project.Prefix
	return refRe.ReplaceAllStringFunc(s, func(m string) string {
		sm := refRe.FindStringSubmatch(m)
		g, n := sm[1], sm[2]
		switch g {
		case "keyframes":
			if l.kf[n] {
				return pre + "-" + n
			}
		case "breakpoints":
			if x, ok := l.v.bpValue[n]; ok {
				return x
			}
		default:
			if l.tok[g+"."+n] {
				return "var(--" + pre + "-" + g + "-" + n + ")"
			}
		}
		l.errf(p, "unknown reference %s", m)
		return m
	})
}

const units = `px|rem|em|ex|ch|lh|rlh|vw|vh|vmin|vmax|vi|vb|svw|svh|lvw|lvh|dvw|dvh|cqw|cqh|cqi|cqb|cqmin|cqmax|cm|mm|q|in|pt|pc|ms|s`

var (
	hexRe    = regexp.MustCompile(`#[0-9A-Fa-f]{3,8}\b`)
	colorFn  = regexp.MustCompile(`(?i)\b(rgb|rgba|hsl|hsla|hwb|lab|lch|oklab|oklch|color)\(`)
	unitRe   = regexp.MustCompile(`(?i)(?:^|[^\w.])-?(?:\d+\.?\d*|\.\d+)(?:` + units + `)\b`)
	quotedRe = regexp.MustCompile(`"[^"]*"|'[^']*'`)
)

// rawProblem names a raw colour, length or time in a value, or "".
func rawProblem(s string) string {
	s = quotedRe.ReplaceAllString(refRe.ReplaceAllString(s, "_"), `""`)
	if m := hexRe.FindString(s); m != "" {
		return "raw colour " + m
	}
	if m := colorFn.FindString(s); m != "" {
		return "raw colour " + m + ")"
	}
	if m := unitRe.FindString(s); m != "" {
		return "raw length or time " + strings.TrimLeft(m, " ,(") + " (move it to [tokens])"
	}
	return ""
}

var classRe = regexp.MustCompile(`\.((?:[A-Za-z0-9_-]+:)?PREFIX-[A-Za-z0-9_-]+(?:\([^)]*\))?)`)

// selector validates the vocabulary classes in a selector and escapes them.
func (l *loader) selector(p []string, key string) string {
	re := regexp.MustCompile(strings.Replace(classRe.String(), "PREFIX", regexp.QuoteMeta(l.v.Project.Prefix), 1))
	return re.ReplaceAllStringFunc(key, func(m string) string {
		c := m[1:]
		if !l.v.valid[c] {
			l.errf(p, "unknown class %s%s", c, l.hint(c))
			return m
		}
		return "." + esc(c)
	})
}

func (l *loader) hint(c string) string {
	fam := c
	if i := strings.Index(fam, "("); i >= 0 {
		fam = fam[:i]
	}
	if t, ok := l.v.Takes(fam); ok && len(t) > 0 {
		return "; " + fam + " takes: " + strings.Join(t, " ")
	}
	return ""
}

// rule builds a nested rule from table t at path p.
func (l *loader) rule(p []string, t map[string]any, key, written string, top bool) *Rule {
	r := &Rule{Key: key, Written: written}
	for _, k := range l.children(p...) {
		kp := append(append([]string{}, p...), k)
		val := t[k]
		if sub, ok := val.(map[string]any); ok {
			r.Nodes = append(r.Nodes, Node{Sub: l.rule(kp, sub, l.subKey(kp, k), k, false)})
			continue
		}
		if top && k == "use" {
			continue
		}
		s, ok := scalar(val)
		if !ok {
			l.errf(kp, "value must be a string")
			continue
		}
		if pr := rawProblem(s); pr != "" {
			l.errf(kp, "%s", pr)
		}
		r.Nodes = append(r.Nodes, Node{Prop: k, Raw: s, Value: l.resolve(kp, s)})
	}
	return r
}

var oneWord = regexp.MustCompile(`^@[A-Za-z0-9_-]+$`)

func (l *loader) subKey(p []string, k string) string {
	if strings.HasPrefix(k, "@") {
		name := k[1:]
		if bp, ok := l.v.bpValue[name]; ok {
			return "@media (min-width: " + bp + ")"
		}
		if strings.HasPrefix(name, "below-") {
			if bp, ok := l.v.bpValue[name[6:]]; ok {
				return "@media (width < " + bp + ")"
			}
			l.errf(p, "unknown breakpoint %q", name[6:])
			return k
		}
		if oneWord.MatchString(k) && k != "@starting-style" {
			l.errf(p, "unknown breakpoint %q", name)
			return k
		}
		return l.resolve(p, k)
	}
	if !strings.Contains(k, "&") {
		l.errf(p, "selector must contain &")
		return k
	}
	return l.selector(p, k)
}

func load(path string) (*Vocabulary, []error) {
	l := &loader{raw: map[string]any{}, kids: map[string][]string{}, tok: map[string]bool{}, kf: map[string]bool{}}
	v := &Vocabulary{Path: path, valid: map[string]bool{}, takes: map[string][]string{}, bpValue: map[string]string{}}
	l.v = v
	md, err := toml.DecodeFile(path, &l.raw)
	if err != nil {
		return nil, []error{err}
	}
	seen := map[string]bool{}
	for _, k := range md.Keys() {
		for i := 1; i <= len(k); i++ { // implicit parent tables are not in Keys
			par, full := jk(k[:i-1]), jk(k[:i])
			if !seen[full] {
				seen[full] = true
				l.kids[par] = append(l.kids[par], k[i-1])
			}
		}
	}
	for _, k := range l.children() {
		switch k {
		case "project", "breakpoints", "tokens", "element", "role", "keyframes", "allow":
		default:
			l.errf([]string{k}, "unknown table")
		}
	}

	// project
	pr := Project{Prefix: "tg", Scan: []string{"."}, Helpers: []string{"cn", "clsx", "cva", "classNames", "twMerge"}}
	dir := filepath.Dir(path)
	pt := l.table("project")
	if s, ok := pt["prefix"].(string); ok {
		pr.Prefix = s
	}
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(pr.Prefix) {
		l.errf([]string{"project", "prefix"}, "must be lowercase letters, digits, dashes")
	}
	if s, ok := pt["out"].(string); ok && s != "" {
		pr.Out = filepath.Join(dir, s)
	} else {
		l.errf([]string{"project", "out"}, "required")
	}
	strs := func(key string, def []string) []string {
		a, ok := pt[key].([]any)
		if !ok {
			return def
		}
		var o []string
		for _, x := range a {
			if s, ok := x.(string); ok {
				o = append(o, s)
			} else {
				l.errf([]string{"project", key}, "must be a list of strings")
			}
		}
		return o
	}
	pr.Scan = strs("scan", pr.Scan)
	for i, s := range pr.Scan {
		pr.Scan[i] = filepath.Join(dir, s)
	}
	pr.Helpers = strs("helpers", pr.Helpers)
	v.Project = pr

	// breakpoints
	for _, k := range l.children("breakpoints") {
		s, ok := l.raw["breakpoints"].(map[string]any)[k].(string)
		if !ok {
			l.errf([]string{"breakpoints", k}, "value must be a string")
			continue
		}
		v.Breakpoints = append(v.Breakpoints, Breakpoint{k, s})
		v.bpValue[k] = s
	}

	// token and keyframe names first, so references resolve in any order
	for _, g := range l.children("tokens") {
		if g == "keyframes" || g == "breakpoints" {
			l.errf([]string{"tokens", g}, "reserved group name")
		}
		if l.table("tokens", g) == nil {
			l.errf([]string{"tokens", g}, "must be a table of name = value")
			continue
		}
		for _, n := range l.children("tokens", g) {
			l.tok[g+"."+n] = true
		}
	}
	for _, k := range l.children("keyframes") {
		l.kf[k] = true
	}

	// tokens
	deps := map[string][]string{}
	for _, g := range l.children("tokens") {
		gt := l.table("tokens", g)
		if gt == nil {
			continue
		}
		tg := TokenGroup{Name: g}
		for _, n := range l.children("tokens", g) {
			p := []string{"tokens", g, n}
			s, ok := scalar(gt[n])
			if !ok {
				l.errf(p, "value must be a string")
				continue
			}
			for _, m := range refRe.FindAllStringSubmatch(s, -1) {
				deps[g+"."+n] = append(deps[g+"."+n], m[1]+"."+m[2])
			}
			tg.Tokens = append(tg.Tokens, Token{n, s, l.resolve(p, s)})
		}
		v.Tokens = append(v.Tokens, tg)
	}
	l.cycles(deps)

	// atoms and class set
	names := map[string][]string{}
	for _, g := range v.Tokens {
		for _, t := range g.Tokens {
			names[g.Name] = append(names[g.Name], t.Name)
		}
	}
	v.Atoms = atoms(v.Project.Prefix, names)
	add := func(c string) {
		if !v.valid[c] {
			v.valid[c] = true
			v.classes = append(v.classes, c)
		}
	}
	for _, a := range v.Atoms {
		add(a.Class)
		fam := v.Project.Prefix + "-" + a.Family
		v.takes[fam] = v.takes[fam] // present even when empty
		if a.Value != "" {
			v.takes[fam] = append(v.takes[fam], a.Value)
		}
	}
	for _, b := range v.Breakpoints {
		for _, a := range v.Atoms {
			add(b.Name + ":" + a.Class)
		}
	}
	// roles: names first
	for _, kind := range l.children("role") {
		kt := l.table("role", kind)
		if kt == nil {
			l.errf([]string{"role", kind}, "must be a table of roles")
			continue
		}
		if isAtomName(kind) {
			l.errf([]string{"role", kind}, "kind collides with atom %s-%s", v.Project.Prefix, kind)
			continue
		}
		if !regexp.MustCompile(`^[A-Za-z0-9_-]+$`).MatchString(kind) {
			l.errf([]string{"role", kind}, "kind must be letters, digits, dashes")
			continue
		}
		fam := v.Project.Prefix + "-" + kind
		v.takes[fam] = nil
		for _, n := range l.children("role", kind) {
			if l.table("role", kind, n) == nil {
				l.errf([]string{"role", kind, n}, "must be a role table")
				continue
			}
			if !bare.MatchString(n) {
				l.errf([]string{"role", kind, n}, "name must be letters, digits, dashes")
				continue
			}
			add(fam + "(" + n + ")")
			v.takes[fam] = append(v.takes[fam], n)
		}
	}

	// elements, keyframes, roles
	for _, sel := range l.children("element") {
		p := []string{"element", sel}
		t := l.table(p...)
		if t == nil {
			l.errf(p, "must be a table")
			continue
		}
		e := Element{Selector: sel}
		e.Use, _ = t["use"].(string)
		e.Rule = l.rule(p, t, sel, sel, true)
		v.Elements = append(v.Elements, e)
	}
	for _, name := range l.children("keyframes") {
		p := []string{"keyframes", name}
		t := l.table(p...)
		if t == nil {
			l.errf(p, "must be a table of stops")
			continue
		}
		kf := Keyframe{Name: name}
		for _, stop := range l.children(p...) {
			sp := append(append([]string{}, p...), stop)
			st := l.table(sp...)
			if st == nil {
				l.errf(sp, "must be a table")
				continue
			}
			kf.Stops = append(kf.Stops, l.rule(sp, st, stop, stop, false))
		}
		v.Keyframes = append(v.Keyframes, kf)
	}
	for _, kind := range l.children("role") {
		if l.table("role", kind) == nil || isAtomName(kind) {
			continue
		}
		for _, name := range l.children("role", kind) {
			p := []string{"role", kind, name}
			t := l.table(p...)
			if t == nil || !bare.MatchString(name) || !bare.MatchString(kind) {
				continue
			}
			class := v.Project.Prefix + "-" + kind + "(" + name + ")"
			r := Role{Kind: kind, Name: name, Class: class}
			if u, ok := t["use"].(string); !ok || strings.TrimSpace(u) == "" {
				l.errf(append(p, "use"), "required: when the role applies")
			} else {
				r.Use = u
			}
			r.Rule = l.rule(p, t, "."+esc(class), class, true)
			v.Roles = append(v.Roles, r)
		}
	}

	// allow
	for _, kind := range []string{"class", "css"} {
		for _, k := range l.children("allow", kind) {
			p := []string{"allow", kind, k}
			s, _ := l.table("allow", kind)[k].(string)
			if strings.TrimSpace(s) == "" {
				l.errf(p, "reason required")
				continue
			}
			if kind == "class" {
				v.AllowClass = append(v.AllowClass, Allow{k, s})
			} else {
				v.AllowCSS = append(v.AllowCSS, Allow{k, s})
			}
		}
	}
	if len(l.errs) > 0 {
		return nil, l.errs
	}
	return v, nil
}

// cycles reports token reference loops.
func (l *loader) cycles(deps map[string][]string) {
	state := map[string]int{}
	var keys []string
	for k := range deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var visit func(k string) bool
	visit = func(k string) bool {
		switch state[k] {
		case 1:
			return true
		case 2:
			return false
		}
		state[k] = 1
		for _, d := range deps[k] {
			if visit(d) {
				return true
			}
		}
		state[k] = 2
		return false
	}
	for _, k := range keys {
		if visit(k) {
			g, n, _ := strings.Cut(k, ".")
			l.errf([]string{"tokens", g, n}, "reference loop")
			return
		}
	}
}
