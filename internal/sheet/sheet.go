// Package sheet writes one self-contained HTML page that shows every token,
// role and atom of the vocabulary with the project's own stylesheet.
package sheet

import (
	"flag"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// layout styles the sheet only. Every selector is a tgs- class, so the
// samples keep the project's own look.
const layout = `
.tgs-page{box-sizing:border-box;max-width:72rem;margin:0 auto;padding:1rem}
.tgs-title{font-size:1.5rem;font-weight:700;margin:1rem 0}
.tgs-h{font-size:1.125rem;font-weight:700;margin:2.5rem 0 .75rem;padding-bottom:.25rem;border-bottom:1px solid #8886}
.tgs-sub{font-size:.95rem;font-weight:700;margin:1.25rem 0 .5rem;opacity:.7}
.tgs-desc{margin:.25rem 0;font-size:.9rem}
.tgs-row{display:flex;flex-wrap:wrap;align-items:center;gap:.25rem .75rem;margin:.4rem 0}
.tgs-chip{display:inline-block;width:3rem;height:2rem;border:1px solid #8886;border-radius:.25rem}
.tgs-item{margin:.6rem 0}
.tgs-bar{height:.75rem;max-width:100%;background:#8888;border-radius:.25rem;margin-bottom:.25rem}
.tgs-value{font:.85rem/1.4 ui-monospace,monospace;opacity:.8}
.tgs-dl{display:grid;grid-template-columns:max-content minmax(0,1fr);gap:.3rem 1rem;margin:0}
.tgs-dl>dt,.tgs-dl>dd{margin:0}
.tgs-role{display:grid;grid-template-columns:minmax(0,1fr) minmax(0,1fr);gap:.5rem 1.5rem;align-items:center;padding:.75rem 0;border-bottom:1px solid #8883}
.tgs-sample{min-width:0;display:flex;align-items:flex-start;position:relative;contain:layout}
.tgs-fill{display:inline-block;padding:1rem 1.5rem}
.tgs-meta{min-width:0}
.tgs-meta>code,.tgs-toggle>code{display:block;word-break:break-all;font-size:.85rem}
.tgs-toggle{display:flex;align-items:center;gap:.5rem}
.tgs-missing{color:#c33;font-weight:600}
@media (max-width:40em){.tgs-role{grid-template-columns:minmax(0,1fr)}}
`

const sentence = "The quick brown fox jumps over the lazy dog."

var styleEnd = regexp.MustCompile(`(?i)</style`)

// safeCSS makes CSS safe to inline in a <style> element: a token value or a
// font file holding "</style" must not close it. Other "<" stay, the
// stylesheet uses them in range media queries.
func safeCSS(css string) string {
	return styleEnd.ReplaceAllString(css, `<\/style`)
}

func hasPadding(r *vocab.Rule) bool {
	if r == nil {
		return false
	}
	for _, n := range r.Nodes {
		if strings.HasPrefix(n.Prop, "padding") {
			return true
		}
	}
	return false
}

// textKinds get a sentence as sample; other kinds get a short label.
var textKinds = map[string]bool{"text": true, "emphasis": true}

const script = `<script>
for (const box of document.querySelectorAll('input[data-state]')) {
  box.addEventListener('change', () => document.body.classList.toggle(box.dataset.state, box.checked));
}
</script>
`

// Run is the sheet command: it writes the self-contained style sheet page.
func Run(args []string) int {
	fs := flag.NewFlagSet("sheet", flag.ContinueOnError)
	out := fs.String("o", "tidgrid-sheet.html", "output file")
	v, code := vocab.LoadFlagSet(fs, args)
	if v == nil {
		return code
	}
	if err := os.WriteFile(*out, []byte(page(v)), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "tidgrid sheet:", err)
		return 2
	}
	fmt.Println("wrote", *out)
	return 0
}

// page is one self-contained page. The project's stylesheet and the
// allow.css files that exist are inlined, so the project's own rules style
// the samples. Output depends only on the vocabulary.
func page(v *vocab.Vocabulary) string {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	b.WriteString("<title>tidgrid sheet</title>\n")
	b.WriteString("<style>" + layout + "</style>\n")
	b.WriteString("<style>\n" + safeCSS(v.Render()) + "</style>\n")
	found := inlineAllowedCSS(v, &b)
	b.WriteString("</head>\n")

	var stateClasses []string
	for _, r := range v.Roles {
		if r.Kind == "state" {
			stateClasses = append(stateClasses, r.Class)
		}
	}
	fmt.Fprintf(&b, "<body class=\"%s\">\n", html.EscapeString(strings.Join(stateClasses, " ")))
	b.WriteString("<main class=\"tgs-page\">\n<h1 class=\"tgs-title\">tidgrid sheet</h1>\n")
	tokens(v, &b)
	roles(v, &b)
	atoms(v, &b)
	allow(v, found, &b)
	b.WriteString("</main>\n" + script + "</body>\n</html>\n")
	return b.String()
}

// inlineAllowedCSS inlines each allow.css file that exists, in file order, and
// reports which ones were found. Paths are relative to the vocabulary.
func inlineAllowedCSS(v *vocab.Vocabulary, b *strings.Builder) map[string]bool {
	found := map[string]bool{}
	// os.Root refuses a path or symlink that leaves the vocabulary's
	// directory: the sheet must not carry a file from elsewhere on the disk.
	root, err := os.OpenRoot(filepath.Dir(v.Path))
	if err != nil {
		return found
	}
	defer root.Close()
	for _, a := range v.AllowCSS {
		f, err := root.Open(filepath.FromSlash(a.Name))
		if err != nil {
			continue
		}
		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			continue
		}
		found[a.Name] = true
		b.WriteString("<style>\n" + safeCSS(string(data)) + "</style>\n")
	}
	return found
}

