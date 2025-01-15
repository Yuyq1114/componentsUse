package cmd

import (
	"context"
	"fmt"
	"test_component/Internal/engine/task1"
	"test_component/Internal/model"
	"test_component/Internal/settings"
	"time"
)

type Arg struct {
	ConfigMap string
}

func Entrance(arg Arg) (err error) {
	fmt.Println("entrance")
	// 初始化执行文件
	config, err := settings.InitConfig(arg.ConfigMap)
	//初始化日志
	//主流程
	runTask(&config)

	return nil
}

func runTask(config *settings.Config) {
	ctx := context.Background()
	//wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//初始化各种组件

	//如redis
	//fmt.Println(config.Redis)
	redisDB, _ := model.InitRedis(config.Redis)
	defer redisDB.Close()
	//初始化pg
	fmt.Println("pg的初始化配置信息为", config.Pg)
	pg, _ := model.InitPostgres(config.Pg)

	//初始化Doris,内存不够，暂不使用
	//dorisDb, _ := model.InitDoris(config.Doris)

	//初始化kafka
	procureTask1 := model.KafkaProcureIns("task1", config.Kafka)
	defer procureTask1.Close()

	consumeTask1 := model.KafkaConsumeIns("task1", config.Kafka)
	defer consumeTask1.Close()

	//初始化nacos
	namingClient, _ := model.InitNacos(config.Nacos)
	defer namingClient.CloseClient()
	// 初始化主任务成员

	//启动一些必要的协程

	//wg.Add(1)
	tak, _ := task1.New(ctx, redisDB, pg, procureTask1, consumeTask1, namingClient)

	go func() error {
		return tak.ProduceTestMessage(ctx)
	}()

	tak.Start()
	//成员的运行方法

	//wg.Done()
	time.Sleep(20)
	fmt.Println("任务结束")
}
