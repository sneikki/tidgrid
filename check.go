package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	srcExts   = map[string]bool{".html": true, ".htm": true, ".astro": true, ".jsx": true, ".tsx": true, ".vue": true, ".svelte": true, ".md": true, ".mdx": true, ".liquid": true, ".njk": true, ".js": true, ".mjs": true, ".cjs": true, ".ts": true}
	cssExts   = map[string]bool{".css": true, ".scss": true, ".sass": true, ".less": true, ".styl": true}
	skipDirs  = map[string]bool{"node_modules": true, ".git": true, "dist": true, "build": true, "_site": true, ".astro": true, "vendor": true}
	classAttr = regexp.MustCompile(`(?i)(?:^|[^\w.:-])((?:v-bind:|:)?class(?::list)?|className)\s*=\s*`)
	assignRe  = regexp.MustCompile(`\.\s*className\s*\+?=`)
	declRe    = regexp.MustCompile(`\b(?:const|let|var)\s+$`)
	classDir  = regexp.MustCompile(`(?i)(?:^|\s)class:([^\s=>/{}"'|]+)`)
	styleAttr = regexp.MustCompile(`(?i)(?:^|[^\w.:-])((?:v-bind:|:)?style)\s*=\s*`)
	styleDir  = regexp.MustCompile(`(?i)(?:^|\s)style:([A-Za-z-]+)`)
	shortRe   = regexp.MustCompile(`[{,]\s*([A-Za-z]\w*)\s*(?:[,}]|$)`)
	styleTag  = regexp.MustCompile(`(?i)<style(?:[\s>/]|$)`)
	propRe    = regexp.MustCompile(`(?:^|[;{,\s'"])(-?[A-Za-z][A-Za-z0-9-]*)['"]?\s*:`)
	capRe     = regexp.MustCompile(`[A-Z]`)
	appearRe  = regexp.MustCompile(`^(?:(?:background|border|font|outline|text-decoration|mask|column-rule)(?:-.*)?|color|box-shadow|text-shadow|opacity|filter|backdrop-filter|letter-spacing|line-height|fill|stroke|accent-color|caret-color|mix-blend-mode)$`)
	rawBlock  = regexp.MustCompile(`(?s)\{%-?\s*(highlight|raw)\b.*?\{%-?\s*end(highlight|raw)\s*-?%\}`)
	ignoreRe  = regexp.MustCompile(`(?s)tidgrid-ignore-start(?:.*?tidgrid-ignore-end|.*$)`)
)

type violation struct {
	path string
	line int
	msg  string
}

// mask blanks out what is not read (code samples), keeping offsets and newlines.
// open lists the offsets of a tidgrid-ignore-start that is never closed.
func mask(src, ext string) (out string, open []int) {
	b := []byte(src)
	blank := func(s, e int) {
		for i := s; i < e; i++ {
			if b[i] != '\n' {
				b[i] = ' '
			}
		}
	}
	for _, re := range []*regexp.Regexp{rawBlock, ignoreRe} {
		for _, m := range re.FindAllIndex(b, -1) {
			if re == ignoreRe && !strings.HasSuffix(string(b[m[0]:m[1]]), "tidgrid-ignore-end") {
				open = append(open, m[0])
			}
			blank(m[0], m[1])
		}
	}
	if ext == ".md" || ext == ".mdx" {
		in, off := false, 0
		for _, ln := range strings.SplitAfter(string(b), "\n") {
			t := strings.TrimLeft(ln, " \t")
			fence := strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")
			if fence || in {
				blank(off, off+len(ln))
			}
			if fence {
				in = !in
			}
			off += len(ln)
		}
	}
	return string(b), open
}

type fileScan struct {
	v       *Vocabulary
	src     string
	nl      []int
	seen    map[int]bool
	allowed map[string]bool
	add     func(line int, msg string)
	use     func(class string)
}

func (f *fileScan) line(off int) int { return 1 + sort.SearchInts(f.nl, off) }

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' }

// skipString returns the index after the string opening at i and the end of
// its content.
func skipString(s string, i, end int) (int, int) {
	q := s[i]
	for k := i + 1; k < end; {
		switch {
		case s[k] == '\\':
			k += 2
		case q == '`' && strings.HasPrefix(s[k:end], "${"):
			if e := delimEnd(s, k+1, end, '{', '}'); e > 0 {
				k = e
			} else {
				k += 2
			}
		case s[k] == q:
			return k + 1, k
		default:
			k++
		}
	}
	return end, end
}

