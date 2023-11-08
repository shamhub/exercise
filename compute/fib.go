package compute

func Fib(nth int) int {
	initialOffset, curr := 1, 0
	kth := 0 // current element position is 0

	var fibRecur func(pred, curr, kth int) int
	fibRecur = func(pred, curr int, kth int) int {
		if kth == nth {
			return curr // 1
		}
		return fibRecur(curr, pred+curr, kth+1)
	}

	return fibRecur(initialOffset, curr, kth)
}
