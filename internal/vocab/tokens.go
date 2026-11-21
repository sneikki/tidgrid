package vocab

import (
	"sort"
	"strings"
)

// declareTokenAndKeyframeNames records every token and keyframes name before
// any value is read, so references resolve whatever the order in the file.
func (l *loader) declareTokenAndKeyframeNames() {
	for _, g := range l.children("tokens") {
		if g == "keyframes" || g == "breakpoints" {
			l.errf([]string{"tokens", g}, "reserved group name")
		}
		if !bareName.MatchString(g) {
			l.errf([]string{"tokens", g}, "group name must be letters, digits, - and _")
		}
		if l.table("tokens", g) == nil {
			l.errf([]string{"tokens", g}, "must be a table of name = value")
			continue
		}
		l.declareGroupTokens(g)
	}
	for _, k := range l.children("keyframes") {
		l.keyframeNames[k] = true
	}
}

func (l *loader) declareGroupTokens(group string) {
	for _, n := range l.children("tokens", group) {
		if !bareName.MatchString(n) {
			l.errf([]string{"tokens", group, n}, "token name must be letters, digits, - and _")
		}
		l.known[group+"."+n] = true
		if s, ok := scalar(l.table("tokens", group)[n]); ok {
			l.written[group+"."+n] = s
		}
	}
}

// loadTokens resolves the references in every token value and rejects reference loops.
func (l *loader) loadTokens() {
	deps := map[string][]string{} // "group.name" -> tokens its value refers to
	for _, g := range l.children("tokens") {
		gt := l.table("tokens", g)
		if gt == nil {
			continue
		}
		group := TokenGroup{Name: g}
		for _, n := range l.children("tokens", g) {
			p := []string{"tokens", g, n}
			s, ok := scalar(gt[n])
			if !ok {
				l.errf(p, "value must be a string")
				continue
			}
			for _, m := range refRe.FindAllStringSubmatch(s, -1) {
				deps[g+"."+n] = append(deps[g+"."+n], m[1]+"."+m[2])
			}
			group.Tokens = append(group.Tokens, Token{n, s, l.resolve(p, s)})
		}
		l.v.Tokens = append(l.v.Tokens, group)
	}
	l.reportReferenceLoop(deps)
}

// reportReferenceLoop reports the first token (in name order) that reaches itself through references.
func (l *loader) reportReferenceLoop(deps map[string][]string) {
	const (
		visiting = 1
		done     = 2
	)
	state := map[string]int{}
	var keys []string
	for k := range deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var loops func(k string) bool
	loops = func(k string) bool {
		switch state[k] {
		case visiting:
			return true
		case done:
			return false
		}
		state[k] = visiting
		for _, dep := range deps[k] {
			if loops(dep) {
				return true
			}
		}
		state[k] = done
		return false
	}
	for _, k := range keys {
		if loops(k) {
			g, n, _ := strings.Cut(k, ".")
			l.errf([]string{"tokens", g, n}, "reference loop")
			return
		}
	}
}
