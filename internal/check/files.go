package check

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// filesUnderScan walks every project.scan root and returns the absolute paths
// of the files found, sorted. A path it cannot read is reported, not skipped.
func (c *projectScan) filesUnderScan() []string {
	files := map[string]bool{}
	for _, root := range c.v.Project.Scan {
		filepath.WalkDir(root, c.visit(root, files))
	}
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// visit is the walk callback for one root: it adds files to the set.
func (c *projectScan) visit(root string, files map[string]bool) fs.WalkDirFunc {
	return func(p string, d fs.DirEntry, err error) error {
		if err != nil { // an unreadable path must not pass as clean
			msg := err.Error()
			if u := errors.Unwrap(err); u != nil {
				msg = u.Error()
			}
			c.add(c.rel(c.abs(p)), 1, "cannot scan: "+msg)
			return nil
		}
		if d.IsDir() {
			if p != root && skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if c.isSymlinkedDir(p, d) {
			return nil
		}
		files[c.abs(p)] = true
		return nil
	}
}

// isSymlinkedDir reports whether the entry is a symlink to a directory. The
// walk does not enter one, and passing over it in silence would leave its
// pages unchecked, so it is reported unless a scan root reaches it anyway.
func (c *projectScan) isSymlinkedDir(p string, d fs.DirEntry) bool {
	if d.Type()&fs.ModeSymlink == 0 || skipDirs[d.Name()] {
		return false
	}
	st, err := os.Stat(p)
	if err != nil || !st.IsDir() {
		return false
	}
	if !underScan(p, c.v.Project.Scan) {
		c.add(c.rel(c.abs(p)), 1, "symlinked directory is not scanned; add its target to project.scan")
	}
	return true
}

// underScan reports whether the directory a symlink points to is walked
// anyway, through one of the scan roots.
func underScan(link string, roots []string) bool {
	target, err := filepath.EvalSymlinks(link)
	if err != nil {
		return false
	}
	for _, r := range roots {
		if r, err = filepath.EvalSymlinks(r); err != nil {
			continue
		}
		if rel, err := filepath.Rel(r, target); err == nil && filepath.IsLocal(rel) {
			return true
		}
	}
	return false
}
