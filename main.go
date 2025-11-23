package main

import (
	"context"
	"fmt"
	"gocode_windows/config"
	"gocode_windows/distributedlock"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

// distributed lock main function
func main() {
	// 加载配置
	_, err := config.LoadConfig("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 创建一个 Redis 后端的锁
	lock := distributedlock.NewDistributedLockInfo(
		"my-resource",  // 锁的键
		"unique-value", // 客户端的唯一标识
		30*time.Second, // 锁的超时时间
	)

	// 设置重试选项
	lock.SetRetry(3, 100*time.Millisecond)

	// 获取锁
	ctx := context.Background()
	acquired, err := lock.AcquireLock(ctx, "redis")
	if err != nil {
		panic(fmt.Sprintf("获取锁失败: %v", err))
	}

	if !acquired {
		fmt.Println("获取锁失败")
		return
	}

	// 确保在函数退出时释放锁
	defer func() {
		err := lock.ReleaseLock(ctx, "redis")
		if err != nil {
			fmt.Printf("释放锁失败: %v\n", err)
		}
	}()

	// 临界区代码
	fmt.Println("成功获取锁，执行任务中...")
	time.Sleep(5 * time.Second)
	fmt.Println("任务完成")
}
