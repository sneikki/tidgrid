package vocab

import "strings"

// checkMixins runs the rules that need every role and mixin: a set needs two
// users, and a role may not write out a whole set by hand.
func (l *loader) checkMixins() {
	for _, r := range l.mix.roles {
		l.checkHandWritten(r)
	}
	for _, m := range l.mix.defs {
		for _, s := range m.sets {
			l.checkUsers(m, s)
		}
	}
}

func (l *loader) checkUsers(m *mixinDef, s *setDef) {
	p := []string{"mixin", m.name, s.name}
	switch len(s.users) {
	case 0:
		l.errf(p, "included by no role; a set needs at least two")
	case 1:
		u := s.users[0]
		l.errf(p, "included only by role %s.%s; a set needs at least two", u.kind, u.name)
	}
}

func (l *loader) checkHandWritten(r roleUse) {
	for _, m := range l.mix.defs {
		for _, s := range m.sets {
			if !contains(r.incl, s) && writesOut(r, s) {
				l.errf(r.p, "own lines repeat set %s; write include = [%s] and drop them", s.label(), s.label())
			}
		}
	}
}

// writesOut is true when every top-level line of s, same property and value as written, is in r.
func writesOut(r roleUse, s *setDef) bool {
	for _, ln := range s.tops {
		if !r.own[strings.Join(ln[:], "\x00")] {
			return false
		}
	}
	return len(s.tops) > 0
}

func contains(list []*setDef, s *setDef) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
