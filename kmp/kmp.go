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
	result := []int{}
	textLength, targetLength := len(text), len(target)
	if targetLength > textLength || targetLength == 0 {
		return result
	}
	textIndex, targetIndex, targetMatchIndex := 0, 0, 0
	partialMatchTable := buildPartialMatchTable(target)

	for textIndex < textLength {
		if text[textIndex] == target[targetIndex] {
			textIndex++
			targetIndex++
			targetMatchIndex++
			if targetMatchIndex == targetLength {
				shift := targetMatchIndex - partialMatchTable[targetMatchIndex-1]
				result = append(result, textIndex-targetMatchIndex)
				targetIndex -= shift
				targetMatchIndex -= shift
			}
		} else {
			if targetMatchIndex == 0 {
				textIndex++
			} else {
				shift := targetMatchIndex - partialMatchTable[targetMatchIndex-1]
				targetIndex -= shift
				targetMatchIndex -= shift
			}
		}
	}

	return result
}
