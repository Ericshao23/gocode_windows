package main

import "fmt"

func dailyTemperatures(T []int) []int {
	res := make([]int, len(T))
	res[len(T)-1] = 0
	for i := len(T) - 1; i >= 0; i-- {
		count := 0
		for j := i - 1; j >= 0; j-- {
			if T[j] < T[i] {
				res[j] = count + 1
				break
			} else {
				count++
			}
		}
	}
	return res
}

func main() {
	fmt.Println(dailyTemperatures([]int{73, 74, 75, 71, 69, 72, 76, 73}))
}
