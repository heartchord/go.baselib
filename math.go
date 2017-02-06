package goblazer

// Max is
func Max(max int, args ...int) int {
	for _, v := range args {
		if max < v {
			max = v
		}
	}
	return max
}

// Min is
func Min(min int, args ...int) int {
	for _, v := range args {
		if v < min {
			min = v
		}
	}

	return min
}
