package check

import "strings"

// Scanning helpers over source text: strings, brackets and statement ends.
// None of them knows about classes or styles.

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' }

// skipString returns the index after the string opening at i and the end of
// its content.
func skipString(s string, i, end int) (int, int) {
	q := s[i]
	for k := i + 1; k < end; {
		switch {
		case s[k] == '\\':
			k += 2
		case q == '`' && strings.HasPrefix(s[k:end], "${"):
			if e := delimEnd(s, k+1, end, '{', '}'); e > 0 {
				k = e
			} else {
				k += 2
			}
		case s[k] == q:
			return k + 1, k
		default:
			k++
		}
	}
	return end, end
}

// delimEnd returns the index after the delimiter matching s[i]==open, or -1.
func delimEnd(s string, i, end int, open, close byte) int {
	depth := 0
	for k := i; k < end; {
		switch c := s[k]; {
		case c == '\'' || c == '"' || c == '`':
			k, _ = skipString(s, k, end)
		case c == open:
			depth++
			k++
		case c == close:
			depth--
			k++
			if depth == 0 {
				return k
			}
		default:
			k++
		}
	}
	return -1
}

// attrValueEnd returns the index of the quote closing the value that opens at i,
// skipping template tags and interpolations that may contain the quote.
func attrValueEnd(s string, i int) int {
	q, end := s[i], len(s)
	for k := i + 1; k < end; {
		switch {
		case strings.HasPrefix(s[k:], "{{"), strings.HasPrefix(s[k:], "{%"):
			cl := "}}"
			if s[k+1] == '%' {
				cl = "%}"
			}
			if j := strings.Index(s[k+2:], cl); j >= 0 {
				k += j + 4
			} else {
				k += 2
			}
		case s[k] == '{', strings.HasPrefix(s[k:], "${"):
			o := k
			if s[k] == '$' {
				o++
			}
			if e := delimEnd(s, o, end, '{', '}'); e > 0 {
				k = e
			} else {
				k++
			}
		case s[k] == q:
			return k
		default:
			k++
		}
	}
	return -1
}

// attrEnd returns the end of an unquoted attribute value starting at p.
func attrEnd(s string, p int) int {
	e := p
	for e < len(s) && !isSpace(s[e]) && s[e] != '>' && !strings.HasPrefix(s[e:], "/>") {
		e++
	}
	return e
}

// continues reports whether the statement from i goes on after the newline
// at nl: the line ends in a binary operator or the next one starts with one.
func continues(s string, i, nl int) bool {
	ls := strings.LastIndexByte(s[i:nl], '\n') + 1 + i
	b := strings.TrimRight(s[ls:nl], " \t\r")
	if b != "" && strings.ContainsRune("+-*/%?:,=&|^~!", rune(b[len(b)-1])) &&
		!strings.HasSuffix(b, "++") && !strings.HasSuffix(b, "--") && !strings.HasSuffix(b, "*/") {
		return true
	}
	a := strings.TrimLeft(s[nl:], " \t\r\n")
	if strings.HasPrefix(a, "//") || strings.HasPrefix(a, "/*") || a == "" {
		return false
	}
	return strings.ContainsRune("?:.+-*/%|&=,", rune(a[0]))
}

// stmtEnd returns the end of the statement that starts at i: a `;` or newline
// outside brackets, or the bracket that closes the enclosing one.
func stmtEnd(s string, i int) int {
	depth := 0
	for k := i; k < len(s); {
		switch c := s[k]; {
		case c == '\'' || c == '"' || c == '`':
			k, _ = skipString(s, k, len(s))
		case c == '(' || c == '[' || c == '{':
			depth++
			k++
		case c == ')' || c == ']' || c == '}':
			if depth--; depth < 0 {
				return k
			}
			k++
		case depth == 0 && c == ';':
			return k
		case depth == 0 && c == '\n':
			if !continues(s, i, k) {
				return k
			}
			k++
		default:
			k++
		}
	}
	return len(s)
}
