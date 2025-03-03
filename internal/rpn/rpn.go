package rpn

import (
	"errors"
	"strconv"

	"github.com/Se623/calc-base-api/pkg/lib"
)

// Переводит выражение в обратную польскую запись
func InfixToPostfix(expression string) ([]string, error) {
	var bufnum string          // Буфер для числа
	var rpnarr []string        // Выражение в виде ПН
	rpnstack := lib.Newstack() // Стэк

	for _, v := range expression {
		if _, err := strconv.Atoi(string(v)); err == nil || string(v) == "." { // Если встречаем цифру (или плавающую точку)...
			bufnum += string(v) // ...Добавляем число в буфер
		} else if string(v) == ")" { // Если встречаем правую скобку...
			if bufnum != "" {
				rpnarr = append(rpnarr, bufnum)
				bufnum = ""
			}

			for rpnstack.GetTop() != "(" {
				if rpnstack.GetTop() == "" {
					return nil, errors.New("one of the brackets is missing a pair")
				}
				rpnarr = append(rpnarr, rpnstack.Pop())
			}
			rpnstack.Pop()
		} else { // Если встречаем оператор...
			flag := false
			allops := [6]string{"/", "*", "-", "+", "(", ")"}
			for _, op := range allops {
				if op == string(v) {
					flag = true
					break
				}
			}
			if !flag {
				return nil, errors.New("detected illigal symbols")
			}

			if bufnum != "" {
				rpnarr = append(rpnarr, bufnum)
				bufnum = ""
			}

			if rpnstack.GetTop() == "(" || string(v) == "(" || rpnstack.GetTop() == "" ||
				((rpnstack.GetTop() == "+" || rpnstack.GetTop() == "-") && (string(v) == "*" || string(v) == "/")) {
				rpnstack.Push(string(v))
			} else {
				for !((rpnstack.GetTop() == "+" || rpnstack.GetTop() == "-") && (string(v) == "*" || string(v) == "/")) && rpnstack.GetTop() != "(" && rpnstack.GetTop() != "" {
					rpnarr = append(rpnarr, rpnstack.Pop())
				}
				rpnstack.Push(string(v))
			}
		}
	}
	if bufnum != "" {
		rpnarr = append(rpnarr, bufnum)
		bufnum = ""
	}
	for rpnstack.GetTop() != "" {
		rpnarr = append(rpnarr, rpnstack.Pop())
	}

	return rpnarr, nil
}
