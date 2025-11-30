package point

func trap(height []int) int {
	stack := []int{}
	var res int
	for i := 0; i < len(height); i++ {
		for len(stack) > 0 && height[i] >= height[stack[len(stack)-1]] {
			h := height[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				break
			}
			left := stack[len(stack)-1]
			w := i - left - 1
			dh := min(height[i], height[left]) - h
			res += dh * w
		}
		stack = append(stack, i)
	}
	return res
}
