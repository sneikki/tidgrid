package vocab

import (
	"path/filepath"
	"strings"
)

// loadAllow reads [allow.class] and [allow.css]: each entry needs a reason.
func (l *loader) loadAllow() {
	for _, kind := range []string{"class", "css"} {
		for _, k := range l.children("allow", kind) {
			p := []string{"allow", kind, k}
			s, _ := l.table("allow", kind)[k].(string)
			if strings.TrimSpace(s) == "" {
				l.errf(p, "reason required")
				continue
			}
			// sheet inlines these files: a path that leaves the project would
			// copy an arbitrary local file into a page people pass around.
			if kind == "css" && !filepath.IsLocal(filepath.FromSlash(k)) {
				l.errf(p, "path must stay inside the vocabulary's directory")
				continue
			}
			if kind == "class" {
				l.v.AllowClass = append(l.v.AllowClass, Allow{k, s})
			} else {
				l.v.AllowCSS = append(l.v.AllowCSS, Allow{k, s})
			}
		}
	}
}
