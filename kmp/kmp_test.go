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

func TestKMP(t *testing.T) {
	testCases := []struct {
		text   string
		target string
		output []int
	}{
		{"a", "a", []int{0}},
		{"ab", "a", []int{0}},
		{"aba", "a", []int{0, 2}},
		{"abcdabca", "dab", []int{3}},
		{"aabaabaaa", "cd", []int{}},
		{"abcdaby", "ab", []int{0, 4}},
		{"ABC ABCDAB ABCDABCDABDE", "ABCDABD", []int{15}},
		{"co", "cococp", []int{}},
		{"cococp", "co", []int{0, 2}},
	}

	for _, testCase := range testCases {
		actual := KMP(testCase.text, testCase.target)
		expected := testCase.output
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("KMP(%s, %s): expected: %v, actual: %v", testCase.text, testCase.target, expected, actual)
		}
	}
}
