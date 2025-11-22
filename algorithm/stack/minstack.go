package main

import (
	"math"
)

type MinStack struct {
	min      int
	minindex int
	stack    []int
	len      int
}

const maxInt32 = 9223372036854775807

func Constructor() MinStack {
	return MinStack{
		min:      maxInt32,
		minindex: -1,
		stack:    []int{},
		len:      0,
	}
}

func (this *MinStack) Push(val int) {
	this.stack = append(this.stack, val)
	this.len++
	if this.min > val {
		this.min = val
		this.minindex = this.len - 1
	}
}

func (this *MinStack) Pop() {
	if this.len > 0 {
		// 先移除栈顶元素
		this.stack = this.stack[:this.len-1]
		this.len--

		// 重新计算最小值
		if this.len == 0 {
			this.min = maxInt32
			this.minindex = -1
		} else {
			this.min = this.stack[0]
			this.minindex = 0
			for i := 1; i < this.len; i++ {
				if this.stack[i] < this.min {
					this.min = this.stack[i]
					this.minindex = i
				}
			}
		}
	}
}

func (this *MinStack) Top() int {
	return this.stack[this.len-1]
}

func (this *MinStack) GetMin() int {
	if this.len == 0 {
		return 0
	}
	return this.min
}

//添加辅助栈，用于保存相应位置的最小值
/*
stack：[-2,1,3,-4,0]
minStack:[-2,-2,-2,-4,-4]
*/
type MinStackV2 struct {
	stack    []int
	minStack []int
}

func ConstructorV2() MinStackV2 {
	return MinStackV2{
		stack:    []int{},
		minStack: []int{math.MaxInt64},
	}
}

func (this *MinStackV2) Push(x int) {
	this.stack = append(this.stack, x)
	top := this.minStack[len(this.minStack)-1]
	this.minStack = append(this.minStack, min(x, top))
}

func (this *MinStackV2) Pop() {
	this.stack = this.stack[:len(this.stack)-1]
	this.minStack = this.minStack[:len(this.minStack)-1]
}

func (this *MinStackV2) Top() int {
	return this.stack[len(this.stack)-1]
}

func (this *MinStackV2) GetMin() int {
	return this.minStack[len(this.minStack)-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// 栈中除了保存添加的元素，还保存前缀最小值。（栈中保存的是 pair）
// 添加元素：设当前栈的大小是 n。添加元素 val 后，额外维护 preMin[n]=min(preMin[n−1],val)，其中 preMin[n−1] 是添加 val 之前，栈顶保存的前缀最小值。
type pair struct{ val, preMin int }

type MinStackv3 []pair

func Constructorv3() MinStackv3 {
	// 这里的 0 写成任意数都可以，反正用不到
	return MinStackv3{{0, math.MaxInt}} // 栈底哨兵
}

func (st *MinStackv3) Push(val int) {
	*st = append(*st, pair{val, min(st.GetMin(), val)})
}

func (st *MinStackv3) Pop() {
	*st = (*st)[:len(*st)-1]
}

func (st MinStackv3) Top() int {
	return st[len(st)-1].val
}

func (st MinStackv3) GetMin() int {
	return st[len(st)-1].preMin
}
