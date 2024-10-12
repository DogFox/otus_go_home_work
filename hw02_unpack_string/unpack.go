package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func isDigit(str string) bool {
	if string('0') <= str && str <= string('9') {
		return true
	}
	return false
}

func Unpack(str string) (string, error) {
	lenQueue := utf8.RuneCountInString(str)
	queue := make([]string, lenQueue)
	for i, v := range str {
		queue[i] = string(v)
	}

	result := strings.Builder{}
	shield := false

	for len(queue) > 0 {
		first := queue[0] // Первый элемент очереди
		second := "1"
		if len(queue) > 1 {
			second = queue[1]
		}

		// fmt.Println(first, " ", second, " len ", len(queue))
		// низя чтобы сразу чиселка была
		if isDigit(first) {
			if !shield {
				return "", ErrInvalidString
			}
		}

		if first == "\\" && !shield {
			// флаг поднять, скипаем
			shield = true
			queue = queue[1:]
			// fmt.Println("shield", shield)
			continue
		}

		times, err := strconv.Atoi(second)
		if err != nil {
			if shield && second != "n" && second != "\\" {
				return "", ErrInvalidString
			}
			// значит тут символ а не цифра, плюсуем и двигаем дальше
			result.WriteString(first)
			queue = queue[1:]
			shield = false

			continue
		}
		result.WriteString(strings.Repeat(first, times))
		// fmt.Println(result.String())

		// Переопределяем очередь сразу на 2 элемента, но если последний, то так и быть...
		if len(queue) >= 2 {
			queue = queue[2:]
		} else {
			queue = queue[1:]
		}
	}
	return result.String(), nil
}
