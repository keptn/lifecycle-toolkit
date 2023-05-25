package operatorcommon

import "strings"

func CreateResourceName(maxLen int, minSubstrLen int, str ...string) string {
	for len(str)*minSubstrLen > maxLen {
		minSubstrLen = minSubstrLen / 2
	}
	for i := 0; i < len(str)-1; i++ {
		newStr := strings.Join(str, "-")
		if len(newStr) > maxLen {
			if len(str[i]) < minSubstrLen {
				continue
			}
			cut := len(newStr) - maxLen
			if cut > len(str[i])-minSubstrLen {
				str[i] = str[i][:minSubstrLen]
			} else {
				str[i] = str[i][:len(str[i])-cut]
			}
		} else {
			return strings.ToLower(newStr)
		}
	}

	return strings.ToLower(strings.Join(str, "-"))
}
