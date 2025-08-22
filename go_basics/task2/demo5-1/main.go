package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup

func main() {
	var counter int64
	for i := 0; i < 10; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done() // 协程结束时减少等待组计数

			// 每个协程执行1000次递增操作
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}

			fmt.Printf("协程 %d 完成\n", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("counter:", counter)
}
