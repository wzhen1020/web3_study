package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	intChan := make(chan int, 100)

	wg.Add(1)
	go producer(intChan)

	wg.Add(1)
	go consumer(intChan)
	wg.Wait()
	fmt.Println("执行完毕...")
}

func producer(intChan chan int) {
	defer wg.Done()
	for i := 1; i <= 100; i++ {
		intChan <- i
		fmt.Printf("生产值: %d\n", i)
	}
	close(intChan)
}
func consumer(intChan chan int) {

	defer wg.Done()
	for v := range intChan {
		fmt.Printf("取到值：%v", v)
	}

}
