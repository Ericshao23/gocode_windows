package main

// 暴力求解
func isValid(s string) bool {
	if len(s)%2 != 0 || len(s) < 2 {
		return false
	} else {
		list := make([]uint8, 0)
		list = append(list, s[0])

		for i := 1; i < len(s); i++ {
			if len(list) == 0 {
				list = append(list, s[i])
			} else {

				switch list[len(list)-1] {
				case '(':
					if s[i] == ')' {
						list = list[:len(list)-1]
					} else {
						list = append(list, s[i])
					}

				case '{':
					if s[i] == '}' {
						list = list[:len(list)-1]
					} else {
						list = append(list, s[i])
					}
				case '[':
					if s[i] == ']' {
						list = list[:len(list)-1]
					} else {
						list = append(list, s[i])
					}
				default:
					return false
				}
			}
		}
		if len(list) == 0 {
			return true
		} else {
			return false
		}
	}
}

// 用key的k-v维护switch-case
func isValid2(s string) bool {
	length := len(s)
	if length == 1 {
		return false
	}
	pairs := map[byte]byte{
		')': '(',
		'}': '{',
		']': '[',
	}
	stack := []byte{}
	for i := 0; i < length; i++ {
		if pairs[s[i]] > 0 {
			if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}
