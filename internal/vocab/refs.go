package vocab

import "regexp"

// maxLiteral caps a literal expansion; tokens that each repeat the next
// one would otherwise cost exponential work and output.
const maxLiteral = 4096

var refRe = regexp.MustCompile(`\{([A-Za-z0-9_-]+)\.([A-Za-z0-9_-]+)\}`)

// resolve substitutes {group.name}, {keyframes.x}, {breakpoints.x}.
func (l *loader) resolve(p []string, s string) string {
	pre := l.v.Project.Prefix
	return refRe.ReplaceAllStringFunc(s, func(m string) string {
		sm := refRe.FindStringSubmatch(m)
		g, n := sm[1], sm[2]
		switch g {
		case "keyframes":
			if l.keyframeNames[n] {
				return pre + "-" + n
			}
		case "breakpoints":
			if x, ok := l.v.bpValue[n]; ok {
				return x
			}
		default:
			if l.known[g+"."+n] {
				return "var(--" + pre + "-" + g + "-" + n + ")"
			}
		}
		l.errf(p, "unknown reference %s", m)
		return m
	})
}

// literal is resolve for places that cannot read custom properties (at-rule
// keys): a token becomes its written value, references inside it resolved too.
func (l *loader) literal(p []string, s string) string {
	tooBig := false
	out := refRe.ReplaceAllStringFunc(s, func(m string) string {
		sm := refRe.FindStringSubmatch(m)
		g, n := sm[1], sm[2]
		if _, ok := l.written[g+"."+n]; ok && g != "keyframes" && g != "breakpoints" {
			e := l.literalToken(g + "." + n)
			tooBig = tooBig || e.big
			return e.s
		}
		return l.resolve(p, m)
	})
	if tooBig || len(out) > maxLiteral {
		l.errf(p, "references expand past %d bytes", maxLiteral)
		return ""
	}
	return out
}

// expansion is the memoised literal text of one token.
type expansion struct {
	s    string
	big  bool // cut at maxLiteral, here or in a token it references
	busy bool
}

// literalToken expands one token once and remembers it.
func (l *loader) literalToken(key string) expansion {
	if e, ok := l.expansions[key]; ok {
		if e.busy {
			return expansion{s: "{" + key + "}"} // loop: reported where it is written
		}
		return *e
	}
	e := &expansion{busy: true}
	l.expansions[key] = e
	e.s = refRe.ReplaceAllStringFunc(l.written[key], func(m string) string {
		sm := refRe.FindStringSubmatch(m)
		k := sm[1] + "." + sm[2]
		if x, ok := l.v.bpValue[sm[2]]; ok && sm[1] == "breakpoints" {
			return x
		}
		if _, ok := l.written[k]; !ok || sm[1] == "keyframes" || sm[1] == "breakpoints" {
			return m // unknown inside a token: reported where it is written
		}
		sub := l.literalToken(k)
		e.big = e.big || sub.big
		return sub.s
	})
	if len(e.s) > maxLiteral {
		e.s, e.big = "", true
	}
	e.busy = false
	return *e
}
