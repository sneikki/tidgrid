package vocab

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var cssComment = regexp.MustCompile(`(?s)/\*.*?\*/`)

// nextDecl finds a second declaration run into the first: a word and a colon
// after whitespace, outside quotes and brackets.
var nextDecl = regexp.MustCompile(`\s[a-zA-Z-]+\s*:`)

// parseRules reads the embedded stylesheets: rules, @-blocks and declarations.
// Comments go; blank runs and line breaks inside a selector list collapse.
// Anything it does not understand is an error, never a silent cut.
func parseRules(src string) ([]*Rule, error) {
	src = cssComment.ReplaceAllString(src, "")
	if strings.Contains(src, "/*") {
		return nil, errors.New("unterminated comment")
	}
	p := &cssParser{src: src}
	var out []*Rule
	for {
		key, delim, err := p.chunk()
		if err != nil {
			return nil, err
		}
		switch delim {
		case 0:
			if strings.TrimSpace(key) != "" {
				return nil, fmt.Errorf("text after the last rule: %q", key)
			}
			return out, nil
		case '{':
			r, err := p.block(key)
			if err != nil {
				return nil, err
			}
			out = append(out, r)
		case ';':
			return nil, fmt.Errorf("top-level statement not understood: %q", strings.TrimSpace(key)+";")
		default:
			return nil, errors.New("unbalanced braces: stray }")
		}
	}
}

type cssParser struct {
	src string
	pos int
}

// block reads the body of the rule whose prelude is key; its '{' is consumed.
func (p *cssParser) block(key string) (*Rule, error) {
	key = strings.Join(strings.Fields(key), " ")
	if key == "" {
		return nil, errors.New("block without a selector")
	}
	r := &Rule{Key: strings.ReplaceAll(key, " ,", ",")}
	for {
		text, delim, err := p.chunk()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", r.Key, err)
		}
		switch delim {
		case 0:
			return nil, fmt.Errorf("unbalanced braces: %s is not closed", r.Key)
		case '{':
			sub, err := p.block(text)
			if err != nil {
				return nil, err
			}
			r.Nodes = append(r.Nodes, Node{Sub: sub})
			continue
		}
		if delim == '}' && strings.TrimSpace(text) == "" {
			return r, nil
		}
		n, err := declNode(text)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", r.Key, err)
		}
		r.Nodes = append(r.Nodes, n)
		if delim == '}' {
			return r, nil
		}
	}
}

func declNode(text string) (Node, error) {
	prop, value, ok := strings.Cut(text, ":")
	prop, value = strings.TrimSpace(prop), strings.Join(strings.Fields(value), " ")
	switch {
	case !ok:
		return Node{}, fmt.Errorf("declaration without a colon: %q", strings.TrimSpace(text))
	case prop == "" || value == "":
		return Node{}, fmt.Errorf("declaration with an empty name or value: %q", strings.TrimSpace(text))
	case nextDecl.MatchString(" " + bare(value)):
		return Node{}, fmt.Errorf("declaration missing its ';': %q", strings.TrimSpace(text))
	}
	return Node{Prop: prop, Value: value}, nil
}

// bare is s without its quoted strings and bracketed parts.
func bare(s string) string {
	var b strings.Builder
	depth, quote := 0, byte(0)
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			depth--
		case depth == 0:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// chunk reads up to a '{', ';' or '}' outside quotes, parentheses and square
// brackets and consumes it; delim is 0 at the end of the input. A quote or
// bracket left open or closed without a match is an error.
func (p *cssParser) chunk() (string, byte, error) {
	start, depth, quote := p.pos, 0, byte(0)
	for ; p.pos < len(p.src); p.pos++ {
		switch c := p.src[p.pos]; {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			if depth--; depth < 0 {
				return "", 0, fmt.Errorf("unmatched %q", c)
			}
		case depth == 0 && (c == '{' || c == ';' || c == '}'):
			p.pos++
			return p.src[start : p.pos-1], c, nil
		}
	}
	if quote != 0 || depth != 0 {
		return "", 0, fmt.Errorf("unclosed quote or bracket near %q", strings.TrimSpace(p.src[start:]))
	}
	return p.src[start:], 0, nil
}
