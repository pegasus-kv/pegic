package parser

import (
	"strings"
)

func hasPrefixNoCase(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

func TrimPrefixNoCase(s, prefix string) string {
	ss := strings.TrimPrefix(strings.ToLower(s), strings.ToLower(prefix))
	return s[len(s)-len(ss):]
}
