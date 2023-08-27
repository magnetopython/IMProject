package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartbeatTime uint64
	LogOutTime    uint64
	IsLoginout    bool
	DeviceInfo    string
}

func main() {
	db, err := gorm.Open(mysql.Open("root:722722@(localhost)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&UserBasic{})

	// Create
	// db.Create(&UserBasic{Code: "D42", Price: 100})
	user := &UserBasic{}
	user.Name = "aaa"
	db.Create(user)

	// Read
	fmt.Printf("db.First(user, 1): %v\n", db.First(user, 1)) // 根据整型主键查找
	// db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(user).Update("Password", "1234")
	// Update - 更新多个字段
	// db.Model(&product).Updates(UserBasic{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	// db.Delete(&product, 1)
}
