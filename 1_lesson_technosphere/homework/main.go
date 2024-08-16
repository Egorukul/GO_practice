package main

import (
	"fmt"
	"sort"
)

func ReturnInt() int {
	return 1
}

func ReturnFloat() float32 {
	return 1.1
}

func ReturnIntArray() [3]int {
	return [...]int{1, 3, 4}
}

func ReturnIntSlice() []int {
	return []int{1, 2, 3}
}

func IntSliceToString(sl []int) string {
	var result string
	for _, val := range sl {
		result += fmt.Sprintf("%d", val)
	}
	return result
}

func MergeSlices(fl_sl []float32, int_sl []int32) []int {
	result := make([]int, 0, len(fl_sl)+len(int_sl))
	for _, val := range fl_sl {
		result = append(result, int(val))
	}
	for _, val := range int_sl {
		result = append(result, int(val))
	}
	return result
}

func GetMapValuesSortedByKey(sl map[int]string) []string {
	result := make([]string, 0, len(sl))
	temp_sl := make([]int, 0, len(sl))
	for key := range sl {
		temp_sl = append(temp_sl, key)
	}

	sort.Slice(temp_sl, func(i, j int) bool {
		return temp_sl[i] < temp_sl[j]
	})

	for _, val := range temp_sl {
		result = append(result, sl[val])
	}
	return result
}
func main() {

}
