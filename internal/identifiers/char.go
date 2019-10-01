package identifiers

func IsAlphaChar(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func IsNumChar(r rune) bool {
	return '0' <= r && r <= '9'
}
