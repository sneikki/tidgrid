package vocab

import (
	"regexp"
	"strings"
)

// buildRule builds the rule for table t at path p; topLevel is true for the
// direct table of a role or element, the only place `use` may appear.
func (l *loader) buildRule(p []string, t map[string]any, key string, topLevel bool) *Rule {
	r := &Rule{Key: key}
	for _, k := range l.children(p...) {
		kp := extend(p, k)
		val := t[k]
		if sub, ok := val.(map[string]any); ok {
			n := len(l.errs)
			key := l.nestedKey(kp, k)
			l.dropSetErrs(kp, n)
			r.Nodes = append(r.Nodes, Node{Sub: l.buildRule(kp, sub, key, false), From: l.origin[pathKey(kp)]})
			continue
		}
		if k == "use" {
			if !topLevel {
				l.errf(kp, "use belongs directly in a role or element table, not in a nested rule")
			} else if _, ok := val.(string); !ok {
				l.errf(kp, "use must be a string")
			}
			continue
		}
		n := len(l.errs)
		s, ok := scalar(val)
		if !ok {
			l.errf(kp, "value must be a string")
		} else {
			if pr := rawProblem(k, s); pr != "" {
				l.errf(kp, "%s", pr)
			}
			r.Nodes = append(r.Nodes, Node{Prop: k, Value: l.resolve(kp, s), From: l.origin[pathKey(kp)]})
		}
		l.dropSetErrs(kp, n)
	}
	return r
}

// dropSetErrs forgets the errors since n when the line at p came from a mixin
// set: checkSets reports a set's content once, where the set is written.
func (l *loader) dropSetErrs(p []string, n int) {
	if l.origin[pathKey(p)] != "" {
		l.errs = l.errs[:n]
	}
}

// checkSets validates each set's content once, at the mixin() call.
func (l *loader) checkSets() {
	for _, m := range l.mix.defs {
		for _, s := range m.sets {
			p := []string{"mixin", m.name, s.name}
			if declCount(s.tree) < 2 {
				l.errf(p, "a set needs at least two lines; a single shared value belongs in a token")
				continue
			}
			t := map[string]any{}
			l.addLines(p, t, s.tree, "", map[string]string{})
			l.buildRule(p, t, "", false)
		}
	}
}

// declCount counts every declaration in ls, nested rules included.
func declCount(ls []cssLine) int {
	n := 0
	for _, ln := range ls {
		if ln.nested {
			n += declCount(ln.sub)
		} else {
			n++
		}
	}
	return n
}

var (
	containerName = regexp.MustCompile(`^(@container)\s+[^\s(]+`) // the name is the author's pick, not a colour
	oneWordAtRule = regexp.MustCompile(`^@[A-Za-z0-9_-]+$`)
	classRe       = regexp.MustCompile(`\.((?:[A-Za-z0-9_-]+:)?PREFIX-[A-Za-z0-9_-]+(?:\([^)]*\))?)`)
)

// nestedKey turns the key k of a nested table into its selector or at-rule prelude.
func (l *loader) nestedKey(p []string, k string) string {
	if strings.HasPrefix(k, "@") {
		return l.atRuleKey(p, k)
	}
	if !strings.Contains(k, "&") {
		l.errf(p, "selector must contain &")
		return k
	}
	return l.selector(p, k)
}

// atRuleKey expands a breakpoint name into its media query, or checks and resolves any other at-rule.
func (l *loader) atRuleKey(p []string, k string) string {
	name := k[1:]
	if bp, ok := l.v.bpValue[name]; ok {
		return "@media (min-width: " + bp + ")"
	}
	if below, ok := strings.CutPrefix(name, "below-"); ok {
		if bp, ok := l.v.bpValue[below]; ok {
			return "@media (width < " + bp + ")"
		}
		l.errf(p, "unknown breakpoint %q", below)
		return k
	}
	if oneWordAtRule.MatchString(k) && k != "@starting-style" {
		l.errf(p, "unknown breakpoint %q", name)
		return k
	}
	if pr := rawProblem("", containerName.ReplaceAllString(k, "$1")); pr != "" {
		l.errf(p, "%s in at-rule key", pr)
	}
	return l.literal(p, k)
}

// selector validates the vocabulary classes in a selector and escapes them.
func (l *loader) selector(p []string, key string) string {
	re := regexp.MustCompile(strings.Replace(classRe.String(), "PREFIX", regexp.QuoteMeta(l.v.Project.Prefix), 1))
	return re.ReplaceAllStringFunc(key, func(m string) string {
		c := m[1:]
		if !l.v.valid[c] {
			l.errf(p, "unknown class %s%s", c, l.takesHint(c))
			return m
		}
		return "." + esc(c)
	})
}

// takesHint lists the values an unknown class's family does take.
func (l *loader) takesHint(c string) string {
	fam, _, _ := strings.Cut(c, "(")
	if t, ok := l.v.Takes(fam); ok && len(t) > 0 {
		return "; " + fam + " takes: " + strings.Join(t, " ")
	}
	return ""
}
