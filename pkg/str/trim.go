package str

import (
	"fmt"
)

func TrimChars(s string, max int) string {
	runes := []rune(s)

	if len(runes) <= max {
		return s
	}

	return string(runes[:max]) +
		fmt.Sprintf("...[+%dchars]", len(runes)-max)
}
