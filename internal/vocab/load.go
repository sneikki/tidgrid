package vocab

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Load parses and validates the vocabulary file at path; errors name the table path and key.
// A .tg file is evaluated as Starlark, anything else is read as TOML.
func Load(path string) (*Vocabulary, []error) {
	l := newLoader(path)
	if filepath.Ext(path) == ".tg" {
		if errs := l.evalStarlark(path); len(errs) > 0 {
			return nil, errs
		}
	} else if err := l.decodeTOML(path); err != nil {
		return nil, []error{err}
	}
	return l.validate()
}

func (l *loader) decodeTOML(path string) error {
	md, err := toml.DecodeFile(path, &l.raw)
	if err != nil {
		return err
	}
	l.indexKeys(md.Keys())
	return nil
}

// validate runs every load step over the filled raw tables and kids.
func (l *loader) validate() (*Vocabulary, []error) {
	path := l.v.Path
	l.checkTopLevelTables()
	l.loadProject(filepath.Dir(path))
	l.loadBreakpoints()
	l.declareTokenAndKeyframeNames() // before any reference is resolved, so order in the file does not matter
	l.loadTokens()
	l.registerAtomClasses()
	l.registerRoleClasses()
	l.loadElements()
	l.loadKeyframes()
	l.checkSets()
	l.loadRoles()
	l.loadAllow()
	l.exportMixins()
	if len(l.errs) > 0 {
		return nil, l.errs
	}
	return l.v, nil
}

// loader holds the decoded file and the state shared by the load steps.
type loader struct {
	v     *Vocabulary
	raw   map[string]any
	kids  map[string][]string // table path -> child keys in file order
	errs  []error
	pos   map[string]string // table path -> "file:line" of the declaring call (Starlark only)
	known map[string]bool   // "group.name" of every declared token
	// written maps "group.name" to the token value as written.
	written map[string]string
	// tidinessAt is "file:line" of the tidiness() call, empty before it.
	tidinessAt string
	// keyframeNames holds every declared keyframes name.
	keyframeNames map[string]bool
	expansions    map[string]*expansion // "group.name" -> literal expansion, see literalToken
	mix           mixinState            // mixins and role includes (Starlark only)
	origin        map[string]string     // role line path -> "mixin.set" it came from
	includes      map[string][]string   // role path -> "mixin.set" included
}

func newLoader(path string) *loader {
	return &loader{
		v:             &Vocabulary{Path: path, valid: map[string]bool{}, takes: map[string][]string{}, bpValue: map[string]string{}},
		raw:           map[string]any{},
		kids:          map[string][]string{},
		pos:           map[string]string{},
		known:         map[string]bool{},
		written:       map[string]string{},
		keyframeNames: map[string]bool{},
		expansions:    map[string]*expansion{},
		origin:        map[string]string{},
		includes:      map[string][]string{},
	}
}

// indexKeys records the children of every table, in file order. md.Keys()
// lacks implicit parent tables, so each prefix of a key path is added.
func (l *loader) indexKeys(keys []toml.Key) {
	seen := map[string]bool{}
	for _, k := range keys {
		for i := 1; i <= len(k); i++ {
			parent, full := pathKey(k[:i-1]), pathKey(k[:i])
			if !seen[full] {
				seen[full] = true
				l.kids[parent] = append(l.kids[parent], k[i-1])
			}
		}
	}
}

func (l *loader) checkTopLevelTables() {
	for _, k := range l.children() {
		switch k {
		case "project", "breakpoints", "tokens", "element", "role", "keyframes", "allow":
		default:
			l.errf([]string{k}, "unknown table")
		}
	}
}

func pathKey(p []string) string { return strings.Join(p, "\x00") }

var bareName = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// showPath renders a table path for an error message, quoting keys that are not bare.
func showPath(p []string) string {
	q := make([]string, len(p))
	for i, k := range p {
		if bareName.MatchString(k) {
			q[i] = k
		} else {
			q[i] = strconv.Quote(k)
		}
	}
	return strings.Join(q, ".")
}

func (l *loader) errf(p []string, f string, a ...any) {
	msg := fmt.Sprintf("%s: %s", showPath(p), fmt.Sprintf(f, a...))
	if at := l.declaredAt(p); at != "" {
		msg = at + ": " + msg
	}
	l.errs = append(l.errs, fmt.Errorf("%s", msg))
}

// declaredAt is the position of the nearest enclosing declaration of path p.
func (l *loader) declaredAt(p []string) string {
	for n := len(p); n > 0; n-- {
		if at, ok := l.pos[pathKey(p[:n])]; ok {
			return at
		}
	}
	return ""
}

func (l *loader) children(p ...string) []string { return l.kids[pathKey(p)] }

// table returns the table at path p, or nil when it is missing or not a table.
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

// scalar reads a string or number as the text written in the file.
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

// extend returns p plus k as a new slice, so appends never share storage.
func extend(p []string, k string) []string {
	return append(append([]string{}, p...), k)
}
