package hw02_unpack_string

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	errIsNumber = errors.New("element is number")
)

type repeater struct {
	nextElement bool
	repeat      bool
}

func (r *repeater) next() {
	r.nextElement = false
	r.repeat = true
}

func Unpack(str string) (string, error) {
	sliceString := []rune(str)
	if len(sliceString) == 0 {
		return "", nil
	}

	if unicode.IsDigit(sliceString[0]) {
		return "", errIsNumber
	}

	var stringBuilder strings.Builder
	repeater := repeater{false, true}

	for _, element := range sliceString {
		switch {
		case repeater.nextElement:
			stringBuilder.WriteString(string(element))
			repeater.next()
		case string(element) == "\\":
			repeater.nextElement = true
			repeater.repeat = true
		case unicode.IsDigit(element):
			if !repeater.repeat {
				return "", errIsNumber
			}

			repeatCount, _ := strconv.Atoi(string(element))
			if repeatCount > 0 {
				stringBuilder.WriteString(repeatPreviousElement(stringBuilder.String(), repeatCount))
			} else if repeatCount == 0 {
				convertedString := removeElement(stringBuilder.String(), 1)
				stringBuilder.Reset()
				stringBuilder.WriteString(convertedString)
			}

			repeater.repeat = false
		default:
			stringBuilder.WriteString(string(element))

			repeater.repeat = true
		}
	}

	return stringBuilder.String(), nil
}

func repeatPreviousElement(str string, repeatCount int) string {
	lastElement := lastElement(str)

	return strings.Repeat(lastElement, repeatCount-1)
}

func lastElement(str string) string {
	if len(str) == 0 {
		return ""
	}

	return str[len(str)-1:]
}

func removeElement(str string, removeCount int) string {
	return str[:len(str)-removeCount]
}
