package plan

import (
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// names maps generated stylesheet names back to the vocabulary's: a custom
// property to group.token, and a breakpoint's @media prelude to its name.
// The zero value keeps every name raw.
type names struct {
	tokens map[string]string // "--tg-color-ink" -> "color.ink"
	media  map[string]string // "@media (min-width: 64em)" -> "md"
}

// newNames takes the mapping from the vocabulary as vocab.Render writes it. A nil vocabulary gives raw names.
func newNames(v *vocab.Vocabulary) names {
	n := names{tokens: map[string]string{}, media: map[string]string{}}
	if v == nil {
		return n
	}
	for _, g := range v.Tokens {
		for _, t := range g.Tokens {
			n.tokens["--"+v.Project.Prefix+"-"+g.Name+"-"+t.Name] = g.Name + "." + t.Name
		}
	}
	for _, bp := range v.Breakpoints {
		n.media["@media (min-width: "+bp.Value+")"] = bp.Name
	}
	return n
}

// tokenName is a custom property's vocabulary name, or the property itself when the vocabulary lacks it.
func (n names) tokenName(prop string) string {
	if s, ok := n.tokens[prop]; ok {
		return s
	}
	return prop
}

// where is the context suffix of an item label: " at md" for a breakpoint, " in <prelude>" otherwise.
func (n names) where(ctx []string) string {
	var b strings.Builder
	for _, c := range ctx {
		if bp, ok := n.media[c]; ok {
			b.WriteString(" at " + bp)
		} else {
			b.WriteString(" in " + c)
		}
	}
	return b.String()
}

// rename writes var(--prop) of every known token as its vocabulary name in a line of output.
func (n names) rename(line string) string {
	for prop, name := range n.tokens {
		line = strings.ReplaceAll(line, "var("+prop+")", name)
	}
	return line
}
