package vocab

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// DefaultConfig is the vocabulary file name used when -c is not given.
const DefaultConfig = "tidgrid.toml"

// LoadFile loads the vocabulary at path. On failure it prints each error to
// stderr and returns a nil vocabulary with exit code 2; else the vocabulary and 0.
func LoadFile(path string) (*Vocabulary, int) {
	v, errs := Load(path)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, e)
		}
		return nil, 2
	}
	return v, 0
}

// LoadFlag parses -c and loads that vocabulary with LoadFile; a bad flag returns exit code 2.
func LoadFlag(name string, args []string) (*Vocabulary, int) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	cfg := fs.String("c", DefaultConfig, "vocabulary file")
	if err := fs.Parse(args); err != nil {
		return nil, 2
	}
	return LoadFile(*cfg)
}

// Load parses and validates the vocabulary file at path; errors name the table path and key.
func Load(path string) (*Vocabulary, []error) {
	l := newLoader(path)
	md, err := toml.DecodeFile(path, &l.raw)
	if err != nil {
		return nil, []error{err}
	}
	l.indexKeys(md.Keys())
	l.checkTopLevelTables()
	l.loadProject(filepath.Dir(path))
	l.loadBreakpoints()
	l.declareTokenAndKeyframeNames() // before any reference is resolved, so order in the file does not matter
	l.loadTokens()
	l.registerAtomClasses()
	l.registerRoleClasses()
	l.loadElements()
	l.loadKeyframes()
	l.loadRoles()
	l.loadAllow()
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
	known map[string]bool // "group.name" of every declared token
	// written maps "group.name" to the token value as written.
	written map[string]string
	// keyframeNames holds every declared keyframes name.
	keyframeNames map[string]bool
	expansions    map[string]*expansion // "group.name" -> literal expansion, see literalToken
}

func newLoader(path string) *loader {
	return &loader{
		v:             &Vocabulary{Path: path, valid: map[string]bool{}, takes: map[string][]string{}, bpValue: map[string]string{}},
		raw:           map[string]any{},
		kids:          map[string][]string{},
		known:         map[string]bool{},
		written:       map[string]string{},
		keyframeNames: map[string]bool{},
		expansions:    map[string]*expansion{},
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
	l.errs = append(l.errs, fmt.Errorf("%s: %s", showPath(p), fmt.Sprintf(f, a...)))
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
