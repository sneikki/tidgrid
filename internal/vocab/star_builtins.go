package vocab

import (
	"fmt"

	"go.starlark.net/starlark"
)

// builtins are the functions a vocabulary file calls. They live Go-side, so
// every file shares them and each call fills the loader in evaluation order.
func (e *evaluator) builtins() starlark.StringDict {
	b := map[string]func(*starlark.Thread, starlark.Tuple, []starlark.Tuple) (starlark.Value, error){
		"project":     e.project,
		"breakpoints": e.breakpoints,
		"tokens":      e.tokens,
		"element":     e.element,
		"keyframes":   e.keyframes,
		"role":        e.role,
		"allow_class": e.allow("class"),
		"allow_css":   e.allow("css"),
		"source":      e.source,
	}
	d := starlark.StringDict{}
	for name, fn := range b {
		name, fn := name, fn
		d[name] = starlark.NewBuiltin(name, func(th *starlark.Thread, _ *starlark.Builtin, a starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
			return fn(th, a, kw)
		})
	}
	return d
}

// at is "file:line" of the builtin call being run on th.
func at(th *starlark.Thread) string {
	pos := th.CallFrame(1).Pos
	return fmt.Sprintf("%s:%d", pos.Filename(), pos.Line)
}

func (e *evaluator) source(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var module string
	if err := starlark.UnpackArgs("source", args, kw, "path", &module); err != nil {
		return nil, err
	}
	name, err := resolveModule(th.Name, module)
	if err != nil {
		return nil, err
	}
	_, err = e.run(name)
	return starlark.None, err
}

func (e *evaluator) project(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var out starlark.String
	var scan, prefix, helpers starlark.Value
	if err := starlark.UnpackArgs("project", args, kw, "out", &out, "scan?", &scan, "prefix?", &prefix, "helpers?", &helpers); err != nil {
		return nil, err
	}
	p := []string{"project"}
	t := map[string]any{}
	e.l.putKey(t, p, "out", string(out))
	for _, f := range []struct {
		key string
		v   starlark.Value
	}{{"scan", scan}, {"prefix", prefix}, {"helpers", helpers}} {
		if f.v != nil {
			e.l.putKey(t, p, f.key, projectValue(f.v))
		}
	}
	return starlark.None, e.l.declare(p, t, at(th))
}

// projectValue converts a string or a list of strings; anything else stays unconverted for the loader to reject.
func projectValue(v starlark.Value) any {
	if s, ok := v.(starlark.String); ok {
		return string(s)
	}
	seq, ok := v.(starlark.Sequence)
	if !ok {
		return nil
	}
	items := []any{}
	for it, x := seq.Iterate(), starlark.Value(nil); it.Next(&x); {
		if s, ok := x.(starlark.String); ok {
			items = append(items, string(s))
		} else {
			items = append(items, nil)
		}
	}
	return items
}

func (e *evaluator) breakpoints(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	d, err := dictArg("breakpoints", args, kw)
	if err != nil {
		return nil, err
	}
	for _, kv := range d.Items() {
		k, ok := kv[0].(starlark.String)
		if !ok {
			return nil, fmt.Errorf("breakpoints: names must be strings, got %s", kv[0].Type())
		}
		v, _ := e.l.cell([]string{"breakpoints", string(k)}, kv[1])
		if err := e.l.declare([]string{"breakpoints", string(k)}, v, at(th)); err != nil {
			return nil, err
		}
	}
	return starlark.None, nil
}

func (e *evaluator) tokens(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("tokens: missing group name")
	}
	group, ok := args[0].(starlark.String)
	if !ok {
		return nil, fmt.Errorf("tokens: group name must be a string, got %s", args[0].Type())
	}
	d, err := dictArg("tokens", args[1:], kw)
	if err != nil {
		return nil, err
	}
	p := []string{"tokens", string(group)}
	t, err := e.l.newTable(p, "", d)
	if err != nil {
		return nil, err
	}
	if err := e.l.declare(p, t, at(th)); err != nil {
		return nil, err
	}
	return newTokenGroup(string(group), d), nil
}

func (e *evaluator) element(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var sel, use string
	var css *starlark.Dict
	if err := starlark.UnpackArgs("element", args, kw, "selector", &sel, "css", &css, "use?", &use); err != nil {
		return nil, err
	}
	p := []string{"element", sel}
	t, err := e.l.newTable(p, use, css)
	if err != nil {
		return nil, err
	}
	return starlark.None, e.l.declare(p, t, at(th))
}

func (e *evaluator) keyframes(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var name string
	var stops *starlark.Dict
	if err := starlark.UnpackArgs("keyframes", args, kw, "name", &name, "stops", &stops); err != nil {
		return nil, err
	}
	p := []string{"keyframes", name}
	t, err := e.l.newTable(p, "", stops)
	if err != nil {
		return nil, err
	}
	if err := e.l.declare(p, t, at(th)); err != nil {
		return nil, err
	}
	return starlark.String("{keyframes." + name + "}"), nil
}

func (e *evaluator) role(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var kind, name, use string
	var css *starlark.Dict
	if err := starlark.UnpackArgs("role", args, kw, "kind", &kind, "name", &name, "use?", &use, "css?", &css); err != nil {
		return nil, err
	}
	p := []string{"role", kind, name}
	t, err := e.l.newTable(p, use, css)
	if err != nil {
		return nil, err
	}
	return starlark.None, e.l.declare(p, t, at(th))
}

func (e *evaluator) allow(kind string) func(*starlark.Thread, starlark.Tuple, []starlark.Tuple) (starlark.Value, error) {
	return func(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
		var name, reason string
		if err := starlark.UnpackArgs("allow_"+kind, args, kw, "name", &name, "reason", &reason); err != nil {
			return nil, err
		}
		return starlark.None, e.l.declare([]string{"allow", kind, name}, reason, at(th))
	}
}

// dictArg is the optional dict argument plus keyword arguments, as one dict in call order.
func dictArg(fn string, args starlark.Tuple, kw []starlark.Tuple) (*starlark.Dict, error) {
	d := starlark.NewDict(len(kw))
	if len(args) > 1 {
		return nil, fmt.Errorf("%s: want one dict, got %d arguments", fn, len(args))
	}
	if len(args) == 1 {
		src, ok := args[0].(*starlark.Dict)
		if !ok {
			return nil, fmt.Errorf("%s: want a dict, got %s", fn, args[0].Type())
		}
		for _, kv := range src.Items() {
			d.SetKey(kv[0], kv[1])
		}
	}
	for _, kv := range kw {
		if err := d.SetKey(kv[0], kv[1]); err != nil {
			return nil, err
		}
	}
	return d, nil
}
