package vocab

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// maxSteps bounds one file's evaluation, so a runaway loop ends instead of hanging a build.
const maxSteps = 10_000_000

// fileOptions turns on top-level for and if; while and recursion stay off.
var fileOptions = &syntax.FileOptions{TopLevelControl: true}

// evaluator runs the entry file and the files it loads or sources.
// Files are named by their slash path relative to the entry's directory.
type evaluator struct {
	l     *loader
	root  string                         // the entry file's directory
	done  map[string]starlark.StringDict // evaluated files: a file runs once
	stack []string                       // files being evaluated, for cycle errors
}

// evalStarlark evaluates the entry file; each builtin call fills l.raw and l.kids.
func (l *loader) evalStarlark(entry string) []error {
	e := &evaluator{l: l, root: filepath.Dir(entry), done: map[string]starlark.StringDict{}}
	if _, err := e.run(filepath.Base(entry)); err != nil {
		return append(l.errs, errors.New(describe(err)))
	}
	return l.errs
}

// describe keeps Starlark's own file, line and column for an evaluation error.
func describe(err error) string {
	var ev *starlark.EvalError
	if errors.As(err, &ev) {
		return ev.Backtrace()
	}
	return err.Error()
}

func (e *evaluator) run(name string) (starlark.StringDict, error) {
	if g, ok := e.done[name]; ok {
		return g, nil
	}
	if i := indexOf(e.stack, name); i >= 0 {
		return nil, fmt.Errorf("load cycle: %s", strings.Join(append(e.stack[i:], name), " -> "))
	}
	src, err := os.ReadFile(filepath.Join(e.root, filepath.FromSlash(name)))
	if err != nil {
		return nil, err
	}
	e.stack = append(e.stack, name)
	defer func() { e.stack = e.stack[:len(e.stack)-1] }()
	th := &starlark.Thread{Name: name, Load: e.load, Print: func(*starlark.Thread, string) {}}
	th.SetMaxExecutionSteps(maxSteps)
	g, err := starlark.ExecFileOptions(fileOptions, th, name, src, e.builtins())
	if err != nil {
		return nil, err
	}
	e.done[name] = g
	return g, nil
}

func (e *evaluator) load(th *starlark.Thread, module string) (starlark.StringDict, error) {
	name, err := resolveModule(th.Name, module)
	if err != nil {
		return nil, err
	}
	return e.run(name)
}

// resolveModule makes module relative to the loading file; it must stay inside the entry's tree.
func resolveModule(from, module string) (string, error) {
	if path.IsAbs(module) {
		return "", fmt.Errorf("%q: path must be relative to the loading file", module)
	}
	name := path.Join(path.Dir(from), module)
	if !filepath.IsLocal(filepath.FromSlash(name)) {
		return "", fmt.Errorf("%q: path leaves the directory of the entry file", module)
	}
	return name, nil
}

func indexOf(s []string, x string) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}
