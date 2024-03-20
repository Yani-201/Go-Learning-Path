package main

import (
	"fmt"
)

type card struct {
	Student_Name string
	Subjects     int
	Grades       map[string]float64
	Total        float64
}

func newCard(name string, subject int) card {
	c := card{
		Student_Name: name,
		Subjects:     subject,
		Grades:       map[string]float64{},
		Total:        0,
	}

	return c

}

func (c card) format() string {
	fs := "Student's Report Card: \n"
	// var total float64

	fs += fmt.Sprintf("%-15v ... %v \n", "Student Name:", c.Student_Name)

	for k, v := range c.Grades {
		fs += fmt.Sprintf("%-15v ... %0.2f \n", k+":", v)
	}

	fs += fmt.Sprintf("%-15v ... %0.2f \n", "Average Grade:", c.Total/float64(c.Subjects))
	return fs
}

func (c *card) updateCard(subject_name string, grade float64) {
	c.Grades[subject_name] = grade
	c.Total += grade

}
