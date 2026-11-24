package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSrc(t *testing.T, toml string) (dir, src string) {
	t.Helper()
	dir = t.TempDir()
	src = filepath.Join(dir, "tidgrid.toml")
	if err := os.WriteFile(src, []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir, src
}

func TestConvertRejectsInvalidTOML(t *testing.T) {
	bad := []string{
		"[tokens]\nstray = \"x\"\n",
		"[allow.weird]\na = \"b\"\n",
		"[role]\nx = \"y\"\n",
		"[element.a]\nuse = 5\n",
		"[allow.class]\nn = 3\n",
	}
	for _, extra := range bad {
		dir, src := writeSrc(t, base+extra)
		if Run([]string{"convert", src}) != 2 {
			t.Errorf("%q: want exit 2", extra)
		}
		if _, err := os.Stat(filepath.Join(dir, "index.tg")); err == nil {
			t.Errorf("%q: index.tg written", extra)
		}
	}
}

func TestConvertLeavesNothingWhenAFileExists(t *testing.T) {
	dir, src := writeSrc(t, base)
	roles := filepath.Join(dir, "roles.tg")
	os.WriteFile(roles, []byte("mine"), 0o644)
	if Run([]string{"convert", src}) != 2 {
		t.Fatal("want exit 2")
	}
	for _, n := range []string{"index.tg", "tokens.tg"} {
		if _, err := os.Stat(filepath.Join(dir, n)); err == nil {
			t.Errorf("%s left behind", n)
		}
	}
	if b, _ := os.ReadFile(roles); string(b) != "mine" {
		t.Fatal("existing file overwritten")
	}
}

func TestConvertBindsGroupsAndSplitsRoles(t *testing.T) {
	dir, src := writeSrc(t, "[project]\nout = \"o.css\"\n[tokens.color]\nink = \"#000\"\n[tokens.for]\na = \"1px\"\n[tokens.my-x]\na = \"1px\"\n"+
		"[role.text.badge]\nuse = \"color.ink\"\n")
	if Run([]string{"convert", src}) != 0 {
		t.Fatal("convert failed")
	}
	tok, _ := os.ReadFile(filepath.Join(dir, "tokens.tg"))
	for _, want := range []string{"color = tokens(\"color\"", "for_ = tokens(\"for\"", "my_x = tokens(\"my-x\""} {
		if !strings.Contains(string(tok), want) {
			t.Errorf("tokens.tg lacks %q:\n%s", want, tok)
		}
	}
	roles, _ := os.ReadFile(filepath.Join(dir, "roles.tg"))
	if !strings.Contains(string(roles), "role(\n    \"text\", \"badge\",\n    use = ") {
		t.Errorf("role not over several lines:\n%s", roles)
	}
}

func TestMissingEntryAndExtraArgsSay(t *testing.T) {
	t.Chdir(t.TempDir())
	if Run([]string{"build"}) != 2 || Run([]string{"build", "a.tg", "b.tg"}) != 2 {
		t.Fatal("want exit 2")
	}
}
