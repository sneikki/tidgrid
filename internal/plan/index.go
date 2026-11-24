package plan

import "strings"

// item is one thing the diff compares: a class, an element rule or a token.
// A token's single declaration has an empty property and the token value.
type item struct {
	label string
	token bool
	decls []decl
}

// table holds items in stylesheet order, keyed by identity.
type table struct {
	keys  []string
	items map[string]*item
}

func newTable() *table {
	return &table{items: map[string]*item{}}
}

// at returns the item under key, appending it to the order when it is new.
func (t *table) at(key, label string) *item {
	it, ok := t.items[key]
	if !ok {
		it = &item{label: label}
		t.items[key] = it
		t.keys = append(t.keys, key)
	}
	return it
}

// index is a stylesheet in comparable form: classes, tokens and element rules.
type index struct {
	classes, tokens, elements *table
}

func buildIndex(blocks []*block) index {
	ix := index{classes: newTable(), tokens: newTable(), elements: newTable()}
	for _, b := range blocks {
		ix.add(b, nil, "")
	}
	return ix
}

// add files b. ctx is the at-rule and nested-rule path above b; owner is the
// selector of the enclosing rule, or "" when b is outside any rule.
func (ix index) add(b *block, ctx []string, owner string) {
	key := unescape(b.key)
	switch {
	case key == "":
		return
	case owner != "":
		ix.nested(b, ctx, key, owner)
	case key == ":root":
		ix.addTokens(b.decls)
	case isAtRule(key) && !isKeyframes(key):
		ix.atRule(b, ctx, key)
	default:
		ix.own(key, ctx, b.decls)
		for _, s := range b.subs {
			ix.add(s, ctx, key)
		}
	}
}

// nested files a rule or at-rule inside owner under the extended context.
func (ix index) nested(b *block, ctx []string, key, owner string) {
	inner := appendCopy(ctx, key)
	if len(b.decls) > 0 {
		ix.own(owner, inner, b.decls)
	}
	for _, s := range b.subs {
		ix.add(s, inner, owner)
	}
}

// atRule files an at-rule outside any rule: its own declarations, then its
// blocks, which name their own selectors.
func (ix index) atRule(b *block, ctx []string, key string) {
	if len(b.decls) > 0 {
		ix.own(key, ctx, b.decls)
	}
	inner := appendCopy(ctx, key)
	for _, s := range b.subs {
		ix.add(s, inner, "")
	}
}

// own adds decls to the item of owner in context ctx; the item exists even without decls.
func (ix index) own(owner string, ctx []string, decls []decl) {
	t := ix.elements
	if strings.HasPrefix(owner, ".") {
		t = ix.classes
	}
	it := t.at(identity(owner, ctx), label(owner, ctx))
	it.decls = append(it.decls, decls...)
}

// addTokens files the declarations of :root, one token each; a repeated token takes its last value.
func (ix index) addTokens(decls []decl) {
	for _, d := range decls {
		it := ix.tokens.at(d.prop, "token "+d.prop)
		it.token = true
		it.decls = []decl{{value: d.value}}
	}
}

func identity(owner string, ctx []string) string {
	return strings.Join(ctx, "\x01") + "\x00" + owner
}

func label(owner string, ctx []string) string {
	l := "element " + owner
	if strings.HasPrefix(owner, ".") {
		l = "class " + owner[1:]
	}
	if len(ctx) > 0 {
		l += " in " + strings.Join(ctx, " ")
	}
	return l
}

func appendCopy(ctx []string, key string) []string {
	out := make([]string, 0, len(ctx)+1)
	return append(append(out, ctx...), key)
}
