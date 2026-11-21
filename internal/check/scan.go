package check

import (
	"regexp"
	"sort"
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// fileScan reads one masked source file and reports through callbacks, so it
// knows nothing of paths or of the run it belongs to.
type fileScan struct {
	v         *vocab.Vocabulary
	src       string
	nl        []int // offsets of the newlines, for line numbers
	seen      map[int]bool
	allowed   map[string]bool
	report    func(line int, msg string)
	usedClass func(class string)
}

func newFileScan(v *vocab.Vocabulary, src string, allowed map[string]bool) *fileScan {
	f := &fileScan{v: v, src: src, seen: map[int]bool{}, allowed: allowed}
	for i := 0; i < len(src); i++ {
		if src[i] == '\n' {
			f.nl = append(f.nl, i)
		}
	}
	return f
}

func (f *fileScan) line(off int) int { return 1 + sort.SearchInts(f.nl, off) }

// scan runs every reader over the file. Order matters only for seen: the
// first reader to reach a value reads it.
func (f *fileScan) scan(helpers []string) {
	f.scanClassAttrs()
	f.scanClassAssignments()
	f.scanClassDirectives()
	f.scanStyleDirectives()
	f.scanStyleAttrs()
	f.scanHelperCalls(helpers)
	f.scanStyleTags()
}

// attrForm is how an attribute value is written.
type attrForm int

const (
	formQuoted attrForm = iota
	formBraced
	formBare
)

// attrSpan finds the value that starts at p. ok is false when a quote or brace
// never closes.
func (f *fileScan) attrSpan(p int) (start, end int, form attrForm, ok bool) {
	s := f.src
	switch s[p] {
	case '"', '\'':
		e := attrValueEnd(s, p)
		return p + 1, e, formQuoted, e > 0
	case '{':
		e := delimEnd(s, p, len(s), '{', '}')
		return p + 1, e - 1, formBraced, e > 0
	}
	return p, attrEnd(s, p), formBare, true
}

func (f *fileScan) scanClassAttrs() {
	s := f.src
	for _, m := range classAttr.FindAllStringSubmatchIndex(s, -1) {
		name, p := strings.ToLower(s[m[2]:m[3]]), m[1]
		if p >= len(s) {
			continue
		}
		if name == "classname" && declRe.MatchString(s[:m[2]]) { // const className = …: script, not an attribute
			if s[p] != '=' {
				f.readLiterals(p, stmtEnd(s, p))
			}
			continue
		}
		dynamic := strings.HasPrefix(name, ":") || strings.HasPrefix(name, "v-bind")
		start, end, form, ok := f.attrSpan(p)
		switch {
		case !ok:
		case form == formBraced, form == formQuoted && dynamic:
			f.readLiterals(start, end)
		default:
			f.readClassValue(start, end, modeAttr)
		}
	}
}

func (f *fileScan) scanClassAssignments() { // el.className = …
	s := f.src
	for _, m := range assignRe.FindAllStringIndex(s, -1) {
		if m[1] < len(s) && s[m[1]] != '=' {
			f.readLiterals(m[1], stmtEnd(s, m[1]))
		}
	}
}

func (f *fileScan) scanClassDirectives() { // Svelte class:name
	s := f.src
	for _, m := range classDir.FindAllStringSubmatchIndex(s, -1) {
		if n := s[m[2]:m[3]]; n != "list" {
			f.judgeClass(n, m[2])
		}
	}
}

func (f *fileScan) scanStyleDirectives() { // Svelte style:prop
	s := f.src
	for _, m := range styleDir.FindAllStringSubmatchIndex(s, -1) {
		f.judgeStyleProp(s[m[2]:m[3]], m[2])
	}
}

func (f *fileScan) scanStyleAttrs() {
	for _, m := range styleAttr.FindAllStringSubmatchIndex(f.src, -1) {
		p := m[1]
		if p >= len(f.src) {
			continue
		}
		if start, end, _, ok := f.attrSpan(p); ok {
			f.judgeStyleValue(start, end)
		}
	}
}

// scanHelperCalls reads the arguments of class helpers (cn, clsx, …) and of
// classList.add and its siblings.
func (f *fileScan) scanHelperCalls(helpers []string) {
	s := f.src
	names := make([]string, 0, len(helpers)+1)
	for _, h := range helpers {
		names = append(names, regexp.QuoteMeta(h))
	}
	names = append(names, `classList\s*\.\s*(?:add|remove|toggle|contains)`)
	call := regexp.MustCompile(`(?:^|[^\w$])(?:` + strings.Join(names, "|") + `)\s*\(`)
	for _, m := range call.FindAllStringIndex(s, -1) {
		p := m[1] - 1
		if e := delimEnd(s, p, len(s), '(', ')'); e > 0 {
			f.readLiterals(p+1, e-1)
		}
	}
}

func (f *fileScan) scanStyleTags() {
	for _, m := range styleTag.FindAllStringIndex(f.src, -1) {
		f.report(f.line(m[0]), "own CSS: <style> block; make it a role in the vocabulary")
	}
}
