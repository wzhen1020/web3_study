package books

import (
	"fmt"
	"task3/datebeas"

	"github.com/jmoiron/sqlx"
)

type Book struct {
	Id     uint
	Title  string
	Author string
	Price  float64
}

var db = datebeas.DB

func Query() {
	dsn := "root:root@tcp(127.0.0.1:3306)/devs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var books []Book
	query := `SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC`

	err = db.Select(&books, query, 50)

	if err == nil {

		for i, book := range books {
			fmt.Printf("%d. ID: %d, 书名: 《%s》, 作者: %s, 价格: ￥%.2f\n",
				i+1, book.Id, book.Title, book.Author, book.Price)
		}

	}
}
