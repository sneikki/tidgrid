package check

import "os"

// Own CSS is a violation twice over: a stylesheet nobody allowed (here) and a
// <style> block (scanStyleTags). The built stylesheet is the one CSS file that
// must exist and match the vocabulary.

// checkStylesheet reports a CSS file that is neither the built output nor
// listed under allow.css.
func (c *checkRun) checkStylesheet(abs, rel string) {
	if !c.allowedCSS[abs] {
		c.add(rel, 1, "own CSS: stylesheet not under allow.css; make it roles in the vocabulary")
	}
}

// checkOutputFresh reports a built stylesheet that is missing or stale.
func (c *checkRun) checkOutputFresh() {
	out := c.v.Project.Out
	b, err := os.ReadFile(out)
	switch {
	case err != nil:
		c.add(c.rel(c.abs(out)), 1, "stylesheet missing; run `tidgrid build`")
	case string(b) != c.v.Render():
		c.add(c.rel(c.abs(out)), 1, "stylesheet differs from the vocabulary; run `tidgrid build`")
	}
}
