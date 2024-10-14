package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func isDigit(str string) bool {
	return string('0') <= str && str <= string('9')
}

func Unpack(str string) (string, error) {
	runes := []rune(str)
	lenStr := len(runes)
	queue := make([]string, lenStr)

	// fmt.Println(str, lenStr)
	for i := 0; i < lenStr; i++ {
		// fmt.Println(i, runes[i])
		queue[i] = string(runes[i])
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
