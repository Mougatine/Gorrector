package trie

func sliceMin(arr []int) int {
	min := 100

	for _, val := range arr {
		if val < min {
			min = val
		}
	}
	return min
}
