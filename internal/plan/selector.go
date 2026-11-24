package plan

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// unescape turns the CSS escapes of a selector back into the class name as
// written in pages: `\(` becomes `(`, `\31 ` becomes `1`.
func unescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			b.WriteByte(s[i])
			continue
		}
		if i+1 >= len(s) {
			break
		}
		if end := hexEnd(s, i+1); end > i+1 {
			n, _ := strconv.ParseUint(s[i+1:end], 16, 32)
			b.WriteRune(rune(n))
			if end < len(s) && s[end] == ' ' {
				end++
			}
			i = end - 1
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i+1:])
		b.WriteRune(r)
		i += size
	}
	return b.String()
}

// hexEnd returns the end of up to six hex digits starting at from.
func hexEnd(s string, from int) int {
	end := from
	for end < len(s) && end-from < 6 && isHex(s[end]) {
		end++
	}
	return end
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isAtRule(key string) bool {
	return strings.HasPrefix(key, "@")
}

func isKeyframes(key string) bool {
	return strings.HasPrefix(key, "@keyframes")
}
