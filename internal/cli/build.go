package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sneikki/tidgrid/internal/vocab"
)

func runBuild(args []string) int {
	v, code := vocab.LoadFlag("build", args)
	if v == nil {
		return code
	}
	out := v.Project.Out
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "tidgrid build:", err)
		return 2
	}
	if err := os.WriteFile(out, []byte(v.Render()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "tidgrid build:", err)
		return 2
	}
	fmt.Println("wrote", out)
	return 0
}
