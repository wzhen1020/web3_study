package datebeas

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {

	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/devs?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	DB = db

}
