package identifiers

import "unicode"

func IsCamelCase(s string) bool {
	i := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			continue
		}
		if i == 0 && !unicode.IsUpper(r) || !IsAlphaChar(r) && !IsNumChar(r) {
			return false
		}
		i++
	}
	return true
}
