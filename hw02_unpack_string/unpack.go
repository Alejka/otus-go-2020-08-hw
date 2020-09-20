package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	valid "github.com/asaskevich/govalidator"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack unpacks the string according to the format string
func Unpack(str string) (string, error) {
	var err error

	if !utf8.ValidString(str) {
		err = ErrInvalidString
		return "", err
	}

	runesCount := utf8.RuneCountInString(str)
	if runesCount == 0 {
		return "", err
	}

	var unpackedStrBuilder strings.Builder

	var repeatCountStrBuilder strings.Builder
	repeatRune := rune(-1)
	repeatCount := 1

	escapedRune := '\\'
	isEscapeLastRune := false

	runeIndex := 0
	strRune := ""
	for _, rune := range str {
		runeIndex++
		strRune = string(rune)

		// first rune cannot be digit
		if runeIndex == 1 && valid.IsInt(strRune) {
			err = ErrInvalidString
			return "", err
		}
		// if escaped char locates before letter
		if isEscapeLastRune && !valid.IsInt(strRune) && rune != escapedRune {
			err = ErrInvalidString
			return "", err
		}

		if valid.IsInt(strRune) && !isEscapeLastRune { // if rune is digit (not escape)
			repeatCountStrBuilder.WriteRune(rune)
			continue
		} else {
			// if rune is "escape" character
			if rune == escapedRune && !isEscapeLastRune {
				// last rune cannot be `\`
				if runeIndex == runesCount {
					err = ErrInvalidString
					return "", err
				}

				isEscapeLastRune = true
			} else {
				// if rune is char: repeat previous rune
				if repeatRune != -1 {
					if repeatCountStrBuilder.Len() > 0 {
						repeatCount, err = strconv.Atoi(repeatCountStrBuilder.String())
						if err != nil {
							return "", err
						}
					} else {
						repeatCount = 1
					}
					unpackedStrBuilder.WriteString(strings.Repeat(string(repeatRune), repeatCount))
				}

				repeatRune = rune
				repeatCountStrBuilder.Reset()
				isEscapeLastRune = false
			}
		}
	}

	// process last rune
	if repeatRune != -1 {
		if repeatCountStrBuilder.Len() > 0 {
			repeatCount, err = strconv.Atoi(repeatCountStrBuilder.String())
			if err != nil {
				return "", err
			}
		} else {
			repeatCount = 1
		}
		unpackedStrBuilder.WriteString(strings.Repeat(string(repeatRune), repeatCount))
	}

	if err != nil {
		unpackedStrBuilder.Reset()
	}

	return unpackedStrBuilder.String(), err
}
