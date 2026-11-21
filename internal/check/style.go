package check

import (
	"fmt"
	"strings"
)

// judgeStyleProp flags one property name (CSS or camelCase, any case) if it is appearance.
func (f *fileScan) judgeStyleProp(name string, off int) {
	p := name
	if strings.ToUpper(p) == p || strings.Contains(p, "-") {
		p = strings.ToLower(p)
	} else {
		p = strings.ToLower(p[:1]) + capRe.ReplaceAllStringFunc(p[1:], func(c string) string { return "-" + strings.ToLower(c) })
	}
	if appearRe.MatchString(p) {
		f.report(f.line(off), fmt.Sprintf("appearance in style attribute: %s; use a role (`tidgrid list`)", p))
	}
}

// judgeStyleValue flags appearance properties in s[start:end], a style attribute value.
func (f *fileScan) judgeStyleValue(start, end int) {
	text := f.src[start:end]
	for _, m := range propRe.FindAllStringSubmatchIndex(text, -1) {
		f.judgeStyleProp(text[m[2]:m[3]], start+m[2])
	}
	if strings.HasPrefix(strings.TrimSpace(text), "{") { // JSX object: { background } is shorthand
		for at := 0; at < len(text); {
			m := shortRe.FindStringSubmatchIndex(text[at:])
			if m == nil {
				break
			}
			f.judgeStyleProp(text[at+m[2]:at+m[3]], start+at+m[2])
			at += m[3] // the comma after the name starts the next match
		}
	}
}
