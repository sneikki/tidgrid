package vocab

func (l *loader) loadElements() {
	for _, sel := range l.children("element") {
		p := []string{"element", sel}
		t := l.table(p...)
		if t == nil {
			l.errf(p, "must be a table")
			continue
		}
		e := element{Selector: sel}
		e.Use, _ = t["use"].(string)
		e.Rule = l.buildRule(p, t, sel, sel, true)
		l.v.elements = append(l.v.elements, e)
	}
}

func (l *loader) loadKeyframes() {
	for _, name := range l.children("keyframes") {
		p := []string{"keyframes", name}
		t := l.table(p...)
		if t == nil {
			l.errf(p, "must be a table of stops")
			continue
		}
		kf := keyframe{Name: name}
		for _, stop := range l.children(p...) {
			sp := extend(p, stop)
			st := l.table(sp...)
			if st == nil {
				l.errf(sp, "must be a table")
				continue
			}
			kf.Stops = append(kf.Stops, l.buildRule(sp, st, stop, stop, false))
		}
		l.v.keyframes = append(l.v.keyframes, kf)
	}
}
