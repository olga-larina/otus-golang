package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	text := "Hello, OTUS!"
	textReversed := stringutil.Reverse(text)
	fmt.Print(textReversed)
}
