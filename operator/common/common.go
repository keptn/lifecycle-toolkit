package operatorcommon

import "strings"

// CreateResourceName is a function that concatenates the parts from the
// input and checks, if the resulting string matches the maxLen condition.
// If it does not match, it reduces the subparts, starting with the first
// one (but leaving its length at least in minSubstrLen so it's not deleted
// completely) adn continuing with the rest if needed.
// Let's take WorkloadInstance as an example (3 parts: app, workload, version).
// First the app name is reduced if needed (only to minSubstrLen),
// afterwards workload and the version is not reduced at all. This pattern is
// chosen to not reduce only one part of the name (that can be completely gone
// afterwards), but to include all of the parts in the resulting string.
func CreateResourceName(maxLen int, minSubstrLen int, str ...string) string {
	// if the minSubstrLen is too long for the number of parts,
	// needs to be reduced
	for len(str)*minSubstrLen > maxLen {
		minSubstrLen = minSubstrLen / 2
	}
	// looping through the subparts and checking if the resulting string
	// matches the maxLen condition
	for i := 0; i < len(str)-1; i++ {
		newStr := strings.Join(str, "-")
		if len(newStr) > maxLen {
			// if the part is already smaller than the minSubstrLen,
			// this part cannot be reduced anymore, so we continue
			if len(str[i]) <= minSubstrLen {
				continue
			}
			// part needs to be reduced
			cut := len(newStr) - maxLen
			// if the needed reduction is bigger than the allowed
			// reduction on the part, it's reduced to the minimum
			if cut > len(str[i])-minSubstrLen {
				str[i] = str[i][:minSubstrLen]
			} else {
				// the needed reduction can be completed fully on this
				// part, so it's reduced accordingly
				str[i] = str[i][:len(str[i])-cut]
			}
		} else {
			return strings.ToLower(newStr)
		}
	}

	return strings.ToLower(strings.Join(str, "-"))
}
