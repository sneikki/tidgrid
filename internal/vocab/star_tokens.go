package vocab

import (
	"fmt"
	"sort"

	"go.starlark.net/starlark"
)

// tokenGroup is the value tokens() returns: `color.ink` and `color["sand-ink"]`
// both give the reference string "{color.ink}". An unknown name is an error.
type tokenGroup struct {
	name  string
	names map[string]bool
}

func newTokenGroup(name string, d *starlark.Dict) *tokenGroup {
	g := &tokenGroup{name: name, names: map[string]bool{}}
	for _, kv := range d.Items() {
		if s, ok := kv[0].(starlark.String); ok {
			g.names[string(s)] = true
		}
	}
	return g
}

func (g *tokenGroup) ref(n string) (starlark.Value, error) {
	if !g.names[n] {
		return nil, fmt.Errorf("token group %s has no token %q", g.name, n)
	}
	return starlark.String("{" + g.name + "." + n + "}"), nil
}

func (g *tokenGroup) Attr(n string) (starlark.Value, error) { return g.ref(n) }

func (g *tokenGroup) AttrNames() []string {
	var out []string
	for n := range g.names {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

func (g *tokenGroup) Get(k starlark.Value) (starlark.Value, bool, error) {
	s, ok := k.(starlark.String)
	if !ok {
		return nil, false, fmt.Errorf("token names are strings, got %s", k.Type())
	}
	v, err := g.ref(string(s))
	return v, err == nil, err
}

func (g *tokenGroup) String() string        { return "tokens(" + g.name + ")" }
func (g *tokenGroup) Type() string          { return "tokens" }
func (g *tokenGroup) Freeze()               {}
func (g *tokenGroup) Truth() starlark.Bool  { return true }
func (g *tokenGroup) Hash() (uint32, error) { return 0, fmt.Errorf("unhashable type: tokens") }
