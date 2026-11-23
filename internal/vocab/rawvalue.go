package vocab

import (
	"regexp"
	"strings"
)

// Every CSS length and time unit (not %, fr or angles).
const units = `px|rem|em|ex|ch|cap|ic|lh|rlh|rcap|rch|rex|ric|` +
	`vw|vh|vmin|vmax|vi|vb|svw|svh|svmin|svmax|svi|svb|lvw|lvh|lvmin|lvmax|lvi|lvb|dvw|dvh|dvmin|dvmax|dvi|dvb|` +
	`cqw|cqh|cqi|cqb|cqmin|cqmax|cm|mm|q|in|pt|pc|ms|s`

const colorNames = `aliceblue|antiquewhite|aqua|aquamarine|azure|beige|bisque|black|blanchedalmond|blue|blueviolet|brown|burlywood|` +
	`cadetblue|chartreuse|chocolate|coral|cornflowerblue|cornsilk|crimson|cyan|darkblue|darkcyan|darkgoldenrod|darkgray|darkgreen|` +
	`darkgrey|darkkhaki|darkmagenta|darkolivegreen|darkorange|darkorchid|darkred|darksalmon|darkseagreen|darkslateblue|darkslategray|` +
	`darkslategrey|darkturquoise|darkviolet|deeppink|deepskyblue|dimgray|dimgrey|dodgerblue|firebrick|floralwhite|forestgreen|` +
	`fuchsia|gainsboro|ghostwhite|gold|goldenrod|gray|green|greenyellow|grey|honeydew|hotpink|indianred|indigo|ivory|khaki|lavender|` +
	`lavenderblush|lawngreen|lemonchiffon|lightblue|lightcoral|lightcyan|lightgoldenrodyellow|lightgray|lightgreen|lightgrey|` +
	`lightpink|lightsalmon|lightseagreen|lightskyblue|lightslategray|lightslategrey|lightsteelblue|lightyellow|lime|limegreen|linen|` +
	`magenta|maroon|mediumaquamarine|mediumblue|mediumorchid|mediumpurple|mediumseagreen|mediumslateblue|mediumspringgreen|` +
	`mediumturquoise|mediumvioletred|midnightblue|mintcream|mistyrose|moccasin|navajowhite|navy|oldlace|olive|olivedrab|orange|` +
	`orangered|orchid|palegoldenrod|palegreen|paleturquoise|palevioletred|papayawhip|peachpuff|peru|pink|plum|powderblue|purple|` +
	`rebeccapurple|red|rosybrown|royalblue|saddlebrown|salmon|sandybrown|seagreen|seashell|sienna|silver|skyblue|slateblue|slategray|` +
	`slategrey|snow|springgreen|steelblue|tan|teal|thistle|tomato|turquoise|violet|wheat|white|whitesmoke|yellow|yellowgreen`

var (
	hexRe     = regexp.MustCompile(`#[0-9A-Fa-f]{3,8}\b`)
	colorFn   = regexp.MustCompile(`(?i)\b(rgb|rgba|hsl|hsla|hwb|lab|lch|oklab|oklch|color)\(`)
	unitRe    = regexp.MustCompile(`(?i)(?:^|[^\w.])-?(?:\d+\.?\d*|\.\d+)(?:e[+-]?\d+)?(?:` + units + `)\b`)
	nameRe    = regexp.MustCompile(`(?i)(?:^|[^\w-])(` + colorNames + `)(?:$|[^\w-])`)
	quotedRe  = regexp.MustCompile(`"[^"]*"|'[^']*'|url\([^)]*\)`) // not looked at
	percentRe = regexp.MustCompile(`(?:^|[^\w.])-?(?:\d+\.?\d*|\.\d+)%`)
)

// sizeProp holds a box size, a spacing or a grid track. A percentage there is
// a size like any other and lives in tokens; in a gradient stop, a colour mix,
// a radius or an offset (`left: 50%`) it is not.
func sizeProp(prop string) bool {
	switch strings.TrimPrefix(strings.TrimPrefix(prop, "min-"), "max-") {
	case "width", "height", "inline-size", "block-size", "flex-basis", "flex", "gap", "row-gap", "column-gap":
		return true
	}
	for _, p := range []string{"margin", "padding", "grid-template", "grid-auto-"} {
		if strings.HasPrefix(prop, p) {
			return true
		}
	}
	return false
}

// identProps hold names the author picks, not colours: `grid-area: gold` is legal.
func identProp(prop string) bool {
	switch prop {
	case "font-family", "animation-name", "container-name", "container", "view-transition-name", "anchor-name", "position-anchor",
		"timeline-scope", "will-change", "transition-property":
		return true
	}
	for _, p := range []string{"grid", "counter-", "scroll-timeline", "view-timeline"} {
		if strings.HasPrefix(prop, p) {
			return true
		}
	}
	return false
}

// rawProblem names a raw colour, length, time or percentage size in a value of
// property prop (empty for an at-rule key), or "". References are masked first.
func rawProblem(prop, s string) string {
	s = quotedRe.ReplaceAllString(refRe.ReplaceAllString(s, "_"), `""`)
	if m := hexRe.FindString(s); m != "" {
		return "raw colour " + m
	}
	if m := colorFn.FindString(s); m != "" {
		return "raw colour " + m + ")"
	}
	if !identProp(prop) {
		if m := nameRe.FindStringSubmatch(s); m != nil {
			return "raw colour " + m[1] + " (move it to [tokens])"
		}
	}
	if m := unitRe.FindString(s); m != "" {
		return "raw length or time " + strings.TrimLeft(m, " ,(") + " (move it to [tokens])"
	}
	if sizeProp(prop) {
		if m := percentRe.FindString(s); m != "" {
			return "raw size " + strings.TrimLeft(m, " ,(") + " (move it to [tokens])"
		}
	}
	return ""
}
