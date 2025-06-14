package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	words := strings.Fields(text)

	return words
}

func main() {
	fmt.Println("Hello, World!")
}
