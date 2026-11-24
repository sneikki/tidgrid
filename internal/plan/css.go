package plan

import "strings"

// maxDepth caps nesting; deeper blocks are skipped, so a hand-edited file cannot exhaust the stack.
const maxDepth = 64

// decl is one declaration. Prop is empty for a token, whose value is the whole line.
type decl struct {
	prop, value string
}

// block is a rule or at-rule: its prelude as written, its declarations and its nested blocks.
type block struct {
	key   string
	decls []decl
	subs  []*block
}

// parseSheet reads a stylesheet into blocks. It never fails: text it does not
// understand is skipped, and an unclosed block ends at the end of the input.
func parseSheet(src string) []*block {
	p := &parser{src: src}
	var top []*block
	for p.pos < len(p.src) {
		_, subs := p.statements(0)
		top = append(top, subs...)
	}
	return top
}

type parser struct {
	src string
	pos int
}

// statements reads declarations and blocks up to the closing brace of its block
// (consumed) or the end of the input.
func (p *parser) statements(depth int) ([]decl, []*block) {
	var decls []decl
	var subs []*block
	for {
		text, delim := p.chunk()
		switch delim {
		case '{':
			subs = append(subs, p.block(text, depth))
		case ';':
			decls = appendDecl(decls, text)
		default: // '}' or the end of the input
			return appendDecl(decls, text), subs
		}
	}
}

// block reads the body of a block whose prelude is key; its '{' is consumed.
func (p *parser) block(key string, depth int) *block {
	if depth+1 > maxDepth {
		p.skipBody()
		return &block{key: key}
	}
	decls, subs := p.statements(depth + 1)
	return &block{key: key, decls: decls, subs: subs}
}

// skipBody moves past the rest of a block whose '{' is already consumed.
func (p *parser) skipBody() {
	open := 1
	for ; p.pos < len(p.src) && open > 0; p.pos++ {
		switch p.src[p.pos] {
		case '{':
			open++
		case '}':
			open--
		}
	}
}

// chunk reads text up to a '{', ';' or '}' outside quotes and parentheses, and
// consumes that delimiter. Comments are dropped; escapes are kept as written.
// It returns delim 0 at the end of the input.
func (p *parser) chunk() (string, byte) {
	var b strings.Builder
	var quote byte
	paren := 0
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if strings.HasPrefix(p.src[p.pos:], "/*") && quote == 0 {
			p.skipComment()
			continue
		}
		if c == '\\' && p.pos+1 < len(p.src) {
			b.WriteString(p.src[p.pos : p.pos+2])
			p.pos += 2
			continue
		}
		p.pos++
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(':
			paren++
		case c == ')' && paren > 0:
			paren--
		case paren == 0 && (c == '{' || c == ';' || c == '}'):
			return strings.TrimSpace(b.String()), c
		}
		b.WriteByte(c)
	}
	return strings.TrimSpace(b.String()), 0
}

// skipComment moves past a /* comment at the current position.
func (p *parser) skipComment() {
	end := strings.Index(p.src[p.pos+2:], "*/")
	if end < 0 {
		p.pos = len(p.src)
		return
	}
	p.pos += 2 + end + 2
}

// appendDecl adds text as a declaration when it has a property and a colon; anything else is skipped.
func appendDecl(decls []decl, text string) []decl {
	i := strings.Index(text, ":")
	if i < 0 {
		return decls
	}
	prop := strings.TrimSpace(text[:i])
	if prop == "" {
		return decls
	}
	return append(decls, decl{prop: prop, value: strings.TrimSpace(text[i+1:])})
}
