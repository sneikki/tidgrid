package vocab

import (
	"strings"
)

// Convert reads the TOML vocabulary at path and returns the Starlark files that
// declare the same vocabulary: index.tg, tokens.tg and roles.tg. Values stay
// strings as written; TOML comments are lost. The TOML is validated as
// Load does first; its errors are returned and nothing is converted.
func Convert(path string) (map[string]string, []error) {
	l := newLoader(path)
	if err := l.decodeTOML(path); err != nil {
		return nil, []error{err}
	}
	if _, errs := l.validate(); len(errs) > 0 {
		return nil, errs
	}
	return map[string]string{
		"index.tg":  l.convertIndex(),
		"tokens.tg": l.convertTokens(),
		"roles.tg":  l.convertRoles(),
	}, nil
}

func (l *loader) convertIndex() string {
	var b strings.Builder
	if t := l.table("project"); t != nil {
		b.WriteString("# Repeated lines are not checked until repeats is removed; then build says which roles share lines.\n")
		b.WriteString("project(" + l.args([]string{"project"}, t, "") + ", repeats = [])\n\n")
	}
	if t := l.table("breakpoints"); t != nil {
		b.WriteString("breakpoints(" + l.dict([]string{"breakpoints"}, t, 0) + ")\n\n")
	}
	for _, sel := range l.children("element") {
		p := []string{"element", sel}
		if t := l.table(p...); t != nil {
			b.WriteString(l.call("element", p, t, quote(sel)))
		}
	}
	for _, kind := range []string{"class", "css"} {
		for _, k := range l.children("allow", kind) {
			s, _ := scalar(l.table("allow", kind)[k])
			b.WriteString("allow_" + kind + "(" + quote(k) + ", " + quote(s) + ")\n")
		}
	}
	return b.String() + "\nsource(\"tokens.tg\")\nsource(\"roles.tg\")\n"
}

func (l *loader) convertTokens() string {
	var b strings.Builder
	taken := map[string]bool{}
	for _, g := range l.children("tokens") {
		if t := l.table("tokens", g); t != nil {
			b.WriteString(safeName(g, taken) + " = tokens(" + quote(g) + ", " + l.dict([]string{"tokens", g}, t, 0) + ")\n\n")
		}
	}
	return b.String()
}

func (l *loader) convertRoles() string {
	var b strings.Builder
	for _, name := range l.children("keyframes") {
		p := []string{"keyframes", name}
		if t := l.table(p...); t != nil {
			b.WriteString("keyframes(" + quote(name) + ", " + l.dict(p, t, 0) + ")\n\n")
		}
	}
	for _, kind := range l.children("role") {
		for _, name := range l.children("role", kind) {
			p := []string{"role", kind, name}
			if t := l.table(p...); t != nil {
				b.WriteString(l.call("role", p, t, quote(kind)+", "+quote(name)))
			}
		}
	}
	return b.String()
}
