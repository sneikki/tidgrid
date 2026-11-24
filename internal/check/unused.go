package check

import (
	"fmt"
	"path/filepath"
)

// unusedRoleMsg names a role no scanned source references, and the ways out.
const unusedRoleMsg = "role %s is not used in the scanned sources: use it, remove it, " +
	"or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()"

// judgeUnusedRoles reports each role that no scanned source references, once.
// A TOML role is reported against the vocabulary file; a .tg role against its own declaring call.
func (c *projectScan) judgeUnusedRoles(use *roleUse) {
	for _, r := range c.v.Roles {
		if use.any[r.Class] {
			continue
		}
		file, line := c.v.Path, 1
		if r.File != "" {
			file, line = filepath.Join(filepath.Dir(c.v.Path), filepath.FromSlash(r.File)), r.Line
		}
		c.add(c.rel(c.abs(file)), line, fmt.Sprintf(unusedRoleMsg, r.Class))
	}
}
