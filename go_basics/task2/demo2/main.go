package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	// num := 10
	// wg.Add(1)
	// go odd(num)
	// wg.Add(1)
	// go even(num)
	// wg.Wait()

	// fmt.Println("执行完毕。。。。")

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go task(i)
	}

	wg.Wait()
	fmt.Println("执行完毕。。。。")
}

func task(num int) {
	startTime := time.Now().Unix()

	for i := ((num - 1) * 1000) - 1; i < num*1000; i++ {
		if i%2 != 0 {
			fmt.Println("打印奇数：", i)
		}
	}

	endTime := time.Now().Unix()
	fmt.Println("任务：", num, "执行时间：", endTime-startTime)
	wg.Done()
}

func odd(num int) {

	for i := 1; i <= num; i++ {
		if i%2 != 0 {
			fmt.Println("打印奇数：", i)
		}
	}
	wg.Done()

}

func even(num int) {

	for i := 2; i <= num; i++ {
		if i%2 == 0 {
			fmt.Println("打印偶数:", i)
		}
	}
	wg.Done()

}
