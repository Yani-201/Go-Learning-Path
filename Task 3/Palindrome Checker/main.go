package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func checkPalindrome(s string) bool {
	var check string
	for _, char := range s {
		if unicode.IsLetter(char) {
			check += string(unicode.ToLower(char))
		}
	}

	i := 0
	j := len(check) - 1

	for j > i {
		if check[i] != check[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func getinput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	sen, _ := getinput("Enter the sentense:  ", reader)
	fmt.Printf("Is the above sentence a palindrome? Answer: %v \n", checkPalindrome(sen))

}
