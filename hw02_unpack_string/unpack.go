package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	valid "github.com/asaskevich/govalidator"
)

// ErrInvalidString shoud be use for testing.
var ErrInvalidString = errors.New("invalid string")

// Unpack unpacks the string according to the format string.
func Unpack(str string) (string, error) {
	if !utf8.ValidString(str) {
		return "", ErrInvalidString
	}

	if utf8.RuneCountInString(str) == 0 {
		return "", nil
	}

	return unpackValidString(str)
}

func unpackValidString(str string) (string, error) {
	var err error
	var unpackedStrBuilder strings.Builder
	var repeatCountStrBuilder strings.Builder // stores string with count of char repeat (ex. 15 and 7 for string "a15b7f")
	runesCount := utf8.RuneCountInString(str)
	strRune := ""            // the rune in string type
	strRepeatedRune := ""    // the repeated rune in string type
	runeToRepeat := rune(-1) // the rune which should repeat
	escapedRune := '\\'
	isEscapeLastRune := false

	runeIndex := 0
	for _, rune := range str {
		runeIndex++
		strRune = string(rune)

		if runeIndex == 1 && valid.IsInt(strRune) { // first rune cannot be digit
			return "", ErrInvalidString
		}

		if isEscapeLastRune && !valid.IsInt(strRune) && rune != escapedRune { // if escaped char locates before letter
			return "", ErrInvalidString
		}

		if valid.IsInt(strRune) && !isEscapeLastRune { // if rune is digit (not escape)
			repeatCountStrBuilder.WriteRune(rune)
			continue
		}

		if rune == escapedRune && !isEscapeLastRune { // if rune is "escape" character
			if runeIndex == runesCount { // last rune cannot be `\`
				return "", ErrInvalidString
			}
			isEscapeLastRune = true
			continue
		}

		if runeToRepeat != -1 { // if rune is char: repeat previous rune
			strRepeatedRune, err = repeatRune(runeToRepeat, repeatCountStrBuilder)
			if err != nil {
				return "", err
			}
			unpackedStrBuilder.WriteString(strRepeatedRune)
		}

		runeToRepeat = rune
		repeatCountStrBuilder.Reset()
		isEscapeLastRune = false
	}

	if runeToRepeat != -1 { // process last rune
		strRepeatedRune, err = repeatRune(runeToRepeat, repeatCountStrBuilder)
		if err != nil {
			return "", err
		}
		unpackedStrBuilder.WriteString(strRepeatedRune)
	}

	return unpackedStrBuilder.String(), err
}

func repeatRune(runeToRepeat rune, repeatCountStrBuilder strings.Builder) (string, error) {
	if repeatCountStrBuilder.Len() == 0 {
		return string(runeToRepeat), nil
	}

	repeatCount, err := strconv.Atoi(repeatCountStrBuilder.String())
	if err != nil {
		return "", err
	}
	return strings.Repeat(string(runeToRepeat), repeatCount), nil
}
