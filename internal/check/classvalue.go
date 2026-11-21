package check

import (
	"fmt"
	"strings"
)

// How a class value is quoted decides which interpolation syntax it can hold.
const (
	modeAttr = iota // HTML-ish attribute text: {% %}, {{ }}, ${ }, { }
	modeTmpl        // JS template literal: ${ }
	modeJS          // plain JS string
)

// valueReader splits one class value into whitespace-separated names. A name
// with an interpolation inside is built from pieces and cannot be checked.
type valueReader struct {
	f        *fileScan
	tok      strings.Builder
	tokStart int // -1 between names
	built    bool
	inner    [][2]int // string literals inside interpolations of this name
}

// readClassValue reads the class names in f.src[start:end].
func (f *fileScan) readClassValue(start, end, mode int) {
	if f.seen[start] {
		return
	}
	f.seen[start] = true
	r := &valueReader{f: f, tokStart: -1}
	for i := start; i < end; {
		i = r.next(i, end, mode)
	}
	r.flush(end)
}

func (r *valueReader) begin(i int) {
	if r.tokStart < 0 {
		r.tokStart = i
	}
}

// flush ends the current name at offset at.
func (r *valueReader) flush(at int) {
	if r.tokStart < 0 {
		return
	}
	switch {
	case !r.built:
		r.f.judgeClass(r.tok.String(), r.tokStart)
	case r.tok.String() == "\x00": // interpolation standing alone: read its string literals
		for _, rg := range r.inner {
			r.f.readLiterals(rg[0], rg[1])
		}
	default:
		r.f.report(r.f.line(r.tokStart), builtNameMessage(r.f.src[r.tokStart:at]))
	}
	r.tok.Reset()
	r.tokStart, r.built, r.inner = -1, false, nil
}

func builtNameMessage(piece string) string {
	text := strings.Join(strings.Fields(piece), " ")
	if r := []rune(text); len(r) > 60 {
		text = string(r[:60]) + "…"
	}
	return fmt.Sprintf("class name built from pieces: %s; write every whole name out in a lookup table", text)
}

// interp takes the interpolation s[i:e], its brace at o, as one placeholder byte.
func (r *valueReader) interp(i, o, e int, lits bool) int {
	r.begin(i)
	r.built = true
	r.tok.WriteByte(0)
	if lits {
		r.inner = append(r.inner, [2]int{o + 1, e - 1})
	}
	return e
}

// literal takes the byte at i as part of the current name.
func (r *valueReader) literal(i int) int {
	r.begin(i)
	r.tok.WriteByte(r.f.src[i])
	return i + 1
}

// next consumes what starts at i and returns where reading goes on.
func (r *valueReader) next(i, end, mode int) int {
	s := r.f.src
	c := s[i]
	switch {
	case isSpace(c):
		r.flush(i)
		return i + 1
	case mode == modeAttr && strings.HasPrefix(s[i:end], "{%"):
		r.flush(i)
		if j := strings.Index(s[i+2:end], "%}"); j >= 0 {
			return i + j + 4
		}
		return end
	case mode == modeAttr && strings.HasPrefix(s[i:end], "{{"):
		if j := strings.Index(s[i+2:end], "}}"); j >= 0 {
			return r.interp(i, i, i+j+4, false)
		}
	case mode != modeJS && strings.HasPrefix(s[i:end], "${"), mode == modeAttr && c == '{':
		o := i
		if c == '$' {
			o++
		}
		if e := delimEnd(s, o, end, '{', '}'); e > 0 {
			return r.interp(i, o, e, true)
		}
	}
	return r.literal(i)
}

// compared reports whether the literal s[i:j] is an operand of ==, ===, != or !==.
func (f *fileScan) compared(i, j int) bool {
	s := f.src
	b := strings.TrimRight(s[:i], " \t\r\n")
	a := strings.TrimLeft(s[j:], " \t\r\n")
	return strings.HasSuffix(b, "==") || strings.HasSuffix(b, "!=") || strings.HasPrefix(a, "==") || strings.HasPrefix(a, "!=")
}

// readLiterals reads every string literal in s[i:end] as a class value.
func (f *fileScan) readLiterals(i, end int) {
	s := f.src
	for i < end {
		switch c := s[i]; {
		case c == '\'' || c == '"' || c == '`':
			j, ie := skipString(s, i, end)
			mode := modeJS
			if c == '`' {
				mode = modeTmpl
			}
			if !f.compared(i, j) {
				f.readClassValue(i+1, ie, mode)
			}
			i = j
		case c == '/' && i+1 < end && s[i+1] == '/':
			for i < end && s[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < end && s[i+1] == '*':
			if j := strings.Index(s[i+2:end], "*/"); j >= 0 {
				i += j + 4
			} else {
				i = end
			}
		default:
			i++
		}
	}
}
