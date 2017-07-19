package kmp

func buildPartialMatchTable(input string) []int {
	result := []int{0}
	i, j := 0, 1

	for j < len(input) {
		if input[j] == input[i] {
			result = append(result, i+1)
			j++
			i++
		} else {
			if i > 0 {
				i = result[i-1]
			} else {
				result = append(result, 0)
				j++
			}
		}
	}
	return result
}

// KMP - A fast substring index finder with O(m + n)
func KMP(text, target string) []int {
	// TODO: complete KMP function
	result := []int{}
	return result
}
