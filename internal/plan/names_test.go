package plan

import (
	"slices"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

// TestVocabularyNames checks that tokens, breakpoints and var() references print in the vocabulary's words.
func TestVocabularyNames(t *testing.T) {
	cfg := writeVocab(t, t.TempDir(), planBase)
	v, code := vocab.LoadFile(cfg)
	if v == nil {
		t.Fatalf("vocabulary does not load: exit %d", code)
	}
	cases := []struct {
		name, old, cur string
		want           []string
	}{
		{
			"token with two readers",
			`:root{--tg-color-ink:#000}.tg-text\(prose\){color:var(--tg-color-ink)}.tg-link\(page\){color:var(--tg-color-ink)}`,
			`:root{--tg-color-ink:#222}.tg-text\(prose\){color:var(--tg-color-ink)}.tg-link\(page\){color:var(--tg-color-ink)}`,
			[]string{"changed token color.ink: #000 -> #222", "  reaches 2 classes: tg-text(prose) tg-link(page)"},
		},
		{
			"removed token with no readers",
			`:root{--tg-color-ink:#000;--tg-space-md:1rem}`,
			`:root{--tg-color-ink:#000}`,
			[]string{"removed token space.md: 1rem"},
		},
		{
			"breakpoint and var() reference, unknown token kept raw",
			`:root{--tg-gone:1}@media (min-width: 64em){.tg-surface\(card\){padding:var(--tg-space-md)}}`,
			`@media (min-width: 64em){.tg-surface\(card\){padding:var(--tg-space-lg)}}`,
			[]string{"changed class tg-surface(card) at md: padding: space.md -> space.lg", "removed token --tg-gone: 1"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Diff(c.old, c.cur, v)
			if !slices.Equal(got, c.want) {
				t.Fatalf("got\n%q\nwant\n%q", got, c.want)
			}
		})
	}
}
