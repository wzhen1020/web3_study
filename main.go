package main

import (
	"fmt"
	"sort"
)

func main() {
	// nums := []int{2, 2, 1}

	// result := singleNumber(nums)
	// fmt.Println(result)

	// fmt.Println(isPalindrome(121))
	// s := "(])"

	// fmt.Println(isValid(s))
	// strs := []string{
	// 	"ab",
	// 	"a",
	// }

	// fmt.Println(longestCommonPrefix(strs))

	// digits := []int{
	// 	9,
	// }

	// fmt.Println(plusOne(digits))

	nums := []int{
		0, 0, 1, 1, 1, 2, 2, 3, 3, 4,
	}

	// fmt.Println(removeDuplicates(nums))

	// intervals := [][]int{
	// 	{1, 4}, {0, 1},
	// }
	// fmt.Println(merge(intervals))

	fmt.Println(twoSum(nums, 7))

}

func twoSum(nums []int, target int) []int {

	// map key=值 value = 下标
	numMap := make(map[int]int)

	// 迭代切片
	for i, num := range nums {

		// 目标值减当前值 得到的结果在map中是否存在
		if val, ok := numMap[target-num]; ok {
			// 存在 则返回当前值下标及结果值下标
			return []int{val, i}
		}

		// 不存在 将当前存放map
		numMap[num] = i

	}
	return make([]int, 0)

}

// 合并区间
func merge(intervals [][]int) [][]int {

	if len(intervals) == 0 {
		return intervals
	}

	// 按区间的起始值排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := make([][]int, 0)
	for _, interval := range intervals {
		// 如果 merged 为空或当前区间不与最后一个区间重叠
		if len(merged) == 0 || interval[0] > merged[len(merged)-1][1] {
			merged = append(merged, interval)
		} else {
			// 重叠时合并区间，更新结束值为最大值
			if interval[1] > merged[len(merged)-1][1] {
				merged[len(merged)-1][1] = interval[1]
			}
		}
	}
	return merged
}

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {

	// 新切片
	newNums := make([]int, 0)
	// 循环切片
	for i := 0; i < len(nums); {
		// 最后一个直接放入新切片
		if i == len(nums)-1 {
			return len(append(newNums, nums[i]))
		}

		// 当前下标值 与 下一个不相等 就将当前下标值放入新切片
		if nums[i] != nums[i+1] {
			newNums = append(newNums, nums[i])
			i++
		} else {
			// 否则就移除当前下标值
			nums = append(nums[:i], nums[i+1:]...)
		}
	}
	return len(newNums)
}

// 加一
func plusOne(digits []int) []int {

	// var num int = 1

	// 从后面开始循环
	for i := len(digits) - 1; i >= 0; i-- {

		// 值+1
		digits[i]++
		// 等于10 进一补0
		if 10 == digits[i] {
			digits[i] = 0
			// 不等于直接返回
		} else {
			return digits
		}
	}
	// 全部循环完 进一补位
	return append([]int{1}, digits...)

	// size := len(digits)

	// for i := 0; i < size-1; i++ {
	// 	for j := 0; j < size-i-1; j++ {
	// 		temp := digits[j]

	// 		// fmt.Println(i)
	// 		if temp < digits[j+1] {
	// 			digits[j] = digits[j+1]
	// 			digits[j+1] = temp
	// 		}

	// 	}
	// }

	// return digits

}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {

	// rune 切片
	prdfixStr := make([]rune, 0)
	// 第一个字符串
	str := strs[0]

	// flag := true
	// 循环第一个字符串字符
	strRune := []rune(str)
	for i := 0; i < len(strRune); i++ {
		// 循环字符串
		for j := 1; j < len(strs); j++ {

			// 当前长度是否超过 字符串长度
			if i >= len(strs[j]) {
				return string(prdfixStr)
			}
			// 判断字符是否相等
			if strRune[i] != []rune(strs[j])[i] {
				return string(prdfixStr)
			}
		}
		// 相等的字符塞到切片中
		prdfixStr = append(prdfixStr, strRune[i])

	}

	return str

}

// 有效的括号
func isValid(s string) bool {

	// 声明切片 模拟栈
	stack := make([]rune, 0)

	// 初始化map的键值关系 左括号为key 右括号为vlaue
	strMap := map[rune]rune{
		'(': ')',
		'{': '}',
		'[': ']',
	}

	// 迭代字符串
	for _, ch := range s {
		// 是否为左括号
		match, isLeft := strMap[ch]
		if isLeft {
			// 为左括号时将 对应的右括号塞进去 最后出现的左括号 所对应的右括号 就会在 切片的末尾
			stack = append(stack, match)
			fmt.Printf("stack--%c\n", stack)
		} else {
			size := len(stack)
			// 当前不是左括号 切片中也没有右括号 直接返回false
			if size == 0 {
				return false
			}
			// 获取切片末尾的下标
			index := size - 1
			// 如果切片末尾的右括号 与当前字符串循环到右括号一致
			if ch == stack[index] {
				// 说明括号闭环  弹出该右括号
				stack = stack[:index]
			} else {
				// 否则返回false
				return false
			}

		}
		// fmt.Printf("%c--%v\n", matching, isRight)

	}
	// fmt.Printf("%c\n", stack)
	// 切片长度为0 说明所有
	return len(stack) == 0

}

// 回文数
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	if x%10 == 0 && x != 0 {
		return false
	}

	var temp int

	fmt.Println(temp)
	for x > temp {

		temp = temp*10 + x%10
		x /= 10
	}

	fmt.Println(x, temp)
	return x == temp || x == temp/10

}

// 136. 只出现一次的数字
func singleNumber(nums []int) int {
	var maps = make(map[int]int)
	for i := 0; i < len(nums); i++ {

		maps[nums[i]]++

	}
	// result := make([]int, 0)
	for k, v := range maps {
		if v == 1 {
			return k
		}
	}
	return -1

}
