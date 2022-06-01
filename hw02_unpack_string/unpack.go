package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrStrange       = errors.New("very strange error of the digit convertation")
)

func Unpack(str string) (string, error) {
	runesIn := []rune(str)
	var b strings.Builder

	for i := 0; i < len(runesIn); i++ {
		if unicode.IsDigit(runesIn[i]) {
			return "", ErrInvalidString
		}
		symbol := string(runesIn[i])
		if runesIn[i] == 92 && i+1 < len(runesIn) { // escaping
			if unicode.IsDigit(runesIn[i+1]) || runesIn[i+1] == 92 { // only digits or backslash may be escaped
				symbol = string(runesIn[i+1])
				i++
			} else {
				return "", ErrInvalidString
			}
		}
		repeatCount := 1
		if i+1 < len(runesIn) && unicode.IsDigit(runesIn[i+1]) {
			repeatCountConverted, err := strconv.Atoi(string(runesIn[i+1]))
			if err != nil {
				return "", ErrStrange
			}
			repeatCount = repeatCountConverted
			i++
		}
		b.WriteString(strings.Repeat(symbol, repeatCount))
	}
	return b.String(), nil
}