func tokens(v *vocab.Vocabulary, b *strings.Builder) {
	pre := v.Project.Prefix
	b.WriteString("<section>\n<h2 class=\"tgs-h\">tokens</h2>\n")
	for _, g := range v.Tokens {
		fmt.Fprintf(b, "<h3 class=\"tgs-sub\">%s</h3>\n", html.EscapeString(g.Name))
		for _, t := range g.Tokens {
			name := html.EscapeString(t.Name)
			value := html.EscapeString(t.Value)
			ref := html.EscapeString("var(--" + pre + "-" + g.Name + "-" + t.Name + ")")
			switch g.Name {
			case "color":
				fmt.Fprintf(b, "<div class=\"tgs-row\"><span class=\"tgs-chip\" style=\"background: %s\"></span><code>%s</code><span class=\"tgs-value\">%s</span></div>\n", ref, name, value)
			case "space", "size":
				fmt.Fprintf(b, "<div class=\"tgs-item\"><div class=\"tgs-bar\" style=\"width: %s\"></div><code>%s</code> <span class=\"tgs-value\">%s</span></div>\n", ref, name, value)
			default:
				fmt.Fprintf(b, "<dl class=\"tgs-dl\"><dt><code>%s</code></dt><dd class=\"tgs-value\">%s</dd></dl>\n", name, value)
			}
		}
	}
	b.WriteString("</section>\n")
}

func roles(v *vocab.Vocabulary, b *strings.Builder) {
	b.WriteString("<section>\n<h2 class=\"tgs-h\">roles</h2>\n")
	kinds, byKind := v.RoleKinds()
	for _, k := range kinds {
		if k == "state" {
			b.WriteString("<h3 class=\"tgs-sub\">set by scripts</h3>\n")
			for _, r := range byKind[k] {
				fmt.Fprintf(b, "<div class=\"tgs-role\"><label class=\"tgs-toggle\"><input type=\"checkbox\" data-state=\"%s\" checked><code>%s</code></label><p class=\"tgs-desc\">%s</p></div>\n",
					html.EscapeString(r.Class), html.EscapeString(r.Class), html.EscapeString(r.Use))
			}
			continue
		}
		fmt.Fprintf(b, "<h3 class=\"tgs-sub\">%s</h3>\n", html.EscapeString(k))
		for _, r := range byKind[k] {
			var sample string
			if textKinds[r.Kind] {
				sample = "<p class=\"" + html.EscapeString(r.Class) + "\">" + sentence + "</p>"
			} else {
				// Roles leave padding to the atoms, so a bare sample is a thin strip.
				// Pad the inside unless the role brings its own.
				inner := html.EscapeString(r.Name)
				if !hasPadding(r.Rule) {
					inner = "<span class=\"tgs-fill\">" + inner + "</span>"
				}
				sample = "<div class=\"" + html.EscapeString(r.Class) + "\">" + inner + "</div>"
			}
			fmt.Fprintf(b, "<div class=\"tgs-role\"><div class=\"tgs-sample\">%s</div><div class=\"tgs-meta\"><code>%s</code><p class=\"tgs-desc\">%s</p></div></div>\n",
				sample, html.EscapeString(r.Class), html.EscapeString(r.Use))
		}
	}
	b.WriteString("</section>\n")
}

func atoms(v *vocab.Vocabulary, b *strings.Builder) {
	b.WriteString("<section>\n<h2 class=\"tgs-h\">atoms</h2>\n")
	if len(v.Breakpoints) > 0 {
		fmt.Fprintf(b, "<p class=\"tgs-desc\">Every atom also takes a breakpoint prefix, mobile first: %s:%s-row applies from %s up.</p>\n",
			html.EscapeString(v.Breakpoints[0].Name), html.EscapeString(v.Project.Prefix), html.EscapeString(v.Breakpoints[0].Name))
	}
	b.WriteString("<dl class=\"tgs-dl\">\n")
	for i := 0; i < len(v.Atoms); {
		j := i + 1
		for j < len(v.Atoms) && v.Atoms[j].Family == v.Atoms[i].Family {
			j++
		}
		base, _, _ := strings.Cut(v.Atoms[i].Class, "(")
		fmt.Fprintf(b, "<dt><code>%s</code></dt>", html.EscapeString(base))
		if v.Atoms[i].Value != "" {
			b.WriteString("<dd>")
			for _, a := range v.Atoms[i:j] {
				fmt.Fprintf(b, "<code>%s</code> ", html.EscapeString(a.Class))
			}
			b.WriteString("</dd>")
		}
		b.WriteString("\n")
		i = j
	}
	b.WriteString("</dl>\n</section>\n")
}

func allow(v *vocab.Vocabulary, found map[string]bool, b *strings.Builder) {
	b.WriteString("<section>\n<h2 class=\"tgs-h\">allow</h2>\n<h3 class=\"tgs-sub\">class</h3>\n<dl class=\"tgs-dl\">\n")
	for _, a := range v.AllowClass {
		fmt.Fprintf(b, "<dt><code>%s</code></dt><dd>%s</dd>\n", html.EscapeString(a.Name), html.EscapeString(a.Reason))
	}
	b.WriteString("</dl>\n<h3 class=\"tgs-sub\">css</h3>\n<dl class=\"tgs-dl\">\n")
	for _, a := range v.AllowCSS {
		note := ""
		if !found[a.Name] {
			note = " <span class=\"tgs-missing\">missing</span>"
		}
		fmt.Fprintf(b, "<dt><code>%s</code></dt><dd>%s%s</dd>\n", html.EscapeString(a.Name), html.EscapeString(a.Reason), note)
	}
	b.WriteString("</dl>\n</section>\n")
}
