package check

import (
	"regexp"
	"strings"
)

// mask blanks out what is not read (code samples), keeping offsets and newlines.
// open lists the offsets of a tidgrid-ignore-start that is never closed.
func mask(src, ext string) (out string, open []int) {
	b := []byte(src)
	// Markdown code goes first: an ignore marker a page only talks about, in
	// a fence or in backticks, must not open a region.
	if ext == ".md" || ext == ".mdx" {
		maskMarkdownCode(b, ext)
	}
	for _, re := range []*regexp.Regexp{rawBlock, ignoreRe} {
		for _, m := range re.FindAllIndex(b, -1) {
			if re == ignoreRe && !strings.HasSuffix(string(b[m[0]:m[1]]), "tidgrid-ignore-end") {
				open = append(open, m[0])
			}
			blank(b, m[0], m[1])
		}
	}
	return string(b), open
}

// blank overwrites b[s:e] with spaces, keeping newlines so line numbers hold.
func blank(b []byte, s, e int) {
	for i := s; i < e; i++ {
		if b[i] != '\n' {
			b[i] = ' '
		}
	}
}

// maskMarkdownCode blanks fenced blocks and, in .md, inline code spans.
func maskMarkdownCode(b []byte, ext string) {
	in, off := false, 0
	for _, ln := range strings.SplitAfter(string(b), "\n") {
		t := strings.TrimLeft(ln, " \t")
		fence := strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")
		switch {
		case fence || in:
			blank(b, off, off+len(ln))
		case ext == ".md":
			// ponytail: not in .mdx, where a backtick is a JS template
			// literal that can build a class name; fences cover MDX samples.
			for _, s := range codeSpans(ln) {
				blank(b, off+s[0], off+s[1])
			}
		}
		if fence {
			in = !in
		}
		off += len(ln)
	}
}

// codeSpans returns the inline code spans of one Markdown line: a run of n
// backticks up to the next run of exactly n.
func codeSpans(ln string) (spans [][2]int) {
	run := func(i int) int {
		j := i
		for j < len(ln) && ln[j] == '`' {
			j++
		}
		return j - i
	}
	for i := 0; i < len(ln); i++ {
		if ln[i] != '`' {
			continue
		}
		n := run(i)
		for j := i + n; j < len(ln); j++ {
			if ln[j] != '`' {
				continue
			}
			if m := run(j); m == n {
				spans = append(spans, [2]int{i, j + n})
				i = j
				break
			} else {
				j += m - 1
			}
		}
		i += n - 1
	}
	return
}
