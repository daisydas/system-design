package main

type Box[T any] struct {
	value T
}

func (b Box[T]) Get() T {
	return b.value
}

func (b Box[T]) Set(val T) {
	b.value = val
}

func (b Box[T]) DoSomething(u U) U {
	return u
}

type justStruct struct {
}

func (b justStruct) DoSomething[U any](u U) U {
	return u
}
