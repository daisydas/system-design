package main

type Maybe[T any] interface {
	Get() T
	GetOrElse(T) T
}

type JustMaybe[A any] struct {
	value A
}

func (j JustMaybe[T]) Get() T {
	return j.value
}

func (j JustMaybe[A]) GetOrElse(def A) A {
	return j.value
}

type NothingMaybe[A any] struct{}

func (n NothingMaybe[A]) GetOrElse(def A) A {
	return def
}

func (n NothingMaybe[A]) Get() A {
	return *new(A)
}

func Nothing[A any]() Maybe[A] {
	return NothingMaybe[A]{}
}
