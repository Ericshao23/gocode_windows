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

func decodeStringV2(s string) string {
	stack := []string{}
	CurrentNum := 0
	currentStr := ""
	for i := 0; i < len(s); i++ {
		char := s[i]
		if char >= '0' && char <= '9' {
			CurrentNum = CurrentNum*10 + int(char-'0')
		} else if char == '[' {
			stack = append(stack, currentStr)
			stack = append(stack, strconv.Itoa(CurrentNum))
			CurrentNum = 0
			currentStr = ""
		} else if char == ']' {
			num, _ := strconv.Atoi(stack[len(stack)-1])
			stack = stack[:len(stack)-1]
			prevStr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			currentStr = prevStr + strings.Repeat(currentStr, num)
		} else {
			currentStr += string(char)
		}

	}
	return currentStr
}

// ] 之前的字符逐个入栈
// 然后逐个出栈，直到遇到[ 这部分出栈内容就是重复子串
// 继续后边遇到的数字就是子串重复的次数，直到遇到非数字
// 将重复子串按次数展开后入栈
func decodeStringV3(s string) string {
	stack := []string{}
	for _, char := range s {
		if char != ']' {
			stack = append(stack, string(char))
			continue
		}
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		subStr := ""
		for cur != "[" {
			subStr = cur + subStr
			cur = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		curNum := ""
		for len(stack) > 0 {
			cur = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if cur < "0" || cur > "9" {
				break
			}
			curNum = cur + curNum
		}
		if cur < "0" || cur > "9" {
			stack = append(stack, cur)
		}
		repeatNum, _ := strconv.Atoi(curNum)
		stack = append(stack, strings.Repeat(subStr, repeatNum))
	}
	return strings.Join(stack, "")
}
