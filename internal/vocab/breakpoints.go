package vocab

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
