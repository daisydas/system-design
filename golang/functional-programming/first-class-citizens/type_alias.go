package main

import (
	"fmt"
	"math"
)

type phoneNumber string

type Person struct {
	name        string
	phoneNumber phoneNumber
}

func (p Person) SetPhoneNumber(number phoneNumber) {
	p.phoneNumber = number
}

type celsius float64

func (c celsius) display() {
	_ = fmt.Sprintf("%0.2f\n", c)
}

type temperature celsius

func (t *temperature) displayMe() {
	fmt.Println(90)
}

func main1() {

	c := celsius(45)
	x := c + 45
	fmt.Println(x)
	math.Abs(float64(c))

	t := temperature(x)
	t.displayMe()

	p := Person{
		name:        "xyz",
		phoneNumber: "1234",
	}

	fmt.Println(p)
	p.SetPhoneNumber("12345")

}
