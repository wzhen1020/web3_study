package transaction

import (
	"errors"
	"fmt"
	"task3/datebeas"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

/*
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账
*/
type Accounts struct {
	Id      uint            `gorm:"primarykey"`
	Balance decimal.Decimal `gorm:"type:decimal(22,2)"`
}

type Transactions struct {
	Id            uint `gorm:"primarykey"`
	FromAccountId uint
	ToAccountId   uint
	Amount        decimal.Decimal `gorm:"type:decimal(22,2)"`
}

var db = datebeas.DB

// 创建表
func CreateTable() {
	db.AutoMigrate(&Accounts{})
	db.AutoMigrate(&Transactions{})
}

func Insert(accounts *Accounts) {

	db.Create(accounts)
}

func Transfer(fromAccountId uint, toAccountId uint, amount decimal.Decimal) bool {

	// tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	err := db.Transaction(func(tx *gorm.DB) error {

		var fromAccount Accounts

		// 校验转出账户余额
		tx.Debug().Where("id = ?", fromAccountId).Find(&fromAccount)

		if fromAccount.Balance.LessThan(amount) {
			return errors.New("余额不足！")
		}

		// 转出
		fromResult := tx.Debug().Model(&Accounts{}).Where("id = ?", fromAccountId).Update("balance", fromAccount.Balance.Sub(amount))
		if fromResult.Error != nil {
			// fmt.Errorf(fromResult.Error.Error())
			return errors.New("转出失败")
		}
		fmt.Println(amount)
		// 转入
		toResult := tx.Debug().Model(&Accounts{}).Where("id = ?", toAccountId).Update("balance", gorm.Expr("balance+?", amount))

		if toResult.Error != nil {
			// fmt.Errorf(toResult.Error.Error())
			return errors.New("转入失败")
		}

		// 记录
		result := tx.Debug().Create(&Transactions{FromAccountId: fromAccountId, ToAccountId: toAccountId, Amount: amount})
		if result.Error != nil {
			// fmt.Errorf(toResult.Error.Error())
			return errors.New("转账记录失败")
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
