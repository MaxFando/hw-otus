package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var builder strings.Builder
	slowPtr := 0
	fastPtr := 0

	if len(input) == 0 {
		return "", nil
	}

	if unicode.IsDigit(rune(input[0])) {
		return "", ErrInvalidString
	}

	for fastPtr < len(input) {
		if input[fastPtr] == '0' && input[slowPtr] == '0' {
			return "", ErrInvalidString
		}

		if input[fastPtr] >= '1' && input[fastPtr] <= '9' {
			builder.WriteString(input[slowPtr : fastPtr-1])
			builder.WriteString(strings.Repeat(string(input[fastPtr-1]), int(input[fastPtr]-'0')))

			slowPtr = fastPtr + 1
		} else if input[fastPtr] == '0' {
			builder.WriteString(input[slowPtr : fastPtr-1])

			slowPtr = fastPtr + 1
		} else if input[fastPtr] == '\\' {
			builder.WriteString(input[slowPtr:fastPtr])
			slowPtr = fastPtr + 1
			fastPtr++
		}

		if fastPtr == len(input)-1 {
			builder.WriteString(input[slowPtr : fastPtr+1])
		}

		fastPtr++
	}

	return builder.String(), nil
}
