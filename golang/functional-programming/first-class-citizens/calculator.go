package main

import (
	"fmt"
	"os"
	"strconv"
)

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) int {
	if b == 0 {
		panic("divide by zero")
	}
	return a / b
}

type operation func(a, b int) int

func main() {
	operations := map[string]operation{
		"+":  Add,
		"-":  Subtract,
		"*":  Multiply,
		"/":  Divide,
	}

	ops := os.Args[1]
	operand1, _ := strconv.Atoi(os.Args[2])
	operand2, _ := strconv.Atoi(os.Args[3])

	fmt.Println(operations[ops](operand1, operand2))
}
