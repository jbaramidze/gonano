package main

func minOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func insertInSlice(arr []rune, p rune, i int) []rune {
	arr1 := append(arr, 0)
	copy(arr1[i+1:], arr1[i:])
	arr1[i] = p
	return arr1
}

func removeFromSlice(arr []rune, i int) []rune {
	n := make([]rune, len(arr)-1)
	copy(n, arr[:i])
	copy(n[i:], arr[i+1:])
	return n
}
