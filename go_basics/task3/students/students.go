package students

import "task3/datebeas"

var db = datebeas.DB

type Students struct {
	Id    uint
	Name  string
	Age   int
	Grade string
}

func Add() {
	db.Create(&Students{1, "里斯", 18, "一年级"})
}

// 初始化表
func InitTable() {
	db.AutoMigrate(&Students{})
}
