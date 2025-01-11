package testConn

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义一个模型
type UserDoris struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func initDorisDB() *gorm.DB {
	// 替换为你的Doris连接信息
	dsn := "username:password@tcp(127.0.0.1:9030)/your_db_name?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	return db
}

// 自动迁移模型
func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("模型自动迁移完成")
}

// 插入数据
func createUser(db *gorm.DB, name string, age int) {
	user := User{Name: name, Age: age}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("无法插入数据: %v", result.Error)
	}
	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
}

// 查询单个用户
func findUser(db *gorm.DB, id uint) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Fatalf("查询用户失败: %v", result.Error)
	}
	fmt.Printf("查询到的用户: ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
}

// 查询多个用户
func findUsers(db *gorm.DB) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatalf("查询用户列表失败: %v", result.Error)
	}
	fmt.Println("用户列表:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

// 更新用户
func updateUser(db *gorm.DB, id uint, newAge int) {
	result := db.Model(&User{}).Where("id = ?", id).Update("age", newAge)
	if result.Error != nil {
		log.Fatalf("更新用户失败: %v", result.Error)
	}
	fmt.Printf("成功更新用户 ID %d 的年龄为 %d\n", id, newAge)
}

// 删除用户
func deleteUser(db *gorm.DB, id uint) {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		log.Fatalf("删除用户失败: %v", result.Error)
	}
	fmt.Printf("成功删除用户 ID %d\n", id)
}

// 使用 GORM 的高级查询方法
func advancedQueries(db *gorm.DB) {
	var users []User

	// 获取所有用户并按年龄降序排序
	db.Order("age desc").Find(&users)
	fmt.Println("按年龄降序排序的用户列表:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// 查询年龄大于 25 的用户
	db.Where("age > ?", 25).Find(&users)
	fmt.Println("年龄大于 25 的用户:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// 使用原生 SQL 查询
	db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&users)
	fmt.Println("使用原生 SQL 查询年龄大于 25 的用户:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}
