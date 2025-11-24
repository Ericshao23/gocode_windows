package main

import "sort"

// 排序
func groupAnagrams(strs []string) [][]string {
	anagramMap := make(map[string][]string)
	var result [][]string
	for _, str := range strs {
		strBytes := []byte(str)
		sort.Slice(strBytes, func(i, j int) bool {
			return strBytes[i] < strBytes[j]
		})
		sortedStr := string(strBytes)
		anagramMap[sortedStr] = append(anagramMap[sortedStr], str)
	}
	for _, group := range anagramMap {
		result = append(result, group)
	}
	return result
}

// 计数，字母出现次数作为key
// 字母异位词的特点是字母种类和出现次数都相同
func groupAnagramsV2(strs []string) [][]string {
	mp := make(map[[26]int][]string)
	for _, str := range strs {
		cnt := [26]int{}
		for _, ch := range str {
			cnt[ch-'a']++
		}
		mp[cnt] = append(mp[cnt], str)
	}
	res := make([][]string, 0, len(mp))
	for _, group := range mp {
		res = append(res, group)
	}
	return res
}
