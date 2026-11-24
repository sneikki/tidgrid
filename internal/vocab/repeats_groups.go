package vocab

// pairGroups gives each pair of roles at most one group: the lines both hold, and every role that
// holds all of them. A line set already seen is skipped, and a group is kept when a pair fires.
func pairGroups(sets [][]int, pairs [][2]int) []repeatGroup {
	seen := map[string]bool{}
	var out []repeatGroup
	for i := range sets {
		for j := i + 1; j < len(sets); j++ {
			lines := intersect(sets[i], sets[j])
			if len(lines) == 0 || seen[key(lines)] {
				continue
			}
			seen[key(lines)] = true
			roles := holders(sets, lines)
			if fires(pairs, len(lines), len(roles)) {
				out = append(out, repeatGroup{lines: lines, roles: roles})
			}
		}
	}
	return out
}

func fires(pairs [][2]int, lines, roles int) bool {
	for _, p := range pairs {
		if lines >= p[0] && roles >= p[1] {
			return true
		}
	}
	return false
}

// intersect of two ascending id lists, ascending.
func intersect(a, b []int) []int {
	var out []int
	for i, j := 0, 0; i < len(a) && j < len(b); {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			out = append(out, a[i])
			i++
			j++
		}
	}
	return out
}

// holders is every role index whose set contains all of lines.
func holders(sets [][]int, lines []int) []int {
	var out []int
	for i, s := range sets {
		if len(intersect(lines, s)) == len(lines) {
			out = append(out, i)
		}
	}
	return out
}
