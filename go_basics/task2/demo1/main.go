package main

import "fmt"

func main() {

	num := 10
	add(&num)
	fmt.Println(num)
	slice := []int{4, 5, 6}
	multiply(&slice)

	fmt.Println(slice)

}

// 加10
func add(num *int) {
	*num += 10
}

// 切片乘10
func multiply(nums *[]int) {
	for i := 0; i < len(*nums); i++ {
		(*nums)[i] *= 10

	}
}
