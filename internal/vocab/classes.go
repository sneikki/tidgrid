package vocab

// registerAtomClasses expands the atom families over the loaded tokens and
// breakpoints, and records each atom class (plain and per breakpoint) as
// valid, with the values its family takes.
func (l *loader) registerAtomClasses() {
	v := l.v
	names := map[string][]string{} // token group -> token names
	for _, g := range v.Tokens {
		for _, t := range g.Tokens {
			names[g.Name] = append(names[g.Name], t.Name)
		}
	}
	v.Atoms = atoms(v.Project.Prefix, names, v.Breakpoints)
	for _, a := range v.Atoms {
		v.addClass(a.Class)
		fam := v.Project.Prefix + "-" + a.Family
		v.takes[fam] = v.takes[fam] // present even when empty
		if a.Value != "" {
			v.takes[fam] = append(v.takes[fam], a.Value)
		}
	}
	for _, b := range v.Breakpoints {
		for _, a := range v.Atoms {
			if a.Below == "" {
				v.addClass(b.Name + ":" + a.Class)
			}
		}
	}
}
