package vocab

import (
	"path/filepath"
	"regexp"
)

var prefixRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// loadProject reads [project]; relative paths are resolved against dir, the vocabulary's directory.
func (l *loader) loadProject(dir string) {
	p := Project{Prefix: "tg", Scan: []string{"."}, Helpers: []string{"cn", "clsx", "cva", "classNames", "twMerge"}}
	t := l.table("project")
	if s, ok := t["prefix"].(string); ok {
		p.Prefix = s
	}
	if !prefixRe.MatchString(p.Prefix) {
		l.errf([]string{"project", "prefix"}, "must be lowercase letters, digits, dashes")
	}
	if s, ok := t["out"].(string); ok && s != "" {
		p.Out = filepath.Join(dir, s)
	} else {
		l.errf([]string{"project", "out"}, "required")
	}
	p.Scan = l.stringList(t, "scan", p.Scan)
	for i, s := range p.Scan {
		// A missing scan path is check's violation, not a load error: build,
		// list and sheet must work in a project that has no sources yet.
		p.Scan[i] = filepath.Join(dir, s)
	}
	p.Helpers = l.stringList(t, "helpers", p.Helpers)
	l.v.Project = p
}

// stringList reads [project] key as a list of strings, def when it is absent.
func (l *loader) stringList(project map[string]any, key string, def []string) []string {
	items, ok := project[key].([]any)
	if !ok {
		if _, set := project[key]; set {
			l.errf([]string{"project", key}, "must be a list of strings")
		}
		return def
	}
	var out []string
	for _, x := range items {
		if s, ok := x.(string); ok {
			out = append(out, s)
		} else {
			l.errf([]string{"project", key}, "must be a list of strings")
		}
	}
	return out
}

func (l *loader) loadBreakpoints() {
	for _, k := range l.children("breakpoints") {
		s, ok := l.raw["breakpoints"].(map[string]any)[k].(string)
		if !ok {
			l.errf([]string{"breakpoints", k}, "value must be a string")
			continue
		}
		l.v.Breakpoints = append(l.v.Breakpoints, Breakpoint{k, s})
		l.v.bpValue[k] = s
	}
}
