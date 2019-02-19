package util

//反转字符串
func ReverseStr(from string) string {
	runes := []rune(from)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
