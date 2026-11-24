package vocab

import (
	"errors"
	"fmt"

	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

// describe renders an evaluation error as `file:line:col: message`, one line,
// plus a `  called from file:line` line for every caller frame, innermost first.
func describe(err error) []error {
	var list resolve.ErrorList
	if errors.As(err, &list) {
		out := make([]error, len(list))
		for i, e := range list {
			out[i] = errors.New(e.Error())
		}
		return out
	}
	var chain []*starlark.EvalError
	for e := err; e != nil; e = errors.Unwrap(e) {
		if ev, ok := e.(*starlark.EvalError); ok {
			chain = append(chain, ev)
		}
	}
	if len(chain) == 0 {
		return []error{err}
	}
	return []error{errors.New(oneLine(chain))}
}

// oneLine is the message of the innermost error (last in chain) with the caller frames of all of them.
func oneLine(chain []*starlark.EvalError) string {
	inner := chain[len(chain)-1]
	msg := inner.Msg
	for i := len(chain) - 1; i >= 0; i-- {
		stack := userFrames(chain[i].CallStack)
		if i == len(chain)-1 && len(stack) > 0 {
			msg = fmt.Sprintf("%s: %s", stack[len(stack)-1].Pos, inner.Msg)
			stack = stack[:len(stack)-1]
		}
		for j := len(stack) - 1; j >= 0; j-- {
			msg += fmt.Sprintf("\n  called from %s:%d", stack[j].Pos.Filename(), stack[j].Pos.Line)
		}
	}
	return msg
}

// userFrames drops the builtin function on top of the stack, if any.
func userFrames(s starlark.CallStack) starlark.CallStack {
	if n := len(s); n > 0 && s[n-1].Pos.Filename() == "<builtin>" {
		return s[:n-1]
	}
	return s
}
