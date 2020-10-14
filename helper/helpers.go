package helper

func MinOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func InsertInSlice(arr []rune, p rune, i int) []rune {
	arr1 := append(arr, 0)
	copy(arr1[i+1:], arr1[i:])
	arr1[i] = p
	return arr1
}
