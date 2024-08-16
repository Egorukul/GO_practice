package main

import (
	"fmt"
)

// TODO: Реализовать вычисление Квадратного корня
func Sqrt(x float64) float64 {

	if x < 0 {
		return x
	}

	var assumption_value float64
	var last_assumption_value float64
	assumption_value = 10
	
	for {
		assumption_value = float64(0.5) * (assumption_value + (x / float64(assumption_value)))
		eps := assumption_value - last_assumption_value
		if eps < 0 {
			eps = -eps
		}
		if eps < 0.01 {
			break
		}
		last_assumption_value = assumption_value

	}

	return assumption_value
}

func main() {
	fmt.Println(Sqrt(4))
}
