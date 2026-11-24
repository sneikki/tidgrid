package check

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

type violation struct {
	path string
	line int
	msg  string
}

// roleUse records which roles each file uses, for the summary.
type roleUse struct {
	byFile map[string]map[string]bool
	any    map[string]bool
}

func newRoleUse() *roleUse {
	return &roleUse{byFile: map[string]map[string]bool{}, any: map[string]bool{}}
}

func (u *roleUse) record(file, class string) {
	if u.byFile[file] == nil {
		u.byFile[file] = map[string]bool{}
	}
	u.byFile[file][class] = true
	u.any[class] = true
}

// sortedViolations orders by path, line, message so output is stable.
func (c *projectScan) sortedViolations() []violation {
	violations := c.violations
	sort.Slice(violations, func(i, j int) bool {
		first, second := violations[i], violations[j]
		if first.path != second.path {
			return first.path < second.path
		}
		if first.line != second.line {
			return first.line < second.line
		}
		return first.msg < second.msg
	})
	return violations
}

func writeViolations(w io.Writer, violations []violation) {
	for _, entry := range violations {
		fmt.Fprintf(w, "%s:%d: %s\n", entry.path, entry.line, entry.msg)
	}
	if len(violations) > 0 {
		fmt.Fprintln(w)
	}
}

func writeRolesUsed(w io.Writer, use *roleUse) {
	fmt.Fprintln(w, "roles used:")
	files := make([]string, 0, len(use.byFile))
	for p := range use.byFile {
		files = append(files, p)
	}
	sort.Strings(files)
	for _, p := range files {
		var roles []string
		for c := range use.byFile[p] {
			roles = append(roles, c)
		}
		sort.Strings(roles)
		fmt.Fprintf(w, "  %s: %s\n", p, strings.Join(roles, " "))
	}
	if len(files) == 0 {
		fmt.Fprintln(w, "  none")
	}
}

func writeAllowed(w io.Writer, v *vocab.Vocabulary) {
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
}
