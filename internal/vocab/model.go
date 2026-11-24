// Package vocab loads and validates the vocabulary file and renders its
// stylesheet.
package vocab

import "strings"

// Project is the [project] table with defaults filled in.
type Project struct {
	Prefix  string
	Out     string
	Scan    []string
	Helpers []string
}

// Breakpoint is one [breakpoints] entry.
type Breakpoint struct{ Name, Value string }

// Token is one entry of a token group; Value as written, CSS with references resolved.
type Token struct{ Name, Value, CSS string }

// TokenGroup is one [tokens.<group>] table.
type TokenGroup struct {
	Name   string
	Tokens []Token
}

// Allow is an [allow] entry: a name (class or css path) and the reason it is allowed.
type Allow struct{ Name, Reason string }

// Node is a declaration (Prop, Value) or a nested rule (Sub). From names the
// mixin set it came from ("type.heading"); it is empty for a role's own lines.
// A nested rule a set opened and the role added to keeps the set's name;
// its own nodes say which line came from where.
type Node struct {
	Prop, Value string
	Sub         *Rule
	From        string
}

// Rule is a block: Key is the rendered selector or at-rule prelude.
type Rule struct {
	Key   string
	Nodes []Node
}

// Role is one [role.<kind>.<name>] table; Use says when the role applies.
type Role struct {
	Kind, Name, Use, Class string
	Rule                   *Rule
	Includes               []string // mixin sets included, as "mixin.set", in include order
	// File is the declaring file's slash path relative to the entry file's
	// directory and Line the line of the call; both are zero for TOML.
	File string
	Line int
}

type element struct {
	Selector, Use string
	Rule          *Rule
}

type keyframe struct {
	Name  string
	Stops []*Rule
}

// Vocabulary is a loaded, validated vocabulary file.
type Vocabulary struct {
	Path        string       // file path as given to Load
	Project     Project      // Prefix, Out, Scan (resolved against the vocabulary dir), Helpers
	Breakpoints []Breakpoint // file order
	Tokens      []TokenGroup // file order
	Tidiness    string       // "", "minimal" or "full": the reset written after :root
	elements    []element
	keyframes   []keyframe
	Mixins      []Mixin // declaration order
	Roles       []Role  // file order
	AllowClass  []Allow // file order
	AllowCSS    []Allow // file order, paths as written
	Atoms       []Atom  // README table order, no breakpoint variants

	valid   map[string]bool
	takes   map[string][]string
	bpValue map[string]string
}

// Valid reports whether c is a class in the vocabulary.
func (v *Vocabulary) Valid(c string) bool { return v.valid[c] }

// Takes returns the allowed values of a class family, e.g. "tg-gap" -> sm md.
func (v *Vocabulary) Takes(family string) ([]string, bool) {
	if i := strings.Index(family, ":"); i >= 0 {
		family = family[i+1:]
	}
	t, ok := v.takes[family]
	return t, ok
}

// RoleKinds returns the role kinds in file order and the roles of each, in file order.
func (v *Vocabulary) RoleKinds() ([]string, map[string][]Role) {
	byKind := map[string][]Role{}
	var kinds []string
	for _, r := range v.Roles {
		if _, ok := byKind[r.Kind]; !ok {
			kinds = append(kinds, r.Kind)
		}
		byKind[r.Kind] = append(byKind[r.Kind], r)
	}
	return kinds, byKind
}

// addClass registers a valid class.
func (v *Vocabulary) addClass(c string) {
	v.valid[c] = true
}
