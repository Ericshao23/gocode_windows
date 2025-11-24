package main

import "sort"

// hashmap
func longestConsecutive(nums []int) int {
	hmap := make(map[int]bool)
	for _, num := range nums {
		hmap[num] = true
	}
	length := len(nums)
	if length == 0 {
		return 0
	}
	maxLen := 0
	for num := range hmap {
		if hmap[num-1] {
			continue
		}
		currentNum := num
		for hmap[currentNum+1] {
			currentNum++
		}
		maxLen = max(maxLen, currentNum+1-num)
	}
	return maxLen
}

func longestConsecutiveV2(nums []int) int {
	hmap := make(map[int]bool)
	for _, num := range nums {
		hmap[num] = true
	}
	length := len(nums)
	if length == 0 {
		return 0
	}
	maxLen := 0
	for num := range hmap {
		if hmap[num-1] {
			continue
		}
		currentNum := num
		for hmap[currentNum+1] {
			currentNum++
		}
		maxLen = max(maxLen, currentNum+1-num)
	}
	return maxLen
}

// 排序
func longestConsecutiveV1(nums []int) int {
	sort.Ints(nums)
	var count, result, last int
	for i := 0; i < len(nums); i++ {
		if i == 0 || nums[i] != last+1 {
			last = nums[i]
			count++
			result = max(result, count)
		} else if nums[i] == nums[i-1] {
			continue
		} else {
			last = nums[i]
			count = 1
		}
	}
	return result
}
