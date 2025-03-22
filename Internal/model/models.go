package model

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"test_component/Internal/settings"
)

func InitRedis(config settings.RedisConfig) (rdb *redis.Client, err error) {

	//fmt.Println(config.Addrs) //
	//fmt.Println(reflect.TypeOf(config.Addrs)) //[]string
	//addrs := strings.Split(config.Addrs, ",")
	//fmt.Println(addrs)
	//rdb = redis.NewClusterClient(
	//	&redis.ClusterOptions{
	//		Addrs:          config.Addrs,
	//		Password:       config.Password,
	//		RouteByLatency: config.RouteByLatency,
	//		DialTimeout:    config.DialTimeout,
	//		ReadTimeout:    config.ReadTimeout,
	//		WriteTimeout:   config.WriteTimeout,
	//	})
	//_, err = rdb.Ping(context.Background()).Result()
	//if err != nil {
	//	fmt.Println("获取rdb错误")
	//	log.Println(err)
	//}
	//--------------------------------------------------------------------
	rdb = redis.NewClient(&redis.Options{
		Addr:         config.Addrs,    // Redis 服务器地址
		Password:     config.Password, // 没有密码，则
		DB:           0,               // 默认数据库
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
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
	//fmt.Println(dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	return
}
func InitDoris(config settings.DorisConfig) (db *gorm.DB, err error) {
	dsn := config.DataSource
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	//-----------------------------------------------------
	//	CREATE TABLE IF NOT EXISTS mydb.example_tbl
	//	(
	//		timestamp DATE NOT NULL COMMENT "['0000-01-01', '9999-12-31']",
	//	type TINYINT NOT NULL COMMENT "[-128, 127]",
	//		error_code INT COMMENT "[-2147483648, 2147483647]",
	//		error_msg VARCHAR(300) COMMENT "[1-65533]",
	//
	//)
	//	DUPLICATE KEY(timestamp, type)
	//	DISTRIBUTED BY HASH(type) BUCKETS 1
	//	PROPERTIES (
	//		"replication_allocation" = "tag.location.default: 1"
	//	);
	//	err = db.AutoMigrate(&UserDoris{})
	//	if err != nil {
	//		log.Fatalf("自动迁移失败: %v", err)
	//	}
	//	fmt.Println("模型自动迁移完成")
	//
	//	user := User{Name: "Tony", Age: 24}
	//	result := db.Create(&user)
	//	if result.Error != nil {
	//		log.Fatalf("无法插入数据: %v", result.Error)
	//	}
	//	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
	//
	//	user = User{Name: "Mike", Age: 14}
	//	result = db.Create(&user)
	//	if result.Error != nil {
	//		log.Fatalf("无法插入数据: %v", result.Error)
	//	}
	//	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
	//
	//	user = User{Name: "Jhon", Age: 32}
	//	result = db.Create(&user)
	//	if result.Error != nil {
	//		log.Fatalf("无法插入数据: %v", result.Error)
	//	}
	//	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
	//
	//	var user1 User
	//	if err := db.Where("name = ?", "Tony").Find(&user1).Error; err != nil {
	//		log.Fatalf("查询失败: %v", err)
	//	}
	//	fmt.Printf("查询到的用户: ID: %d, Name: %s, Age: %d\n", user1.ID, user1.Name, user1.Age)
	//
	//	// 使用原生 SQL 查询
	//	var user2 User
	//	db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&user2)
	//	fmt.Println("使用原生 SQL 查询年龄大于 25 的用户:", user2)
	//------------------------------------------------------
	return
}
func InitNacos(config settings.NacosConfig) (namingClient naming_client.INamingClient, err error) {
	//配置服务端
	ServerConfig := []constant.ServerConfig{
		{
			IpAddr: config.ServerAddr,
			Port:   config.ServerPort,
		},
	}

	// 配置 Nacos 客户端
	ClientConfig := constant.ClientConfig{
		NamespaceId: config.NamespaceId,
		Username:    config.Username,
		Password:    config.Password,
	}
	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: ServerConfig,
			ClientConfig:  &ClientConfig,
		},
	)
	if err != nil {
		//logger.Error(err)
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("初始化nacos成功")
	// 创建 Nacos 服务发现客户端
	return namingClient, nil
}

//func InitKafka(config settings.KafkaConfig) (conn *kafka.Conn, error error) {
//	conn, error = kafka.Dial(config.Protocol, config.Addr)
//	return
//}

func KafkaProcureIns(topic string, conf settings.KafkaConfig) *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(conf.Addr),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}
	fmt.Println("初始化kafka生产者成功")
	return w
	//w.WriteMessages(ctx, message...)
	//defer w.Close()
}

func KafkaConsumeIns(topic string, conf settings.KafkaConfig) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{conf.Addr},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	//r.SetOffset(42)
	fmt.Println("初始化kafka消费者成功")
	return r
	//for {
	//	m, err := r.ReadMessage(ctx)
	//	if err != nil {
	//		break
	//	}
	//	fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	//}
	//
	//if err := r.Close(); err != nil {
	//	log.Fatal("failed to close reader:", err)
	//}

}
