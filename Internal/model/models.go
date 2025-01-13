package model

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"test_component/Internal/settings"
)

func InitRedis(config settings.RedisConfig) (rdb *redis.ClusterClient, err error) {

	fmt.Println(config.Addrs) //
	//fmt.Println(reflect.TypeOf(config.Addrs)) //[]string
	//addrs := strings.Split(config.Addrs, ",")
	//fmt.Println(addrs)
	rdb = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:          config.Addrs,
			Password:       config.Password,
			RouteByLatency: config.RouteByLatency,
			DialTimeout:    config.DialTimeout,
			ReadTimeout:    config.ReadTimeout,
			WriteTimeout:   config.WriteTimeout,
		})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("获取rdb错误")
		log.Println(err)
	}
	//---------------------------------------------------------------------
	//err = rdb.Set(ctx, "key1", "value2", 0).Err()
	//if err != nil {
	//	log.Fatalf("could not set key: %v", err)
	//}
	//
	//val, err := rdb.Get(ctx, "key1").Result()
	//if err != nil {
	//	log.Fatalf("could not get key: %v", err)
	//}
	//
	//fmt.Println("key:", val)
	//rdb.Close()
	//------------------------------------------------------------------------
	return
}
func InitPostgres(config settings.PgConfig) (db *gorm.DB, err error) {
	dsn := config.DataSource
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	//-----------------------------------------------------------------------------
	//err = db.AutoMigrate(&User{}, &Order{})
	//if err != nil {
	//	log.Fatalf("自动迁移失败: %v", err)
	//}
	//fmt.Println("模型自动迁移完成")
	//
	//user := User{Name: "Tony", Age: 24}
	//result := db.Create(&user)
	//if result.Error != nil {
	//	log.Fatalf("无法插入数据: %v", result.Error)
	//}
	//fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
	//
	//user = User{Name: "Mike", Age: 14}
	//result = db.Create(&user)
	//if result.Error != nil {
	//	log.Fatalf("无法插入数据: %v", result.Error)
	//}
	//fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
	//
	//user = User{Name: "Jhon", Age: 32}
	//result = db.Create(&user)
	//if result.Error != nil {
	//	log.Fatalf("无法插入数据: %v", result.Error)
	//}
	//fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
	//
	//var user1 User
	//if err := db.Where("name = ?", "Tony").Find(&user1).Error; err != nil {
	//	log.Fatalf("查询失败: %v", err)
	//}
	//fmt.Printf("查询到的用户: ID: %d, Name: %s, Age: %d\n", user1.ID, user1.Name, user1.Age)
	//
	//// 使用原生 SQL 查询
	//var user2 User
	//db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&user2)
	//fmt.Println("使用原生 SQL 查询年龄大于 25 的用户:", user2)
	//------------------------------------------------------------------------
	return
}
func InitDoris(config settings.DorisConfig) (db *gorm.DB, err error) {
	dsn := config.DataSource
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	//-----------------------------------------------------
	//CREATE TABLE users (
	//	id BIGINT NOT NULL AUTO_INCREMENT,
	//	name VARCHAR(255) NOT NULL,
	//	email VARCHAR(255) NOT NULL,
	//	age INT NOT NULL,
	//	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	//	PRIMARY KEY (id)
	//)
	//DISTRIBUTED BY HASH(id) BUCKETS 3
	//PROPERTIES (
	//	"replication_num" = "1"
	//);
	err = db.AutoMigrate(&UserDoris{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("模型自动迁移完成")

	user := User{Name: "Tony", Age: 24}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("无法插入数据: %v", result.Error)
	}
	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)

	user = User{Name: "Mike", Age: 14}
	result = db.Create(&user)
	if result.Error != nil {
		log.Fatalf("无法插入数据: %v", result.Error)
	}
	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)

	user = User{Name: "Jhon", Age: 32}
	result = db.Create(&user)
	if result.Error != nil {
		log.Fatalf("无法插入数据: %v", result.Error)
	}
	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)

	var user1 User
	if err := db.Where("name = ?", "Tony").Find(&user1).Error; err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	fmt.Printf("查询到的用户: ID: %d, Name: %s, Age: %d\n", user1.ID, user1.Name, user1.Age)

	// 使用原生 SQL 查询
	var user2 User
	db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&user2)
	fmt.Println("使用原生 SQL 查询年龄大于 25 的用户:", user2)
	//------------------------------------------------------
	return
}
