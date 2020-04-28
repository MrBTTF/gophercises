package phone

import (
	"bytes"
	"unicode"
)

func Normalize(phone string) string {
	var buf bytes.Buffer
	for _, c := range phone {
		if unicode.IsDigit(c) {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
