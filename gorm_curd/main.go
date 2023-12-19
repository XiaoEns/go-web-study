package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type User struct {
	gorm.Model
	Name         string
	Age          sql.NullInt64
	Birthday     *time.Time
	Email        string  `gorm:"type:varchar(100);unique_index"`
	Role         string  `gorm:"size:255"`        //设置字段的大小为255个字节
	MemberNumber *string `gorm:"unique;not null"` // 设置 memberNumber 字段唯一且不为空
	Num          int     `gorm:"AUTO_INCREMENT"`  // 设置 Num字段自增
	Address      string  `gorm:"index:addr"`      // 给Address 创建一个名字是  `addr`的索引
	IgnoreMe     int     `gorm:"-"`               //忽略这个字段
}

const (
	server   = "127.0.0.1"
	port     = "3306"
	user     = "root"
	password = "123456"
	database = "demo"
)

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, server, port, database)
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

}
