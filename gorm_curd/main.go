package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type User struct {
	gorm.Model
	Name     string
	Age      int
	Birthday time.Time
	Email    string `gorm:"type:varchar(100);unique_index"`
	Address  string `gorm:"index:addr"` // 给Address 创建一个名字是  `addr`的索引
}

func (User) TableName() string {
	return "test"
}

// BeforeCreate 在执行SQL语句前，可以做一些统一的处理，比如将 user.id 设置为 uuid
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	fmt.Println("beforeCreate user....")
	//scope.SetColumn("ID", uuid.New())
	return nil
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

	// 根据 User struct 创建表
	db.AutoMigrate(&User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

	// 创建数据
	//create(db)

	// 查询数据
	//query(db)

	// 原生查询
	//protozoaQuery(db)

	// 删除数据
	//delete(db)
}

func delete(db *gorm.DB) {
	db.Where("id = ?", "3").Delete(User{})
}

func protozoaQuery(db *gorm.DB) {
	var user User
	db.Where("name = ?", "lishi").First(&user)
	fmt.Printf("%v\n", user)
}

func query(db *gorm.DB) {
	var user User

	// 获取第一条记录，按主键排序
	db.First(&user)
	//// SELECT * FROM users ORDER BY id LIMIT 1;
	fmt.Printf("first: %v\n", user)

	// 获取一条记录，不指定排序
	var user2 User
	db.Take(&user2)
	//// SELECT * FROM users LIMIT 1;
	fmt.Printf("Take: %v\n", user2)

	// 获取最后一条记录，按主键排序
	var user3 User
	db.Last(&user3)
	//// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	fmt.Printf("Last: %v\n", user3)

	// 获取所有的记录
	var users []User
	db.Find(&users)
	//// SELECT * FROM users;
	fmt.Printf("Find all:\n")
	for _, u := range users {
		fmt.Printf("%v\n", u)
	}

	// 通过主键进行查询 (仅适用于主键是数字类型)
	var user4 User
	db.First(&user4, 2)
	//// SELECT * FROM users WHERE id = 10;
	fmt.Printf("Last: %v\n", user4)
}

func create(db *gorm.DB) {
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	result := db.Create(&user)
	if result.Error != nil {
		fmt.Printf("err:%v", result.Error)
	} else {
		fmt.Printf("success, id = %d ,RowsAffected = %d", user.ID, result.RowsAffected)
	}
}