// delimEnd returns the index after the delimiter matching s[i]==open, or -1.
func delimEnd(s string, i, end int, open, close byte) int {
	depth := 0
	for k := i; k < end; {
		switch c := s[k]; {
		case c == '\'' || c == '"' || c == '`':
			k, _ = skipString(s, k, end)
		case c == open:
			depth++
			k++
		case c == close:
			depth--
			k++
			if depth == 0 {
				return k
			}
		default:
			k++
		}
	}
	return -1
}

// attrValueEnd returns the index of the quote closing the value that opens at i,
// skipping template tags and interpolations that may contain the quote.
func attrValueEnd(s string, i int) int {
	q, end := s[i], len(s)
	for k := i + 1; k < end; {
		switch {
		case strings.HasPrefix(s[k:], "{{"), strings.HasPrefix(s[k:], "{%"):
			cl := "}}"
			if s[k+1] == '%' {
				cl = "%}"
			}
			if j := strings.Index(s[k+2:], cl); j >= 0 {
				k += j + 4
			} else {
				k += 2
			}
		case s[k] == '{', strings.HasPrefix(s[k:], "${"):
			o := k
			if s[k] == '$' {
				o++
			}
			if e := delimEnd(s, o, end, '{', '}'); e > 0 {
				k = e
			} else {
				k++
			}
		case s[k] == q:
			return k
		default:
			k++
		}
	}
	return -1
}

const (
	modeAttr = iota // HTML-ish attribute text: {% %}, {{ }}, ${ }, { }
	modeTmpl        // JS template literal: ${ }
	modeJS          // plain JS string
)

// value reads the class names in s[start:end].
func (f *fileScan) value(start, end, mode int) {
	if f.seen[start] {
		return
	}
	f.seen[start] = true
	s := f.src
	var tok strings.Builder
	tokStart, built := -1, false
	var inner [][2]int
	begin := func(i int) {
		if tokStart < 0 {
			tokStart = i
		}
	}
	flush := func(at int) {
		if tokStart < 0 {
			return
		}
		switch {
		case !built:
			f.class(tok.String(), tokStart)
		case tok.String() == "\x00": // interpolation standing alone: read its string literals
			for _, r := range inner {
				f.collect(r[0], r[1])
			}
		default:
			text := strings.Join(strings.Fields(s[tokStart:at]), " ")
			if r := []rune(text); len(r) > 60 {
				text = string(r[:60]) + "…"
			}
			f.add(f.line(tokStart), fmt.Sprintf("class name built from pieces: %s; write every whole name out in a lookup table", text))
		}
		tok.Reset()
		tokStart, built, inner = -1, false, nil
	}
	interp := func(i, o, e int, lits bool) int { // interpolation s[i:e] with its brace at o
		begin(i)
		built = true
		tok.WriteByte(0)
		if lits {
			inner = append(inner, [2]int{o + 1, e - 1})
		}
		return e
	}
	for i := start; i < end; {
		c := s[i]
		switch {
		case isSpace(c):
			flush(i)
			i++
		case mode == modeAttr && strings.HasPrefix(s[i:end], "{%"):
			flush(i)
			if j := strings.Index(s[i+2:end], "%}"); j >= 0 {
				i += j + 4
			} else {
				i = end
			}
		case mode == modeAttr && strings.HasPrefix(s[i:end], "{{"):
			if j := strings.Index(s[i+2:end], "}}"); j >= 0 {
				i = interp(i, i, i+j+4, false)
			} else {
				begin(i)
				tok.WriteByte(c)
				i++
			}
		case mode != modeJS && strings.HasPrefix(s[i:end], "${"), mode == modeAttr && c == '{':
			o := i
			if c == '$' {
				o++
			}
			if e := delimEnd(s, o, end, '{', '}'); e > 0 {
				i = interp(i, o, e, true)
			} else {
				begin(i)
				tok.WriteByte(c)
				i++
			}
		default:
			begin(i)
			tok.WriteByte(c)
			i++
		}
	}
	flush(end)
}

