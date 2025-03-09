package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	checkRune := rune(str[0])       // буфер для записи предидущего значения в цикле
	if unicode.IsDigit(checkRune) { // проверка на цифру в начале строки
		return "", ErrInvalidString
	}
	resultHw2 := strings.Builder{}
	var prVal rune // переменная для проверки последнего символа в строке

	for _, val := range str {
		if unicode.IsDigit(checkRune) { // проверка на число
			if unicode.IsDigit(val) {
				return "", ErrInvalidString
			}
		}
		if unicode.IsDigit(val) { // если текущий символ в цикле это цифра, ио записываем предидущий символ n раз
			resultHw2.WriteString(strings.Repeat(string(checkRune), int(val-'0')))
		} else if !unicode.IsDigit(checkRune) { // если не число то записываем символ
			resultHw2.WriteString(string(checkRune))
		}
		checkRune = val
		prVal = val
	}
	if !unicode.IsDigit(prVal) { // проверка последнего символа
		resultHw2.WriteString(string(prVal))
	}
	resStr := resultHw2.String()
	return resStr, nil
}
