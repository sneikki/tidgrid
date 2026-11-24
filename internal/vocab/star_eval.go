package vocab

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// maxSteps bounds the whole evaluation, all files together, so a runaway loop
// ends instead of hanging a build. Memory cannot be capped in-process with
// go.starlark.net, so a huge allocation is not stopped.
var maxSteps uint64 = 10_000_000

// fileOptions turns on top-level for and if; while and recursion stay off.
var fileOptions = &syntax.FileOptions{TopLevelControl: true}

// evaluator runs the entry file and the files it loads or sources.
// Files are named by their slash path relative to the entry's directory.
type evaluator struct {
	l     *loader
	root  *os.Root                       // the entry file's directory; reads cannot leave it, symlinks included
	done  map[string]starlark.StringDict // evaluated files: a file runs once
	stack []string                       // files being evaluated, for cycle errors
	// threads are the threads of the files being evaluated; spent is the steps of finished ones.
	threads []*starlark.Thread
	spent   uint64
}

// remaining is the step budget left for the whole evaluation; at least 1, as 0 means unlimited.
func (e *evaluator) remaining() uint64 {
	used := e.spent
	for _, t := range e.threads {
		used += t.ExecutionSteps()
	}
	if used >= maxSteps {
		return 1
	}
	return maxSteps - used
}

// finish takes th off the running threads and charges its steps; the thread that resumes gets the smaller budget.
func (e *evaluator) finish(th *starlark.Thread) {
	e.threads = e.threads[:len(e.threads)-1]
	e.spent += th.ExecutionSteps()
	if n := len(e.threads); n > 0 {
		top := e.threads[n-1]
		top.SetMaxExecutionSteps(top.ExecutionSteps() + e.remaining())
	}
}

// evalStarlark evaluates the entry file; each builtin call fills l.raw and l.kids.
func (l *loader) evalStarlark(entry string) []error {
	root, err := os.OpenRoot(filepath.Dir(entry))
	if err != nil {
		return []error{err}
	}
	defer root.Close()
	e := &evaluator{l: l, root: root, done: map[string]starlark.StringDict{}}
	if _, err := e.run(filepath.Base(entry)); err != nil {
		return append(l.errs, describe(err)...)
	}
	return l.errs
}

func (e *evaluator) run(name string) (starlark.StringDict, error) {
	if g, ok := e.done[name]; ok {
		return g, nil
	}
	if i := indexOf(e.stack, name); i >= 0 {
		return nil, fmt.Errorf("load cycle: %s", strings.Join(append(e.stack[i:], name), " -> "))
	}
	src, err := e.read(name)
	if err != nil {
		return nil, err
	}
	e.stack = append(e.stack, name)
	defer func() { e.stack = e.stack[:len(e.stack)-1] }()
	th := &starlark.Thread{Name: name, Load: e.load, Print: func(*starlark.Thread, string) {}}
	th.SetMaxExecutionSteps(e.remaining())
	e.threads = append(e.threads, th)
	defer e.finish(th)
	g, err := starlark.ExecFileOptions(fileOptions, th, name, src, e.builtins())
	if err != nil {
		return nil, err
	}
	e.done[name] = g
	return g, nil
}

// read opens name through the root, so a symlink that points outside the tree is refused.
func (e *evaluator) read(name string) ([]byte, error) {
	f, err := e.root.Open(filepath.FromSlash(name))
	if err != nil {
		return nil, fmt.Errorf("%s: cannot read inside the entry file's directory: %w", name, err)
	}
	defer f.Close()
	return io.ReadAll(f)
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
