package main

import (
	"fmt"
	"math/rand"
)

func main() {

	x := randomClosure(func(i int, i2 int) int {
		return i + i2
	})(5)
	fmt.Println(x)
	randomFunctionCurrying(1)(2)(3)
}

func randomClosure(f func(a, b int) int) func(int) int {
	rand1 := rand.Int()
	return func(i int) int {
		return f(rand1, i)
	}
}

func randomFunctionCurrying(a int) func(int) func(int) int {
	return func(b int) func(int) int {
		return func(c int) int {
			return a + b + c
		}
	}
}

func randomFunctionNoCurrying(a, b, c int) int {
	return ((50 + a) + b) + c
}
