package main

import "fmt"

//
// Person 结构体，包含 Name 和 Age 字段，
// 再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	EmployeeId string
	person     Person
}

func (employee *Employee) PrintInfo() {
	fmt.Println(employee)
}

func main() {
	var employee Employee
	employee.EmployeeId = "1"
	employee.person.Name = "golang"
	employee.person.Age = 20
	employee.PrintInfo()
}
