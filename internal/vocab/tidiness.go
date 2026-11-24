package vocab

import (
	"embed"
	"fmt"

	"go.starlark.net/starlark"
)

//go:embed tidiness/minimal.css tidiness/full.css
var tidinessFS embed.FS

// TidinessLevels are the values tidiness() takes, in the order their files are written.
var TidinessLevels = []string{"minimal", "full"}

func (e *evaluator) tidiness(th *starlark.Thread, args starlark.Tuple, kw []starlark.Tuple) (starlark.Value, error) {
	var level starlark.Value
	if err := starlark.UnpackArgs("tidiness", args, kw, "level", &level); err != nil {
		return nil, err
	}
	s, _ := level.(starlark.String)
	rules, err := tidinessUpTo(string(s))
	if rules == nil && err == nil {
		return nil, fmt.Errorf("tidiness: level must be %q or %q, got %s", TidinessLevels[0], TidinessLevels[1], level.String())
	}
	if err != nil {
		return nil, fmt.Errorf("tidiness(%q): %w", string(s), err)
	}
	if e.l.tidinessAt != "" {
		return nil, fmt.Errorf("tidiness: called twice (%s and %s)", e.l.tidinessAt, at(th))
	}
	e.l.tidinessAt = at(th)
	e.l.v.Tidiness, e.l.v.tidyRules = string(s), rules
	return starlark.None, nil
}

// tidinessUpTo parses the files from the first level through level; nil rules
// and no error when level is not one of TidinessLevels.
func tidinessUpTo(level string) ([]*Rule, error) {
	var out []*Rule
	for _, lvl := range TidinessLevels {
		rules, err := parseLevel(lvl)
		if err != nil {
			return nil, err
		}
		out = append(out, rules...)
		if lvl == level {
			return out, nil
		}
	}
	return nil, nil
}

func parseLevel(level string) ([]*Rule, error) {
	src, err := tidinessFS.ReadFile("tidiness/" + level + ".css")
	if err != nil {
		return nil, err
	}
	rules, err := parseRules(string(src))
	if err != nil {
		return nil, fmt.Errorf("%s.css: %w", level, err)
	}
	return rules, nil
}

// tidinessRules are the chosen level's rules, parsed when tidiness() ran.
func (v *Vocabulary) tidinessRules() []*Rule { return v.tidyRules }
