package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var mu sync.Mutex

func main() {
	counter := 0
	for i := 0; i < 10; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done() // 协程结束时减少等待组计数

			// 每个协程执行1000次递增操作
			for j := 0; j < 1000; j++ {
				mu.Lock()   // 获取锁
				counter++   // 递增计数器
				mu.Unlock() // 释放锁
			}

			fmt.Printf("协程 %d 完成\n", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("counter:", counter)
}
