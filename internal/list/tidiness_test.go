package list

import (
	"strings"
	"testing"

	"github.com/sneikki/tidgrid/internal/vocab"
)

func TestTidinessLinePrintedOnlyWhenSet(t *testing.T) {
	v, errs := vocab.Load("../vocab/vocabtest/example.tg")
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if strings.Contains(text(v), "# tidiness") {
		t.Fatal("tidiness printed when unset")
	}
	v.Tidiness = "full"
	if !strings.Contains(text(v), "full: nothing has margin, padding, border, heading size, list marker or link colour") {
		t.Fatal("full line missing")
	}
	v.Tidiness = "minimal"
	if !strings.Contains(text(v), "minimal: browser fixes only") {
		t.Fatal("minimal line missing")
	}
}
