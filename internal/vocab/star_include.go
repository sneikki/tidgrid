package vocab

import (
	"go.starlark.net/starlark"
)

// declareRole declares the role table: the included sets' lines first, in
// include order, then the role's own css. A property given twice is an error.
func (l *loader) declareRole(p []string, use string, css *starlark.Dict, include starlark.Value, at string) error {
	t := map[string]any{}
	if err := l.declare(p, t, at); err != nil {
		return err
	}
	if use != "" {
		l.putKey(t, p, "use", use)
	}
	owner := map[string]string{}
	sets := l.includedSets(p, include)
	for _, s := range sets {
		l.addLines(p, t, s.tree, s.label(), owner)
		l.noteUse(p, s)
	}
	ru := roleUse{p: p, incl: sets, own: map[string]bool{}}
	if css != nil {
		own := l.snapshot(p, css)
		l.addLines(p, t, own, "", owner)
		for _, ln := range topLines(own) {
			ru.own[ln[0]+"\x00"+ln[1]] = true
		}
	}
	l.mix.roles = append(l.mix.roles, ru)
	return nil
}

// noteUse records that the role at p includes s.
func (l *loader) noteUse(p []string, s *setDef) {
	l.includes[pathKey(p)] = append(l.includes[pathKey(p)], s.label())
	u := roleRef{p[1], p[2]}
	if n := len(s.users); n == 0 || s.users[n-1] != u {
		s.users = append(s.users, u)
	}
}

// includedSets reads the include argument; only mixin sets are accepted.
func (l *loader) includedSets(p []string, include starlark.Value) []*setDef {
	if include == nil {
		return nil
	}
	ip := extend(p, "include")
	seq, ok := include.(starlark.Sequence)
	if !ok {
		l.errf(ip, "must be a list of mixin sets, such as [type.heading]")
		return nil
	}
	var out []*setDef
	for it, x := seq.Iterate(), starlark.Value(nil); it.Next(&x); {
		if s, ok := x.(*setDef); ok {
			if contains(out, s) {
				l.errf(ip, "set %s is included twice", s.label())
				continue
			}
			out = append(out, s)
		} else {
			l.errf(ip, "%s is not a mixin set; write a set such as type.heading", x.String())
		}
	}
	return out
}

// addLines copies lines into the table t. from names the set the lines
// belong to, empty for the role's own css; owner says who gave each
// property, so a second giver is an error at any depth. A nested rule that
// is already there is merged into: the first giver's lines stay first.
func (l *loader) addLines(p []string, t map[string]any, ls []cssLine, from string, owner map[string]string) {
	for _, ln := range ls {
		kp := extend(p, ln.key)
		k := pathKey(kp)
		if sub, ok := t[ln.key].(map[string]any); ok && ln.nested {
			l.addLines(kp, sub, ln.sub, from, owner)
			continue
		}
		if prev, dup := owner[k]; dup {
			l.errf(kp, "given twice: by %s and by %s", giver(prev), giver(from))
			continue
		}
		owner[k] = from
		if from != "" {
			l.origin[k] = from
		}
		if ln.nested {
			sub := map[string]any{}
			l.putKey(t, p, ln.key, sub)
			l.addLines(kp, sub, ln.sub, from, owner)
			continue
		}
		v, err := l.cell(kp, ln.val)
		if err != nil {
			l.errf(kp, "%s", err.Error())
			continue
		}
		l.putKey(t, p, ln.key, v)
	}
}

func giver(set string) string {
	if set == "" {
		return "the role's own css"
	}
	return "set " + set
}
