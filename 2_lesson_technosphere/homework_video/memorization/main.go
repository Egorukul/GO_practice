package main

import "fmt"

type memoizeFunction func(int, ...int) interface{}

// TODO реализовать
var fibonacci memoizeFunction
var romanForDecimal memoizeFunction

func fibonacciReal(r int, v ...int) interface{} {
	sl := make([]int, r)
	for i := range r {
		if i < 2 {
			sl[i] = i
		} else {
			sl[i] = sl[i-1] + sl[i-2]
		}
	}
	return sl[r-1]
}

//TODO Write memoization function

func memoize(function memoizeFunction) memoizeFunction {
	cache := make(map[int]int)

	wrapper := func(r int, v ...int) interface{} {
		if result, ok := cache[r]; ok {
			return result
		} else {
			result := function(r, v...).(int)
			cache[r] = result
			return result
		}

	}

	return wrapper

}

// TODO обернуть функции fibonacci и roman в memoize
func init() {

}

func main() {

	fibonacci = fibonacciReal
	fmt.Println("Fibonacci(45) =", fibonacci(45).(int))
	// for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
	// 	14, 15, 16, 17, 18, 19, 20, 25, 30, 40, 50, 60, 69, 70, 80,
	// 	90, 99, 100, 200, 300, 400, 500, 600, 666, 700, 800, 900,
	// 	1000, 1009, 1444, 1666, 1945, 1997, 1999, 2000, 2008, 2010,
	// 	2012, 2500, 3000, 3999} {
	// 	fmt.Printf("%4d = %s\n", x, romanForDecimal(x).(string))
	// }
}
