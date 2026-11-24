// Package check reports source that uses classes, style values or own CSS
// outside the vocabulary, and a built stylesheet that is stale.
package check

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// Run is the check command: it reports source that uses classes or values outside the vocabulary.
func Run(args []string) int {
	v, code := vocab.LoadFlag("check", args)
	if v == nil {
		return code
	}
	return reportProject(v, os.Stdout)
}

// reportProject checks the project, writes the report to w and returns the
// exit code: 1 when any violation was found, else 0.
func reportProject(v *vocab.Vocabulary, w io.Writer) int {
	c := newProjectScan(v)
	use := c.scanFiles(c.filesUnderScan())
	c.judgeOutputFresh()
	c.judgeUnusedRoles(use)
	violations := c.sortedViolations()
	writeViolations(w, violations)
	writeRolesUsed(w, use)
	writeAllowed(w, v)
	if len(violations) > 0 {
		return 1
	}
	return 0
}

// projectScan is one check of one project: the violations found so far and the
// lookups the file readers share.
type projectScan struct {
	v            *vocab.Vocabulary
	dir          string // absolute directory of the vocabulary file; paths in the report are relative to it
	violations   []violation
	seenViolated map[violation]bool
	allowedClass map[string]bool
	allowedCSS   map[string]bool // absolute paths
	roleClass    map[string]bool
}

func newProjectScan(v *vocab.Vocabulary) *projectScan {
	c := &projectScan{v: v, seenViolated: map[violation]bool{}, allowedClass: map[string]bool{}, roleClass: map[string]bool{}}
	c.dir = c.abs(filepath.Dir(v.Path))
	for _, a := range v.AllowClass {
		c.allowedClass[a.Name] = true
	}
	c.allowedCSS = map[string]bool{c.abs(v.Project.Out): true}
	for _, a := range v.AllowCSS {
		c.allowedCSS[c.abs(filepath.Join(c.dir, a.Name))] = true
	}
	for _, r := range v.Roles {
		c.roleClass[r.Class] = true
	}
	return c
}

func (c *projectScan) abs(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return a
}

// rel is the path as the report prints it.
func (c *projectScan) rel(p string) string {
	if r, err := filepath.Rel(c.dir, p); err == nil {
		p = r
	}
	return filepath.ToSlash(p)
}

// add records a violation once; the same one found twice is reported once.
func (c *projectScan) add(path string, line int, msg string) {
	found := violation{path, line, msg}
	if !c.seenViolated[found] {
		c.seenViolated[found] = true
		c.violations = append(c.violations, found)
	}
}

// scanFiles judges each file by its kind and returns which roles were used.
func (c *projectScan) scanFiles(paths []string) *roleUse {
	use := newRoleUse()
	for _, p := range paths {
		ext := strings.ToLower(filepath.Ext(p))
		switch {
		case cssExts[ext]:
			c.judgeStylesheet(p, c.rel(p))
		case srcExts[ext]:
			c.scanSource(p, ext, use)
		}
	}
	return use
}

// scanSource reads one source file for class names and style attributes.
func (c *projectScan) scanSource(p, ext string, use *roleUse) {
	r := c.rel(p)
	raw, err := os.ReadFile(p)
	if err != nil {
		c.add(r, 1, "cannot read file: "+errors.Unwrap(err).Error())
		return
	}
	src, open := mask(string(raw), ext)
	f := newFileScan(c.v, src, c.allowedClass)
	f.report = func(line int, msg string) { c.add(r, line, msg) }
	f.usedClass = func(class string) {
		if c.roleClass[class] {
			use.record(r, class)
		}
	}
	for _, o := range open {
		f.report(f.line(o), "tidgrid-ignore-start is never closed; the rest of the file is not checked")
	}
	f.scan(c.v.Project.Helpers)
}
