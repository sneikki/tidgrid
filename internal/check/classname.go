package check

import (
	"fmt"
	"strings"
)

// judgeClass checks one class name against the vocabulary and records its use.
func (f *fileScan) judgeClass(t string, off int) {
	if f.v.Valid(t) {
		f.usedClass(t)
		return
	}
	if f.allowed[t] {
		return
	}
	fam := t
	if i := strings.Index(fam, "("); i >= 0 {
		fam = fam[:i]
	}
	if i := strings.Index(fam, ":"); i >= 0 && !f.isBreakpoint(fam[:i]) {
		f.report(f.line(off), fmt.Sprintf("unknown class %s: %s is not a breakpoint", t, fam[:i]))
		return
	}
	f.report(f.line(off), f.unknownClassMessage(t, fam))
}

func (f *fileScan) isBreakpoint(name string) bool {
	for _, b := range f.v.Breakpoints {
		if b.Name == name {
			return true
		}
	}
	return false
}

// unknownClassMessage says what the class family takes, when it is a known family.
func (f *fileScan) unknownClassMessage(t, fam string) string {
	vals, ok := f.v.Takes(fam)
	switch {
	case !ok:
		return fmt.Sprintf("unknown class %s: not in the vocabulary; run `tidgrid list`", t)
	case len(vals) == 0:
		return fmt.Sprintf("unknown class %s; %s takes no value", t, fam)
	}
	return fmt.Sprintf("unknown class %s; %s takes: %s", t, fam, strings.Join(vals, " "))
}
