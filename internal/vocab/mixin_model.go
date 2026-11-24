package vocab

// Mixin is a named group of css sets; a role includes sets by name.
type Mixin struct {
	Name, Use string
	Sets      []MixinSet
}

// MixinSet is one set: its lines as written (nested rules as `sel { a: b }`)
// and the classes of the roles that include it.
type MixinSet struct {
	Name  string
	Lines []string
	Roles []string
}

// exportMixins fills Vocabulary.Mixins once the roles are loaded.
func (l *loader) exportMixins() {
	for _, m := range l.mix.defs {
		out := Mixin{Name: m.name, Use: m.use}
		for _, s := range m.sets {
			ms := MixinSet{Name: s.name, Lines: s.lines}
			for _, u := range s.users {
				ms.Roles = append(ms.Roles, l.v.Project.Prefix+"-"+u.kind+"("+u.name+")")
			}
			out.Sets = append(out.Sets, ms)
		}
		l.v.Mixins = append(l.v.Mixins, out)
	}
}
