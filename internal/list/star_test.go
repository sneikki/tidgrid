package list

import (
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

func TestListSameForTOMLAndStarlark(t *testing.T) {
	a, errs := vocab.Load("../vocab/vocabtest/example.toml")
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	b, errs := vocab.Load("../vocab/vocabtest/example.tg")
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if text(a) != text(b) {
		t.Fatal("list differs between example.toml and example.tg")
	}
}
