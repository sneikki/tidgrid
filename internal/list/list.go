package list

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// listRules are the three lines of README "Working with an agent".
const listRules = "Style only with the classes `tidgrid list` prints. No <style>, no style= for appearance, no own CSS.\n" +
	"Run `tidgrid check` before every commit; it must exit 0.\n" +
	"A new role is a change to tidgrid.toml in its own commit, never mixed into page work.\n"

// Run is the list command: the agent-facing summary of the vocabulary.
func Run(args []string) int {
	v, code := vocab.LoadFlag("list", args)
	if v == nil {
		return code
	}
	fmt.Print(listText(v))
	return 0
}

// listText is the agent-facing summary: plain text, class names and use texts,
// no token values.
func listText(v *vocab.Vocabulary) string {
	var b strings.Builder
	listBreakpoints(v, &b)
	b.WriteString("\n# tokens\n")
	listTokens(v, &b)
	b.WriteString("\n# atoms\n")
	listAtoms(v, &b)
	b.WriteString("\n# roles\n")
	listRoles(v, &b)
	b.WriteString("\n# allow\n")
	listAllow(v, &b)
	b.WriteString("\n# rules\n")
	b.WriteString(listRules)
	return b.String()
}

func listBreakpoints(v *vocab.Vocabulary, b *strings.Builder) {
	b.WriteString("# breakpoints\n")
	if len(v.Breakpoints) == 0 {
		b.WriteString("none\n")
	} else {
		names := make([]string, 0, len(v.Breakpoints))
		for _, bp := range v.Breakpoints {
			names = append(names, bp.Name)
		}
		fmt.Fprintf(b, "%s (ascending). Any atom takes a prefix: %s:%s-row applies from %s up.\n",
			strings.Join(names, " "), names[0], v.Project.Prefix, names[0])
	}
}

func listTokens(v *vocab.Vocabulary, b *strings.Builder) {
	for _, g := range v.Tokens {
		names := make([]string, 0, len(g.Tokens))
		for _, t := range g.Tokens {
			names = append(names, t.Name)
		}
		fmt.Fprintf(b, "%s: %s\n", g.Name, strings.Join(names, " "))
	}
}

// listAtoms prints one line per atom family; valueless atoms share a line.
func listAtoms(v *vocab.Vocabulary, b *strings.Builder) {
	b.WriteString("S = a name in space, Z = a name in size, B = a breakpoint name (no breakpoint prefix on this one)\n")
	var statics []string
	flush := func() {
		if len(statics) > 0 {
			b.WriteString(strings.Join(statics, " ") + "\n")
			statics = nil
		}
	}
	for i := 0; i < len(v.Atoms); {
		j := i + 1
		for j < len(v.Atoms) && v.Atoms[j].Family == v.Atoms[i].Family {
			j++
		}
		first := v.Atoms[i]
		base, _, _ := strings.Cut(first.Class, "(")
		switch {
		case first.Value == "":
			statics = append(statics, base)
		case first.Below != "":
			flush()
			fmt.Fprintf(b, "%s(B)\n", base)
		case strings.HasPrefix(first.Decls[0][1], "var("):
			flush()
			fmt.Fprintf(b, "%s(%s)\n", base, scaleLetter(listScale(v, first.Decls[0][1])))
		default:
			flush()
			vals := make([]string, 0, j-i)
			for _, a := range v.Atoms[i:j] {
				vals = append(vals, a.Value)
			}
			fmt.Fprintf(b, "%s(%s)\n", base, strings.Join(rangeNums(vals), "|"))
		}
		i = j
	}
	flush()
}

func listRoles(v *vocab.Vocabulary, b *strings.Builder) {
	kinds, byKind := v.RoleKinds()
	for _, k := range kinds {
		if k == "state" {
			b.WriteString("## set by scripts\n")
		} else {
			fmt.Fprintf(b, "## %s\n", k)
		}
		for _, r := range byKind[k] {
			fmt.Fprintf(b, "%s  %s\n", r.Class, oneLine(r.Use))
		}
	}
}

func listAllow(v *vocab.Vocabulary, b *strings.Builder) {
	for _, a := range v.AllowClass {
		fmt.Fprintf(b, "class %s  %s\n", a.Name, a.Reason)
	}
	for _, a := range v.AllowCSS {
		fmt.Fprintf(b, "css %s  %s\n", a.Name, a.Reason)
	}
}

// listScale names the token group a valued atom reads, from its var() declaration.
func listScale(v *vocab.Vocabulary, decl string) string {
	for _, g := range v.Tokens {
		if strings.HasPrefix(decl, "var(--"+v.Project.Prefix+"-"+g.Name+"-") {
			return g.Name
		}
	}
	return ""
}

// scaleLetter is the short name the atoms section gives the two scales.
func scaleLetter(group string) string {
	if group == "size" {
		return "Z"
	}
	return "S"
}

// rangeNums folds runs of three or more consecutive integers: 1 2 3 4 5 6 12 -> 1-6 12.
func rangeNums(vals []string) []string {
	var out []string
	for i := 0; i < len(vals); {
		n, err := strconv.Atoi(vals[i])
		j := i
		for err == nil && j+1 < len(vals) {
			if m, e := strconv.Atoi(vals[j+1]); e != nil || m != n+(j+1-i) {
				break
			}
			j++
		}
		if j-i >= 2 {
			out = append(out, vals[i]+"-"+vals[j])
			i = j + 1
		} else {
			out = append(out, vals[i])
			i++
		}
	}
	return out
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
