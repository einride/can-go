package identifiers

import "unicode"

func IsCamelCase(s string) bool {
	for i, r := range s {
		if i == 0 && !unicode.IsUpper(r) || !IsAlphaChar(r) && !IsNumChar(r) {
			return false
		}
	}
	return true
}
