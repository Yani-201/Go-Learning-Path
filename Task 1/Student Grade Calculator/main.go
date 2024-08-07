package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getinput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}

func createCard() card {
	reader := bufio.NewReader(os.Stdin)
	name, _ := getinput("Enter your full name:  ", reader)
	subjects, _ := getinput("Enter the number of courses you took this semester:  ", reader)

	s, err := strconv.Atoi(subjects)
	if err != nil {
		fmt.Println("The number of subjects must be a number. Please enter again.")
		createCard()
	}

	c := newCard(name, s)
	fmt.Println("Report Card is succesfully created.")
	return c
}

func Validateinput(inp string) float64{
	g, err := strconv.ParseFloat(inp, 64)
	if err != nil {
		fmt.Println("The Grade value must be a number.")
		reader := bufio.NewReader(os.Stdin)
		new, _ := getinput("Enter your grade for this Subject:  ", reader)

		return Validateinput(new)
	}
	if g < 0 || g > 100 {
		fmt.Println("The Grade value must be in a range from 0 to 100.")
		reader := bufio.NewReader(os.Stdin)
		new, _ := getinput("Enter your grade for this Subject:  ", reader)

		return Validateinput(new)
	}

	return g

}
func populateGrades(c *card){
	reader := bufio.NewReader(os.Stdin)
	for i:= 1; i <= c.Subjects; i++{
		prompt := fmt.Sprintf("Enter Name for Course %v:  ", i)
		Subj_Name, _ := getinput(prompt, reader)
		input_grade, _ := getinput("Enter your grade for this Subject:  ", reader)
		grade := Validateinput(input_grade)

		c.updateCard(Subj_Name, grade)
	}
	// name, _ := getinput("Enter your full name:", reader)

}
func main() {
	myCard := createCard()
	populateGrades(&myCard)
	fmt.Println(myCard.format())
}
