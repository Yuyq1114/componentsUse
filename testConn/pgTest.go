package testConn

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User 定义用户模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Orders    []Order        `gorm:"foreignKey:UserID"`
}

// Order 定义订单模型
type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Product   string
	Price     float64
	CreatedAt time.Time
}

func initPgDB() *gorm.DB {
	// 替换为你的PostgreSQL连接信息
	dsn := "host=localhost user=postgres password=password dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	return db
}

// 自动迁移模型
func autoMigratePg(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Order{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("模型自动迁移完成")
}

// 插入数据
func createUserPg(db *gorm.DB, name string, age int) {
	user := User{Name: name, Age: age}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("无法插入数据: %v", result.Error)
	}
	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
}

// 查询单个用户
func findUserPg(db *gorm.DB, id uint) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Fatalf("查询用户失败: %v", result.Error)
	}
	fmt.Printf("查询到的用户: ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
}

// 查询多个用户
func findUsersPg(db *gorm.DB) {
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
func updateUserPg(db *gorm.DB, id uint, newAge int) {
	result := db.Model(&User{}).Where("id = ?", id).Update("age", newAge)
	if result.Error != nil {
		log.Fatalf("更新用户失败: %v", result.Error)
	}
	fmt.Printf("成功更新用户 ID %d 的年龄为 %d\n", id, newAge)
}

// 删除用户
func deleteUserPg(db *gorm.DB, id uint) {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		log.Fatalf("删除用户失败: %v", result.Error)
	}
	fmt.Printf("成功删除用户 ID %d\n", id)
}

// 创建订单
func createOrderPg(db *gorm.DB, userID uint, product string, price float64) {
	order := Order{UserID: userID, Product: product, Price: price}
	result := db.Create(&order)
	if result.Error != nil {
		log.Fatalf("创建订单失败: %v", result.Error)
	}
	fmt.Printf("成功创建订单: 用户ID: %d, 产品: %s, 价格: %.2f\n", userID, product, price)
}

// 查询用户的订单
func findUserOrdersPg(db *gorm.DB, userID uint) {
	var orders []Order
	result := db.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		log.Fatalf("查询订单失败: %v", result.Error)
	}
	fmt.Printf("用户 ID %d 的订单:\n", userID)
	for _, order := range orders {
		fmt.Printf("订单ID: %d, 产品: %s, 价格: %.2f\n", order.ID, order.Product, order.Price)
	}
}

// 事务处理
func transactionalCreatePg(db *gorm.DB, userName string, userAge int, product string, price float64) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := User{Name: userName, Age: userAge}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 创建订单
		order := Order{UserID: user.ID, Product: product, Price: price}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("事务操作失败: %v", err)
	}
	fmt.Println("事务操作成功")
}

// 使用 GORM 的高级查询方法
func advancedQueriesPg(db *gorm.DB) {
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
