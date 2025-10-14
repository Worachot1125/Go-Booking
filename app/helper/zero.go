package helper

import "strings"

func FirstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func FirstNonZero(a, b int64) int64 {
	if a > 0 {
		return a
	}
	return b
}
