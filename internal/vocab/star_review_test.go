package vocab

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStepBudgetIsSharedByAllFiles(t *testing.T) {
	old := maxSteps
	t.Cleanup(func() { maxSteps = old })
	spin := "for i in range(100):\n    pass\n"
	one := map[string]string{"index.tg": projectLine + spin}
	if _, errs := Load(filepath.Join(writeTree(t, one), "index.tg")); len(errs) > 0 {
		t.Fatal(errs)
	}
	maxSteps = 3000 // one file fits, four do not
	files := map[string]string{"index.tg": projectLine + "source(\"a.tg\")\nsource(\"b.tg\")\nsource(\"c.tg\")\nsource(\"d.tg\")\n"}
	for _, n := range []string{"a.tg", "b.tg", "c.tg", "d.tg"} {
		files[n] = spin
	}
	if got := loadErr(t, files); !strings.Contains(got, "too many steps") {
		t.Fatalf("got %q, want too many steps", got)
	}
	if _, errs := Load(filepath.Join(writeTree(t, one), "index.tg")); len(errs) > 0 {
		t.Fatalf("one file must fit: %v", errs)
	}
}

func TestSourceInFunctionResolvesAgainstItsOwnFile(t *testing.T) {
	dir := writeTree(t, map[string]string{
		"index.tg": projectLine + "load(\"sub/a.tg\", \"f\")\nf()\n",
		"sub/a.tg": "def f():\n    source(\"b.tg\")\n",
		"sub/b.tg": "role(\"text\", \"b\", use = \"x\")\n",
	})
	if _, errs := Load(filepath.Join(dir, "index.tg")); len(errs) > 0 {
		t.Fatal(errs)
	}
}

func TestEvalErrorIsOneLinePerFrame(t *testing.T) {
	got := loadErr(t, map[string]string{
		"index.tg": projectLine + "def f():\n    role(\"text\", 5)\n\ndef g():\n    f()\n\ng()\n",
	})
	want := "index.tg:3:9: role: for parameter name: got int, want string\n  called from index.tg:6\n  called from index.tg:8"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSyntaxErrorIsOneLine(t *testing.T) {
	got := loadErr(t, map[string]string{"index.tg": "x = (\n"})
	if !strings.HasPrefix(got, "index.tg:") || strings.Contains(got, "Traceback") {
		t.Fatalf("got %q", got)
	}
}

func TestTokensKeywordRepeatingADictKeyFails(t *testing.T) {
	got := loadErr(t, map[string]string{"index.tg": projectLine + "tokens(\"c\", {\"a\": \"1\"}, a = \"2\")\n"})
	if !strings.Contains(got, "given twice") {
		t.Fatalf("got %q", got)
	}
}
