package list

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

func TestListMixins(t *testing.T) {
	src := `project(out = "o.css")
type = mixin("type", use = "A text style.", sets = {
    "heading": {"font-size": "700", "&:hover": {"opacity": "0.5"}},
})
role("text", "a", use = "A.", include = [type.heading])
role("text", "b", use = "B.", include = [type.heading])
`
	path := filepath.Join(t.TempDir(), "index.tg")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	v, errs := vocab.Load(path)
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	want := "# mixins\ntype  A text style.\n  heading: font-size: 700; &:hover { opacity: 0.5 }\n    included by: tg-text(a) tg-text(b)\n"
	if out := text(v); !strings.Contains(out, want) {
		t.Fatalf("list lacks %q in:\n%s", want, out)
	}
}
