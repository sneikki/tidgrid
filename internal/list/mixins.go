package list

import (
	"fmt"
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// mixins prints each mixin, its sets with their lines as written and the
// roles that include them. Nothing is printed when the file has none.
func mixins(v *vocab.Vocabulary, b *strings.Builder) {
	if len(v.Mixins) == 0 {
		return
	}
	b.WriteString("\n# mixins\n")
	for _, m := range v.Mixins {
		fmt.Fprintf(b, "%s  %s\n", m.Name, oneLine(m.Use))
		for _, s := range m.Sets {
			fmt.Fprintf(b, "  %s: %s\n", s.Name, strings.Join(s.Lines, "; "))
			fmt.Fprintf(b, "    included by: %s\n", strings.Join(s.Roles, " "))
		}
	}
}
