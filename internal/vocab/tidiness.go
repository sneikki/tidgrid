package vocab

import (
	"embed"
	"fmt"
	"regexp"
	"strings"

	"go.starlark.net/starlark"
)

//go:embed tidiness/minimal.css tidiness/full.css
var tidinessFS embed.FS

// TidinessLevels are the values tidiness() takes, in the order their files are written.
var TidinessLevels = []string{"minimal", "full"}

func (e *evaluator) tidiness(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var level starlark.Value
	if err := starlark.UnpackArgs("tidiness", args, kw, "level", &level); err != nil {
		return nil, err
	}
	s, ok := level.(starlark.String)
	if !ok || (string(s) != "minimal" && string(s) != "full") {
		return nil, fmt.Errorf(`tidiness: level must be "minimal" or "full", got %s`, level.String())
	}
	if e.l.tidinessAt != "" {
		return nil, fmt.Errorf("tidiness: called twice (%s and %s)", e.l.tidinessAt, at(th))
	}
	e.l.tidinessAt = at(th)
	e.l.v.Tidiness = string(s)
	return starlark.None, nil
}

// tidinessRules are the rules of the chosen level: minimal's, then full's.
func (v *Vocabulary) tidinessRules() []*Rule {
	var out []*Rule
	for _, lvl := range TidinessLevels {
		src, _ := tidinessFS.ReadFile("tidiness/" + lvl + ".css")
		out = append(out, parseRules(string(src))...)
		if lvl == v.Tidiness {
			break
		}
	}
	if v.Tidiness == "" {
		return nil
	}
	return out
}

var cssComment = regexp.MustCompile(`(?s)/\*.*?\*/`)

// parseRules reads the embedded stylesheets: rules, @-blocks and declarations.
// Comments go; blank runs and line breaks inside a selector list collapse.
func parseRules(src string) []*Rule {
	p := &cssParser{src: cssComment.ReplaceAllString(src, "")}
	return p.rules()
}

type cssParser struct {
	src string
	pos int
}

// rules reads blocks up to the closing brace (consumed) or the end.
func (p *cssParser) rules() []*Rule {
	var out []*Rule
	for p.pos < len(p.src) {
		key, delim := p.chunk()
		switch delim {
		case '{':
			out = append(out, p.block(key))
		case '}':
			return out
		}
	}
	return out
}

// block reads the body of the rule whose prelude is key; its '{' is consumed.
func (p *cssParser) block(key string) *Rule {
	r := &Rule{Key: strings.ReplaceAll(strings.Join(strings.Fields(key), " "), " ,", ",")}
	p.body(r)
	return r
}

// body fills r with declarations and nested blocks up to its closing brace.
func (p *cssParser) body(r *Rule) {
	for p.pos < len(p.src) {
		text, delim := p.chunk()
		switch delim {
		case ';':
			r.Nodes = append(r.Nodes, declNode(text))
		case '{':
			r.Nodes = append(r.Nodes, Node{Sub: p.block(text)})
		default:
			if strings.TrimSpace(text) != "" {
				r.Nodes = append(r.Nodes, declNode(text))
			}
			return
		}
	}
}

func declNode(text string) Node {
	prop, value, _ := strings.Cut(text, ":")
	return Node{Prop: strings.TrimSpace(prop), Value: strings.Join(strings.Fields(value), " ")}
}

// chunk reads up to a '{', ';' or '}' outside parentheses and consumes it.
func (p *cssParser) chunk() (string, byte) {
	start, depth := p.pos, 0
	for ; p.pos < len(p.src); p.pos++ {
		switch c := p.src[p.pos]; {
		case c == '(':
			depth++
		case c == ')':
			depth--
		case depth == 0 && (c == '{' || c == ';' || c == '}'):
			p.pos++
			return p.src[start : p.pos-1], c
		}
	}
	return p.src[start:], 0
}
