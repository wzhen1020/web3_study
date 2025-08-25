package employees

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employees struct {
	Id         uint
	Name       string
	Department string
	Salary     float64
}

/*
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

func Query1() {
	dsn := "root:root@tcp(127.0.0.1:3306)/devs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var employeeSlice []Employees

	err = db.Select(&employeeSlice, "select * from employees where department = ?", "技术部")

	if err == nil {

		for _, v := range employeeSlice {
			fmt.Printf("ID: %d, Name: %s, Department: %s, Salary: %.2f\n",
				v.Id, v.Name, v.Department, v.Salary)
		}

	}
}

func Query2() {
	dsn := "root:root@tcp(127.0.0.1:3306)/devs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var employeeSlice []Employees

	err = db.Select(&employeeSlice, "select * from employees order by salary desc limit 1")

	if err == nil {

		for _, v := range employeeSlice {
			fmt.Printf("ID: %d, Name: %s, Department: %s, Salary: %.2f\n",
				v.Id, v.Name, v.Department, v.Salary)
		}

	}
}
