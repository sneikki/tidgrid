package check

import "regexp"

// File kinds: what is read for class names and style attributes, what counts
// as a stylesheet, and the directories that hold generated or foreign files.
var (
	srcExts  = map[string]bool{".html": true, ".htm": true, ".astro": true, ".jsx": true, ".tsx": true, ".vue": true, ".svelte": true, ".md": true, ".mdx": true, ".liquid": true, ".njk": true, ".js": true, ".mjs": true, ".cjs": true, ".ts": true}
	cssExts  = map[string]bool{".css": true, ".scss": true, ".sass": true, ".less": true, ".styl": true}
	skipDirs = map[string]bool{"node_modules": true, ".git": true, "dist": true, "build": true, "_site": true, ".astro": true, "vendor": true}
)

// Where a class name or a style value can start.
var (
	classAttr = regexp.MustCompile(`(?i)(?:^|[^\w.:-])((?:v-bind:|:)?class(?::list)?|className)\s*=\s*`)
	assignRe  = regexp.MustCompile(`\.\s*className\s*\+?=`)
	declRe    = regexp.MustCompile(`\b(?:const|let|var)\s+$`)
	classDir  = regexp.MustCompile(`(?i)(?:^|\s)class:([^\s=>/{}"'|]+)`)
	styleAttr = regexp.MustCompile(`(?i)(?:^|[^\w.:-])((?:v-bind:|:)?style)\s*=\s*`)
	styleDir  = regexp.MustCompile(`(?i)(?:^|\s)style:([A-Za-z-]+)`)
	styleTag  = regexp.MustCompile(`(?i)<style(?:[\s>/]|$)`)
)

// Reading a style value.
var (
	shortRe  = regexp.MustCompile(`[{,]\s*([A-Za-z]\w*)\s*(?:[,}]|$)`)
	propRe   = regexp.MustCompile(`(?:^|[;{,\s'"])(-?[A-Za-z][A-Za-z0-9-]*)['"]?\s*:`)
	capRe    = regexp.MustCompile(`[A-Z]`)
	appearRe = regexp.MustCompile(`^(?:(?:background|border|font|outline|text-decoration|mask|column-rule)(?:-.*)?|color|box-shadow|text-shadow|opacity|filter|backdrop-filter|letter-spacing|line-height|fill|stroke|accent-color|caret-color|mix-blend-mode)$`)
)

// Regions that are not read.
var (
	rawBlock = regexp.MustCompile(`(?s)\{%-?\s*(highlight|raw)\b.*?\{%-?\s*end(highlight|raw)\s*-?%\}`)
	ignoreRe = regexp.MustCompile(`(?s)tidgrid-ignore-start(?:.*?tidgrid-ignore-end|.*$)`)
)
