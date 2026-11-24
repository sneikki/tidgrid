package plan

import "strings"

// item is one thing the diff compares: a class, an element rule, a keyframe stop or a token.
// A token's single declaration has an empty property and the token value.
type item struct {
	label string
	name  string // a class's name without the dot; empty for other items
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
	names                     names
}

func buildIndex(blocks []*block, nm names) index {
	ix := index{classes: newTable(), tokens: newTable(), elements: newTable(), names: nm}
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
	case key == ":root" && len(ctx) == 0:
		ix.root(b)
	case isKeyframes(key):
		ix.keyframes(key, b)
	case isAtRule(key):
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

// root files a top-level :root: its --* declarations are tokens, everything else
// in it, nested rules included, is an element rule.
func (ix index) root(b *block) {
	var rest []decl
	for _, d := range b.decls {
		if strings.HasPrefix(d.prop, "--") {
			ix.token(d)
		} else {
			rest = append(rest, d)
		}
	}
	if len(rest) > 0 {
		ix.own(":root", nil, rest)
	}
	for _, s := range b.subs {
		ix.add(s, nil, ":root")
	}
}

// token files one token; a repeated token takes its last value.
func (ix index) token(d decl) {
	it := ix.tokens.at(d.prop, "token "+ix.names.tokenName(d.prop))
	it.token = true
	it.decls = []decl{{value: d.value}}
}

// keyframes files each stop of a keyframes rule as its own element item.
func (ix index) keyframes(key string, b *block) {
	name := strings.TrimPrefix(key, "@keyframes ")
	for _, s := range b.subs {
		stop := unescape(s.key)
		it := ix.elements.at(identity(key, []string{stop}), "keyframes "+name+" at "+stop)
		it.decls = append(it.decls, s.decls...)
	}
}

// own adds decls to the item of owner in context ctx; the item exists even without decls.
func (ix index) own(owner string, ctx []string, decls []decl) {
	t := ix.elements
	name, isClass := strings.CutPrefix(owner, ".")
	if isClass {
		t = ix.classes
	}
	it := t.at(identity(owner, ctx), ix.label(owner, ctx))
	if isClass {
		it.name = name
	}
	it.decls = append(it.decls, decls...)
}

func identity(owner string, ctx []string) string {
	return strings.Join(ctx, "\x01") + "\x00" + owner
}

func (ix index) label(owner string, ctx []string) string {
	if name, ok := strings.CutPrefix(owner, "."); ok {
		return "class " + name + ix.names.where(ctx)
	}
	return "element " + owner + ix.names.where(ctx)
}

func appendCopy(ctx []string, key string) []string {
	out := make([]string, 0, len(ctx)+1)
	return append(append(out, ctx...), key)
}
