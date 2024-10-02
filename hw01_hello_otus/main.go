package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	str := "Hello, OTUS!"
	reverseString := reverse.String(str)
	fmt.Println(reverseString)
}