// compared reports whether the literal s[i:j] is an operand of ==, ===, != or !==.
func (f *fileScan) compared(i, j int) bool {
	s := f.src
	b := strings.TrimRight(s[:i], " \t\r\n")
	a := strings.TrimLeft(s[j:], " \t\r\n")
	return strings.HasSuffix(b, "==") || strings.HasSuffix(b, "!=") || strings.HasPrefix(a, "==") || strings.HasPrefix(a, "!=")
}

// stmtEnd returns the end of the statement that starts at i: a `;` or newline
// outside brackets, or the bracket that closes the enclosing one.
func stmtEnd(s string, i int) int {
	depth := 0
	for k := i; k < len(s); {
		switch c := s[k]; {
		case c == '\'' || c == '"' || c == '`':
			k, _ = skipString(s, k, len(s))
		case c == '(' || c == '[' || c == '{':
			depth++
			k++
		case c == ')' || c == ']' || c == '}':
			if depth--; depth < 0 {
				return k
			}
			k++
		case depth == 0 && (c == ';' || c == '\n'):
			return k
		default:
			k++
		}
	}
	return len(s)
}

// attrEnd returns the end of an unquoted attribute value starting at p.
func attrEnd(s string, p int) int {
	e := p
	for e < len(s) && !isSpace(s[e]) && s[e] != '>' && !strings.HasPrefix(s[e:], "/>") {
		e++
	}
	return e
}

