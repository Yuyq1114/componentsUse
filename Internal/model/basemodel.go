package model

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Product   string
	Price     float64
	CreatedAt time.Time
}

type UserDoris struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 使用grom创建表时会自动的创建某个表的名字，这个可以自定义表名
func (u *User) TableName() string {
	return "users"
}

// InsertUser 插入用户表
func InsertUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
func (u *User) InsertUser(db *gorm.DB) error {
	return db.Create(u).Error
}

// SelectUserByAge 查找年龄大于指定年龄的用户
func (u *User) SelectUserByAge(db *gorm.DB, age int) (userSelect *User, err error) {
	err = db.Raw("SELECT * FROM users WHERE age > ?", age).Scan(&userSelect).Error
	return
}

// SelectUserByName 从my_user table查找姓名为name的用户
func (u *User) SelectUserByName(db *gorm.DB, name string) (user1 []*User, err error) {
	if err := db.Where("name = ?", name).Find(&user1).Error; err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	return
}

//err = db.AutoMigrate(&User{}, &Order{})
//if err != nil {
//	log.Fatalf("自动迁移失败: %v", err)
//}
//fmt.Println("模型自动迁移完成")
