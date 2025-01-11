package main

import "test_component/cmd"

func main() {

	cmd.Execute()
	//// 创建一个 context
	//ctx := context.Background()
	//
	//// 配置 Redis 集群的节点
	//rdb := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: []string{
	//		"118.178.127.89:7001", // Redis 集群节点 1
	//		"118.178.127.89:7002", // Redis 集群节点 2
	//		"118.178.127.89:7003", // Redis 集群节点 3
	//		"118.178.127.89:7004",
	//		"118.178.127.89:7005",
	//		"118.178.127.89:7006",
	//		// 可以添加更多节点
	//	},
	//	Password:       "mypassword", // 如果需要密码，则提供密码
	//	RouteByLatency: true,         //根据延迟选择节点
	//	//DialTimeout:    10 * time.Second, // 设置连接超时时间
	//	//ReadTimeout:    10 * time.Second, // 设置读超时时间
	//})
	//
	//// 测试连接是否正常
	//_, err := rdb.Ping(ctx).Result()
	//if err != nil {
	//	log.Fatalf("could not connect to Redis: %v", err)
	//}
	//
	//fmt.Println("Connected to Redis cluster!")
	//
	//// 执行简单的命令
	//err = rdb.Set(ctx, "key", "value2", 0).Err()
	//if err != nil {
	//	log.Fatalf("could not set key: %v", err)
	//}
	//
	//val, err := rdb.Get(ctx, "key").Result()
	//if err != nil {
	//	log.Fatalf("could not get key: %v", err)
	//}
	//
	//fmt.Println("key:", val)
	//
	//// 关闭连接
	//err = rdb.Close()
	//if err != nil {
	//	log.Fatalf("could not close Redis connection: %v", err)
	//}
}
