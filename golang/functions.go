package main

type predicate func(int) bool

func filter(nums []int, p predicate) []int {
	result := make([]int, len(nums))
	for _, n := range nums {
		if p(n) {
			result = append(result, n)
		}
	}
	return result
}

func largerThenTwo(x int) bool {
	return x > 2
}

func filter_x(f func(int) bool) {
	f(5)
}

func utilizeFunctions() {
	filter([]int{1, 2, 3, 3}, largerThenTwo)
	filter([]int{1, 2, 3}, func(n int) bool { return n > 6 })

}
