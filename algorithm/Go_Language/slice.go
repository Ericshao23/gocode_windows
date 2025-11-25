package golanguage

import (
	"fmt"
)

func sliceBasic() {
	/* 切片的定义方式*/
	// 内置函数len返回有效长度，cap函数返回切片容量大小，len <= cap
	var (
		a []int               //nil 切片,和nil相等
		b = []int{}           //空切片，长度为0，不等于nil
		c = []int{1, 2, 3}    //含有初始值的切片, 有三个元素，len和cap均为3
		d = c[:2]             //切片d引用了切片c的前两个元素，len为2，cap为3
		e = c[0:2:cap(c)]     //切片e引用了切片c的前两个元素，len为2，cap为3
		f = c[:0]             //长度为0，cap为3
		g = make([]int, 3)    //长度为3，cap为3，初始值为0
		h = make([]int, 2, 3) //长度为2，cap为3，初始值为0
		i = make([]int, 0, 3) //长度为0，cap为3
	)

	/*切片遍历*/
	for i := range a {
		fmt.Printf("a[%d]=%d\n", i, a[i])
	}
	for i, v := range b {
		fmt.Printf("b[%d]=%d\n", i, v)
	}
	for i := 0; i < len(c); i++ {
		fmt.Printf("c[%d]=%d\n", i, c[i])
	}

	/*添加切片*/
	a = append(a, 1)                    //向nil切片添加单个元素
	a = append(a, 2, 3, 4)              //向切片添加多个元素，手写解包方式
	a = append(a, []int{5, 6}...)       //向切片添加另一个切片，切片解包方式
	a = append([]int{0}, a...)          //在切片开头添加元素，需要重新分配底层数组
	a = append([]int{-3, -2, -1}, a...) //在切片开头添加多个元素，需要重新分配底层数组
	// 通过链式操作，给切片中间添加元素
	j, x := 5, 99
	a = append(a[:j], append([]int{x}, a[j:]...)...)       //在切片中间位置j添加元素x
	a = append(a[:j], append([]int{1, 2, 3}, a[j:]...)...) //在切片中间位置j添加切片

	// copy函数，将一个切片复制到另一个切片中
	a = append(a, 0)     // 切片向后扩展一个元素位置
	copy(a[j+1:], a[j:]) // a[j:]部分向后移动一位
	a[j] = x             // 在位置j插入元素x
	// copy和append实现中间位置插入多个元素
	a = append(a, c...)       // 为c切片扩充足够的空间
	copy(a[j+len(c):], a[j:]) // 将a[j:]部分向后移动len(c)位
	copy(a[j:], c)            // 将c切片复制到a[j:]位置

	/*删除切片*/
	// 从尾部删除
	a = []int{1, 2, 3}
	a = a[:len(a)-1] // 删除最后一个元素
	a = a[:len(a)-j] // 删除最后j个元素
	// 从头部删除
	a = []int{1, 2, 3}
	a = a[1:] // 删除第一个元素
	a = a[j:] // 删除前j个元素
	/*原地址实现删除*/
	// 基于append实现删除开头元素
	a = []int{1, 2, 3, 4, 5}
	a = append(a[:0], a[1:]...) // 删除第一个元素
	a = append(a[:0], a[j:]...) // 删除前j个元素
	// 基于copy实现删除开头元素
	a = []int{1, 2, 3, 4, 5}
	a = a[:copy(a, a[1:])] // 删除第一个元素
	a = a[:copy(a, a[j:])] // 删除前j个元素
	// 删除中间元素
	a = []int{1, 2, 3, 4, 5}
	a = append(a[:2], a[3:]...)    // 删除索引2的元素
	a = append(a[:j], a[j+x:]...)  // 删除索引j开始的x个元素
	a = a[:j+copy(a[j:], a[j+x:])] // 删除索引j开始的x个元素
}
