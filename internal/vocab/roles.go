package vocab

import "strings"

// registerRoleClasses adds every role class to the valid set, so selectors in
// any rule can name any role. Problems are reported here; loadRoles skips
// the same entries without reporting them again.
func (l *loader) registerRoleClasses() {
	v := l.v
	for _, kind := range l.children("role") {
		if l.table("role", kind) == nil {
			l.errf([]string{"role", kind}, "must be a table of roles")
			continue
		}
		if isAtomName(kind) {
			l.errf([]string{"role", kind}, "kind collides with atom %s-%s", v.Project.Prefix, kind)
			continue
		}
		if !bareName.MatchString(kind) {
			l.errf([]string{"role", kind}, "kind must be letters, digits, dashes")
			continue
		}
		fam := v.Project.Prefix + "-" + kind
		v.takes[fam] = nil
		for _, n := range l.children("role", kind) {
			if l.table("role", kind, n) == nil {
				l.errf([]string{"role", kind, n}, "must be a role table")
				continue
			}
			if !bareName.MatchString(n) {
				l.errf([]string{"role", kind, n}, "name must be letters, digits, dashes")
				continue
			}
			v.addClass(fam + "(" + n + ")")
			v.takes[fam] = append(v.takes[fam], n)
		}
	}
}

// loadRoles appends each valid [role.<kind>.<name>] table to l.v.Roles, in file
// order. Invalid entries are skipped here without a second report;
// registerRoleClasses reports them.
func (l *loader) loadRoles() {
	for _, kind := range l.children("role") {
		if l.table("role", kind) == nil || isAtomName(kind) {
			continue
		}
		for _, name := range l.children("role", kind) {
			p := []string{"role", kind, name}
			t := l.table(p...)
			if t == nil || !bareName.MatchString(name) || !bareName.MatchString(kind) {
				continue
			}
			l.v.Roles = append(l.v.Roles, l.buildRole(p, t, kind, name))
		}
	}
}

func (l *loader) buildRole(p []string, t map[string]any, kind, name string) Role {
	class := l.v.Project.Prefix + "-" + kind + "(" + name + ")"
	r := Role{Kind: kind, Name: name, Class: class}
	r.File, r.Line = splitPos(l.pos[pathKey(p)])
	r.Includes = l.includes[pathKey(p)]
	if u, ok := t["use"].(string); !ok || strings.TrimSpace(u) == "" {
		l.errf(extend(p, "use"), "required: when the role applies")
	} else {
		r.Use = u
	}
	r.Rule = l.buildRule(p, t, "."+esc(class), true)
	return r
}
