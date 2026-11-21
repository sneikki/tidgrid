package vocab

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
			kf.Stops = append(kf.Stops, l.buildRule(sp, st, stop, false))
		}
		l.v.keyframes = append(l.v.keyframes, kf)
	}
}
