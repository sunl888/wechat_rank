package utils

import "strings"

// wxName => wx_name
func Snake(str, delimiter string) string {
	if delimiter == "" {
		delimiter = "_"
	}
	sb := strings.Builder{}
	sb.WriteString(str[:1])

	for i := 1; i < len(str); i++ {
		if str[i] >= 'A' && str[i] <= 'Z' {
			sb.WriteString(delimiter)
			sb.WriteByte(str[i] + 32)
		} else {
			sb.WriteByte(str[i])
		}
	}
	return sb.String()
}
