package plan

import (
	"slices"
	"testing"
)

// TestDiffCases covers the stylesheet shapes the diff must read: repeated selectors,
// repeated properties, at-rule contexts and an unescaped class name.
func TestDiffCases(t *testing.T) {
	cases := []struct {
		name, old, cur string
		want           []string
	}{
		{"same selector in two blocks", ".x{a:1}.x{b:2}", ".x{a:1}", []string{"changed class x: removed b: 2"}},
		{"repeated property", ".x{color:red;color:blue}", ".x{color:red}", []string{"changed class x: removed color: blue"}},
		{"change inside at-rule", "@media (min-width: 1px){.y{margin:0}}", "@media (min-width: 1px){.y{margin:1px}}", []string{"changed class y in @media (min-width: 1px): margin: 0 -> 1px"}},
		{"escaped class name", `.tg-surface\(card\){a:1}`, `.tg-surface\(card\){a:2}`, []string{"changed class tg-surface(card): a: 1 -> 2"}},
		{"token added", "", ":root{--a: 1}", []string{"added token --a: 1"}},
		{"element removed", "body{margin:0}", "", []string{"removed element body"}},
		{"order-only change", ".x{a:1;a:2}", ".x{a:2;a:1}", nil},
		{"identical", ".x{a:1}", ".x{a:1}", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Diff(c.old, c.cur)
			if !slices.Equal(got, c.want) {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestUnescape(t *testing.T) {
	cases := map[string]string{
		`.tg-surface\(card\)`: ".tg-surface(card)",
		`.md\:tg-row`:         ".md:tg-row",
		`.\31 0`:              ".10",
		`a\`:                  "a",
		`.é\(`:                ".é(",
	}
	for in, want := range cases {
		if got := unescape(in); got != want {
			t.Errorf("unescape(%q) = %q, want %q", in, got, want)
		}
	}
}
