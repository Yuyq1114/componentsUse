package testConn

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

const redisAddr = "120.26.84.229:6379"

func Init() *redis.Client {
	ctx := context.Context(context.Background())
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Redis 服务器地址
		Password: "",        // 没有密码，默认
		DB:       0,         // 默认数据库
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("无法连接到 Redis: %v", err)
	}
	fmt.Printf("连接成功: %v\n", pong)
	return rdb
}

func Pub(rdb *redis.Client, ctx context.Context) {
	// 发布消息
	err := rdb.Publish(ctx, "mychannel", "Hello, Redis!").Err()
	if err != nil {
		log.Fatalf("无法发布消息: %v", err)
	}
}

func Sub(rdb *redis.Client, ctx context.Context) {
	// 发布/订阅（Pub/Sub）模式
	pubsub := rdb.Subscribe(ctx, "mychannel")

	// 订阅消息
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Printf("接收到消息: %s\n", msg.Payload)
		}
	}()
}
