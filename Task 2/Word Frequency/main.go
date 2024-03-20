package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func wordCount(sentence string) map[string]int {
	count := make(map[string]int)
	words := strings.FieldsFunc(sentence, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	for _, word := range words {
		count[strings.ToLower(word)]++
	}
	return count
}

func getinput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	sen, _ := getinput("Enter your favourite sentense:  ", reader)
	fmt.Println("The words in your sentence are: ", wordCount(sen))

}
