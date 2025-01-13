package model

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"test_component/Internal/settings"
)

func InitRedis(ctx context.Context, config settings.RedisConfig) (rdb *redis.ClusterClient, err error) {

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
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("获取rdb错误")
		log.Println(err)
	}
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
	return
}
func InitPostgres(ctx context.Context, config settings.PgConfig) (db *gorm.DB, err error) {
	dsn := config.DataSource
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	return
}
