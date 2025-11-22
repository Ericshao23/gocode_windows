package main

import (
	"strconv"
	"strings"
)

// 递归
func decodeString(s string) string {
	stack := []string{}
	currentNum := 0
	currentStr := ""
	for _, char := range s {
		switch {
		case char >= '0' && char <= '9':
			currentNum = currentNum*10 + int(char-'0')
		case char == '[':
			stack = append(stack, currentStr)
			stack = append(stack, strconv.Itoa(currentNum))
			currentNum = 0
			currentStr = ""
		case char == ']':
			num, _ := strconv.Atoi(stack[len(stack)-1])
			stack = stack[:len(stack)-1]
			prevStr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			currentStr = prevStr + strings.Repeat(currentStr, num)
		default:
			currentStr += string(char)
		}
	}
	return currentStr
}
