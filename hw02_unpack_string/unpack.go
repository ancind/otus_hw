package unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func Unpack(stingForUnpack string) (string, error) {
	runeSliceString := []rune(stingForUnpack)
	if len(runeSliceString) == 0 {
		return "", nil
	}

	if unicode.IsDigit(runeSliceString[0]) {
		return "", errors.New("element is number")
	}

	var stringBuilder strings.Builder

	canWriteNextElement := false
	canRepeatPreviousElement := true

	for _, element := range runeSliceString {
		switch {
		case canWriteNextElement:
			stringBuilder.WriteString(string(element))

			canWriteNextElement = false
			canRepeatPreviousElement = true
		case string(element) == "\\" && !canWriteNextElement:
			canWriteNextElement = true
			canRepeatPreviousElement = true
		case unicode.IsDigit(element):
			if !canRepeatPreviousElement {
				return "", errors.New("element is number")
			}

			repeatCount, _ := strconv.Atoi(string(element))
			if repeatCount > 0 {
				stringBuilder.WriteString(repeatPreviousElement(stringBuilder.String(), repeatCount))
			} else if repeatCount == 0 {
				convertedString := removeElement(stringBuilder.String(), 1)
				stringBuilder.Reset()
				stringBuilder.WriteString(convertedString)
			}

			canRepeatPreviousElement = false
		default:
			stringBuilder.WriteString(string(element))

			canRepeatPreviousElement = true
		}
	}

	return stringBuilder.String(), nil
}

func repeatPreviousElement(str string, repeatCount int) string {
	lastWrittenElement := getLastElementString(str)

	return strings.Repeat(lastWrittenElement, repeatCount-1)
}

func getLastElementString(str string) string {
	if len(str) == 0 {
		return ""
	}

	return str[len(str)-1:]
}

func removeElement(str string, removeCount int) string {
	return str[:len(str)-removeCount]
}
