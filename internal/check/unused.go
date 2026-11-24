package check

import "fmt"

// unusedRoleMsg names a role no scanned source references, and the ways out.
const unusedRoleMsg = "role %s is not used in the scanned sources: use it, remove it, " +
	"or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()"

// judgeUnusedRoles reports each role that no scanned source references, once,
// against the vocabulary file.
func (c *projectScan) judgeUnusedRoles(use *roleUse) {
	path := c.rel(c.abs(c.v.Path))
	for _, r := range c.v.Roles {
		if !use.any[r.Class] {
			c.add(path, 1, fmt.Sprintf(unusedRoleMsg, r.Class))
		}
	}
}
