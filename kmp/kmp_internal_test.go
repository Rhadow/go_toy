package kmp

import (
	"reflect"
	"testing"
)

func TestPartialMatchTable(t *testing.T) {
	testCases := []struct {
		input  string
		output []int
	}{
		{"a", []int{0}},
		{"ab", []int{0, 0}},
		{"aba", []int{0, 0, 1}},
		{"abcdabca", []int{0, 0, 0, 0, 1, 2, 3, 1}},
		{"aabaabaaa", []int{0, 1, 0, 1, 2, 3, 4, 5, 2}},
		{"abcaby", []int{0, 0, 0, 1, 2, 0}},
	}

	for _, testCase := range testCases {
		actual := buildPartialMatchTable(testCase.input)
		expected := testCase.output
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("buildPartialMatchTable(%s): expected: %v, actual: %v", testCase.input, expected, actual)
		}
	}
}