// collect reads every string literal in s[i:end].
func (f *fileScan) collect(i, end int) {
	s := f.src
	for i < end {
		switch c := s[i]; {
		case c == '\'' || c == '"' || c == '`':
			j, ie := skipString(s, i, end)
			mode := modeJS
			if c == '`' {
				mode = modeTmpl
			}
			if !f.compared(i, j) {
				f.value(i+1, ie, mode)
			}
			i = j
		case c == '/' && i+1 < end && s[i+1] == '/':
			for i < end && s[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < end && s[i+1] == '*':
			if j := strings.Index(s[i+2:end], "*/"); j >= 0 {
				i += j + 4
			} else {
				i = end
			}
		default:
			i++
		}
	}
}

// class checks one name and records its use.
func (f *fileScan) class(t string, off int) {
	if f.v.Valid(t) {
		f.use(t)
		return
	}
	if f.allowed[t] {
		return
	}
	fam := t
	if i := strings.Index(fam, "("); i >= 0 {
		fam = fam[:i]
	}
	msg := fmt.Sprintf("unknown class %s: not in the vocabulary; run `tidgrid list`", t)
	if i := strings.Index(fam, ":"); i >= 0 {
		bp := false
		for _, b := range f.v.Breakpoints {
			bp = bp || b.Name == fam[:i]
		}
		if !bp {
			f.add(f.line(off), fmt.Sprintf("unknown class %s: %s is not a breakpoint", t, fam[:i]))
			return
		}
	}
	if vals, ok := f.v.Takes(fam); ok {
		if len(vals) == 0 {
			msg = fmt.Sprintf("unknown class %s; %s takes no value", t, fam)
		} else {
			msg = fmt.Sprintf("unknown class %s; %s takes: %s", t, fam, strings.Join(vals, " "))
		}
	}
	f.add(f.line(off), msg)
}

// styleProp flags one property name (CSS or camelCase, any case) if it is appearance.
func (f *fileScan) styleProp(name string, off int) {
	p := name
	if strings.ToUpper(p) == p || strings.Contains(p, "-") {
		p = strings.ToLower(p)
	} else {
		p = strings.ToLower(p[:1]) + capRe.ReplaceAllStringFunc(p[1:], func(c string) string { return "-" + strings.ToLower(c) })
	}
	if appearRe.MatchString(p) {
		f.add(f.line(off), fmt.Sprintf("appearance in style attribute: %s; use a role (`tidgrid list`)", p))
	}
}

// style flags appearance properties in s[start:end], a style attribute value.
func (f *fileScan) style(start, end int) {
	text := f.src[start:end]
	flag := f.styleProp
	for _, m := range propRe.FindAllStringSubmatchIndex(text, -1) {
		flag(text[m[2]:m[3]], start+m[2])
	}
	if strings.HasPrefix(strings.TrimSpace(text), "{") { // JSX object: { background } is shorthand
		for at := 0; at < len(text); {
			m := shortRe.FindStringSubmatchIndex(text[at:])
			if m == nil {
				break
			}
			flag(text[at+m[2]:at+m[3]], start+at+m[2])
			at += m[3] // the comma after the name starts the next match
		}
	}
}

func (f *fileScan) scan(helpers []string) {
	s := f.src
	for _, m := range classAttr.FindAllStringSubmatchIndex(s, -1) {
		name, p := strings.ToLower(s[m[2]:m[3]]), m[1]
		dynamic := strings.HasPrefix(name, ":") || strings.HasPrefix(name, "v-bind")
		if p >= len(s) {
			continue
		}
		if name == "classname" && declRe.MatchString(s[:m[2]]) { // const className = …: script, not an attribute
			if p < len(s) && s[p] != '=' {
				f.collect(p, stmtEnd(s, p))
			}
			continue
		}
		switch s[p] {
		case '"', '\'':
			if e := attrValueEnd(s, p); e > 0 {
				if dynamic {
					f.collect(p+1, e)
				} else {
					f.value(p+1, e, modeAttr)
				}
			}
		case '{':
			if e := delimEnd(s, p, len(s), '{', '}'); e > 0 {
				f.collect(p+1, e-1)
			}
		default:
			f.value(p, attrEnd(s, p), modeAttr)
		}
	}
	for _, m := range assignRe.FindAllStringIndex(s, -1) { // el.className = …
		if m[1] < len(s) && s[m[1]] != '=' {
			f.collect(m[1], stmtEnd(s, m[1]))
		}
	}
	for _, m := range classDir.FindAllStringSubmatchIndex(s, -1) { // Svelte class:name
		if n := s[m[2]:m[3]]; n != "list" {
			f.class(n, m[2])
		}
	}
	for _, m := range styleDir.FindAllStringSubmatchIndex(s, -1) { // Svelte style:prop
		f.styleProp(s[m[2]:m[3]], m[2])
	}
	for _, m := range styleAttr.FindAllStringSubmatchIndex(s, -1) {
		p := m[1]
		if p >= len(s) {
			continue
		}
		switch s[p] {
		case '"', '\'':
			if e := attrValueEnd(s, p); e > 0 {
				f.style(p+1, e)
			}
		case '{':
			if e := delimEnd(s, p, len(s), '{', '}'); e > 0 {
				f.style(p+1, e-1)
			}
		default:
			f.style(p, attrEnd(s, p))
		}
	}
	hs := make([]string, 0, len(helpers)+1)
	for _, h := range helpers {
		hs = append(hs, regexp.QuoteMeta(h))
	}
	hs = append(hs, `classList\s*\.\s*(?:add|remove|toggle|contains)`)
	call := regexp.MustCompile(`(?:^|[^\w$])(?:` + strings.Join(hs, "|") + `)\s*\(`)
	for _, m := range call.FindAllStringIndex(s, -1) {
		p := m[1] - 1
		if e := delimEnd(s, p, len(s), '(', ')'); e > 0 {
			f.collect(p+1, e-1)
		}
	}
	for _, m := range styleTag.FindAllStringIndex(s, -1) {
		f.add(f.line(m[0]), "own CSS: <style> block; make it a role in the vocabulary")
	}
}

func runCheck(args []string) int {
	v, code := loadFlag("check", args)
	if v == nil {
		return code
	}
	return doCheck(v, os.Stdout)
}

func doCheck(v *Vocabulary, w io.Writer) int {
	dir := filepath.Dir(v.Path)
	rel := func(p string) string {
		if r, err := filepath.Rel(dir, p); err == nil {
			p = r
		}
		return filepath.ToSlash(p)
	}
	abs := func(p string) string {
		a, err := filepath.Abs(p)
		if err != nil {
			return p
		}
		return a
	}
	dir = abs(dir)
	var vios []violation
	seenV := map[violation]bool{}
	addV := func(p string, line int, msg string) {
		x := violation{p, line, msg}
		if !seenV[x] {
			seenV[x] = true
			vios = append(vios, x)
		}
	}

	allowedClass := map[string]bool{}
	for _, a := range v.AllowClass {
		allowedClass[a.Name] = true
	}
	allowedCSS := map[string]bool{abs(v.Project.Out): true}
	for _, a := range v.AllowCSS {
		allowedCSS[abs(filepath.Join(dir, a.Name))] = true
	}
	roleClass := map[string]bool{}
	for _, r := range v.Roles {
		roleClass[r.Class] = true
	}

	files := map[string]bool{}
	for _, root := range v.Project.Scan {
		filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil { // an unreadable path must not pass as clean
				msg := err.Error()
				if u := errors.Unwrap(err); u != nil {
					msg = u.Error()
				}
				addV(rel(abs(p)), 1, "cannot scan: "+msg)
				return nil
			}
			if d.IsDir() {
				if p != root && skipDirs[d.Name()] {
					return filepath.SkipDir
				}
				return nil
			}
			files[abs(p)] = true
			return nil
		})
	}
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	usedBy := map[string]map[string]bool{}
	usedAny := map[string]bool{}
	for _, p := range paths {
		ext := strings.ToLower(filepath.Ext(p))
		r := rel(p)
		if cssExts[ext] {
			if !allowedCSS[p] {
				addV(r, 1, "own CSS: stylesheet not under allow.css; make it roles in the vocabulary")
			}
			continue
		}
		if !srcExts[ext] {
			continue
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			addV(r, 1, "cannot read file: "+errors.Unwrap(err).Error())
			continue
		}
		src, open := mask(string(raw), ext)
		f := &fileScan{v: v, src: src, seen: map[int]bool{}, allowed: allowedClass}
		for i := 0; i < len(src); i++ {
			if src[i] == '\n' {
				f.nl = append(f.nl, i)
			}
		}
		f.add = func(line int, msg string) { addV(r, line, msg) }
		f.use = func(c string) {
			if roleClass[c] {
				if usedBy[r] == nil {
					usedBy[r] = map[string]bool{}
				}
				usedBy[r][c] = true
				usedAny[c] = true
			}
		}
		for _, o := range open {
			f.add(f.line(o), "tidgrid-ignore-start is never closed; the rest of the file is not checked")
		}
		f.scan(v.Project.Helpers)
	}

	out := v.Project.Out
	if b, err := os.ReadFile(out); err != nil {
		addV(rel(abs(out)), 1, "stylesheet missing; run `tidgrid build`")
	} else if string(b) != v.render() {
		addV(rel(abs(out)), 1, "stylesheet differs from the vocabulary; run `tidgrid build`")
	}

	sort.Slice(vios, func(i, j int) bool {
		a, b := vios[i], vios[j]
		if a.path != b.path {
			return a.path < b.path
		}
		if a.line != b.line {
			return a.line < b.line
		}
		return a.msg < b.msg
	})
	for _, x := range vios {
		fmt.Fprintf(w, "%s:%d: %s\n", x.path, x.line, x.msg)
	}
	if len(vios) > 0 {
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "roles used:")
	var us []string
	for p := range usedBy {
		us = append(us, p)
	}
	sort.Strings(us)
	for _, p := range us {
		var rs []string
		for c := range usedBy[p] {
			rs = append(rs, c)
		}
		sort.Strings(rs)
		fmt.Fprintf(w, "  %s: %s\n", p, strings.Join(rs, " "))
	}
	if len(us) == 0 {
		fmt.Fprintln(w, "  none")
	}
	fmt.Fprintln(w, "roles unused:")
	n := 0
	for _, r := range v.Roles {
		if !usedAny[r.Class] {
			fmt.Fprintf(w, "  %s\n", r.Class)
			n++
		}
	}
	if n == 0 {
		fmt.Fprintln(w, "  none")
	}
	fmt.Fprintln(w, "allow:")
	for _, a := range v.AllowClass {
		fmt.Fprintf(w, "  class %s: %s\n", a.Name, a.Reason)
	}
	for _, a := range v.AllowCSS {
		fmt.Fprintf(w, "  css %s: %s\n", a.Name, a.Reason)
	}
	if len(v.AllowClass)+len(v.AllowCSS) == 0 {
		fmt.Fprintln(w, "  none")
	}
	if len(vios) > 0 {
		return 1
	}
	return 0
}
