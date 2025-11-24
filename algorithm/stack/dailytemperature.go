package main

import (
	"slices"
)

// 暴力
func dailyTemperatures(T []int) []int {
	stack := []int{}
	res := make([]int, len(T))
	for i := len(T) - 1; i >= 0; i-- {
		if len(stack) == 0 {
			res[i] = 0
			// stack = append(stack, T[i])
		} else {
			count := 0
			for j := len(stack) - 1; j >= 0; j-- {
				if T[i] < stack[j] {
					count++
					res[i] = count
					break
				} else if j == 0 {
					count = 0
				} else {
					count++
				}
			}
			res[i] = count
		}
		stack = append(stack, T[i])
	}
	return res
}

// 维护一个单调栈（单调增或减），栈顶至栈底单调递减
// 当元素破坏递减性时，统一结算
func dailyTemperaturesv2(T []int) []int {
	stack := []int{}
	res := make([]int, len(T))
	for i := 0; i < len(T); i++ {
		for len(stack) > 0 && T[i] > T[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res[top] = i - top
		}
		stack = append(stack, i)
	}
	return res
}

// 从右至左单调栈 单调增（栈顶最小）
func dailyTemperaturesv3(T []int) []int {
	res := make([]int, len(T))
	stack := []int{}
	for i, t := range slices.Backward(T) {
		for len(stack) > 0 && t >= T[stack[len(stack)-1]] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			res[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}
	return res
}
