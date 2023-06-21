package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	runeSlice := []rune(input)

	if unicode.IsDigit(runeSlice[0]) {
		return "", ErrInvalidString
	}

	var builder strings.Builder
	slowPtrIndex := 0
	fastPtrIndex := 0
	needToBeShielded := false

	for fastPtrIndex < len(runeSlice) {
		switch {
		case unicode.IsDigit(runeSlice[fastPtrIndex]) && runeSlice[fastPtrIndex] != '0':
			if unicode.IsDigit(runeSlice[fastPtrIndex-1]) && !needToBeShielded {
				return "", ErrInvalidString
			}

			for slowPtrIndex < fastPtrIndex-1 {
				builder.WriteRune(runeSlice[slowPtrIndex])
				slowPtrIndex++
			}

			builder.WriteString(strings.Repeat(string(runeSlice[fastPtrIndex-1]), int(runeSlice[fastPtrIndex]-'0')))

			slowPtrIndex = fastPtrIndex + 1
			needToBeShielded = false
		case unicode.IsDigit(runeSlice[fastPtrIndex]) && runeSlice[fastPtrIndex] == '0':
			if unicode.IsDigit(runeSlice[fastPtrIndex-1]) {
				return "", ErrInvalidString
			}

			for slowPtrIndex < fastPtrIndex-1 {
				builder.WriteRune(runeSlice[slowPtrIndex])
				slowPtrIndex++
			}

			slowPtrIndex = fastPtrIndex + 1
		case runeSlice[fastPtrIndex] == '\\':
			for slowPtrIndex < fastPtrIndex {
				builder.WriteRune(runeSlice[slowPtrIndex])
				slowPtrIndex++
			}

			slowPtrIndex = fastPtrIndex + 1
			fastPtrIndex++
			needToBeShielded = true
		}

		if fastPtrIndex == len(runeSlice)-1 {
			for slowPtrIndex < fastPtrIndex+1 {
				builder.WriteRune(runeSlice[slowPtrIndex])
				slowPtrIndex++
			}
		}

		fastPtrIndex++
	}

	return builder.String(), nil
}
