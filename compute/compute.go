package compute

import "fmt"

func SumDigits(n int) (int, error) {
	if n < 0 {
		return -1, fmt.Errorf("invalid input")
	}
	var sumRecur func(int) int
	sumRecur = func(n int) int {
		if n < 10 {
			return n
		}
		allButLast, last := split(n)
		return last + sumRecur(allButLast)
	}
	return sumRecur(n), nil
}

func split(n int) (int, int) {
	return n / 10, n % 10
}
