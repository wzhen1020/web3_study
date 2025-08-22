package main

import "fmt"

//  Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
}

func (rectangle *Rectangle) Area() {
	fmt.Println("Rectangle--Area")

}

func (rectangle *Rectangle) Perimeter() {
	fmt.Println("Rectangle--Perimeter")

}

type Circle struct {
}

func (circle *Circle) Area() {
	fmt.Println("Circle--Area")

}

func (circle *Circle) Perimeter() {
	fmt.Println("Circle--Perimeter")

}

func main() {

	var circle = &Circle{}
	circle.Area()
	circle.Perimeter()

	var rectangle = &Rectangle{}
	rectangle.Area()
	rectangle.Perimeter()

}
