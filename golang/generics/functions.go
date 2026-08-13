package main

import "fmt"

// func (r *receiver)Name[type](args...)return{}

func Sum[T float32 | float64 | int | int64 | int16 | int8](a T, b T) T {
	return a + b
}

func Identity[A any](value A) A {
	return value
}

// MakePair Multiple parameters
func MakePair[K any, V any](a K, b V) any {
	return Pair[K, V]{a, b}
}

type Pair[K, V any] struct {
	K any
	V any
}

// constraints: 1
type Number interface {
	int | int8 | int16 | int32 | int64
}

func SumNumber[T Number](a T, b T) T {
	return a + b
}

// constraints: 2
func UseComparable[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// constraint: 3
type NumberConstraint interface {
	~int32 | ~int
}

type MyInt int32

func SumNumberConstraint[T NumberConstraint](a, b T) T {
	return a + b
}

func main() {
	fmt.Println(Sum(1, 2))
	fmt.Println(Identity(42))
	fmt.Println(Identity("hello"))
	fmt.Println(MakePair(1, 2))
	fmt.Println(SumNumber(1, 2))
	fmt.Println(SumNumberConstraint(MyInt(4), MyInt(5)))
}
