package goblazer

// BinarySearchIntsEQ searches the 'idx' of any element in slice 's' where s[idx] == v. The slice 's' must be sorted in ascending order.
func BinarySearchIntsEQ(s []int, v int) (idx int) {
	idx = -1
	low, high := 0, len(s)-1

	for low <= high {
		mid := (low + high) >> 1
		if s[mid] < v { // 这里说明s[low]到s[mid]之间所有数都比v小
			low = mid + 1
		} else if s[mid] > v { // 这里说明s[mid]到s[high]之间所有数都比v大
			high = mid - 1
		} else {
			idx = mid
			break
		}
	}
	return
}

// BinarySearchIntsGE searches the 'idx' of first element in slice 's' where s[idx] >= v. The slice 's' must be sorted in ascending order.
func BinarySearchIntsGE(s []int, v int) (idx int) {
	idx = -1
	num := len(s)
	low, high := 0, num

	for low < high {
		mid := (low + high) >> 1 // mid ≤ mid < high
		if s[mid] < v {          // 这里说明s[low]到s[mid]之间所有数都比v小
			low = mid + 1
		} else { // 这里说明s[mid]到s[high]之间所有数都比v大，且s[mid]是第一个比v大的数
			high = mid
		}
	}

	if low != num {
		idx = low
	}
	return
}

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
