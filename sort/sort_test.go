package sort

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{[]int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{[]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, testCase := range testCases {
		quickSort(testCase.input, 0, len(testCase.input)-1)
		if !reflect.DeepEqual(testCase.input, testCase.expected) {
			t.Errorf("期望得到 %v, 但实际得到 %v", testCase.expected, testCase.input)
		}
	}
}

func TestMergeSort(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{[]int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{[]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, testCase := range testCases {
		sortedArr := mergeSort(testCase.input)
		if !reflect.DeepEqual(sortedArr, testCase.expected) {
			t.Errorf("期望得到 %v，但实际得到 %v", testCase.expected, sortedArr)
		}
	}
}
