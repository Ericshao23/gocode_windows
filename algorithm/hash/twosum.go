package main

// 基于hashmap 两次循环实现
func twoSum(nums []int, target int) []int {
	hmap := make(map[int]int)
	for i, num := range nums {
		hmap[num] = i
	}
	for i, num := range nums {
		index := target - num
		if j, ok := hmap[index]; ok && i != j {
			return []int{i, j}
		}
	}
	return []int{}
}

// 基于hashmap 1次循环实现
func twoSumV2(nums []int, target int) []int {
	hmap := make(map[int]int)
	for i, num := range nums {
		index := target - num
		if j, ok := hmap[index]; ok && j != i {
			return []int{i, j}
		}
		hmap[num] = i
	}
	return []int{}
}
