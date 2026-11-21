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
		e.Rule = l.buildRule(p, t, sel, true)
		l.v.elements = append(l.v.elements, e)
	}
}
