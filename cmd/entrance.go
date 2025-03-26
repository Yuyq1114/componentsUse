package cmd

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"test_component/Internal/engine/access"
	"test_component/Internal/engine/asset"
	"test_component/Internal/model"
	"test_component/Internal/settings"
)

var chanCap = 100

// Arg is the argument for Entrance
type Arg struct {
	ConfigMap string
}

func Entrance(arg Arg) (err error) {
	fmt.Println("entrance")
	// 初始化执行文件
	config, err := settings.InitConfig(arg.ConfigMap)
	//初始化日志
	//主流程
	err = runTask(&config)
	//fmt.Println("任务结束in Entrance")
	if err != nil {
		fmt.Println("任务结束 in Entrance")
	}
	return

}

func runTask(config *settings.Config) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//通过errgroup优雅退出
	eg, egCtx := errgroup.WithContext(ctx)
	//wg := sync.WaitGroup{}
	//ctx, cancel := context.WithCancel(ctx)

	//初始化各种组件

	//如redis
	//fmt.Println(config.Redis)
	redisDB, err := model.InitRedis(config.Redis)
	if err != nil {
		fmt.Println("初始化redis失败")
	} else {
		fmt.Println("初始化redis成功")
	}
	defer redisDB.Close()
	//初始化pg
	//fmt.Println("pg的初始化配置信息为", config.Pg)
	pg, err := model.InitPostgres(config.Pg)
	if err != nil {
		fmt.Println("初始化pg失败")
	} else {
		fmt.Println("初始化pg成功")
	}
	//初始化Doris,内存不够，暂不使用
	dorisDb, err := model.InitDoris(config.Doris)
	if err != nil {
		fmt.Println("初始化doris失败")
	} else {
		fmt.Println("初始化doris成功")
	}

	//初始化kafka
	//procureTask1 := model.KafkaProcureIns("task1", config.Kafka)
	//defer procureTask1.Close()
	//
	//consumeTask1 := model.KafkaConsumeIns("task1", config.Kafka)
	//defer consumeTask1.Close()

	//初始化nacos
	//namingClient, _ := model.InitNacos(config.Nacos)
	//defer namingClient.CloseClient()

	//初始化所有的通道
	accessOutputCh := make(chan *model.TLV, chanCap)

	// 初始化主任务成员
	//  初始化接入
	acc := access.New(egCtx, eg, accessOutputCh, config.Kafka, dorisDb, pg)
	acc.Start()
	//启动一些必要的协程

	//初始化tag
	//tag := discover.New(egCtx, eg, accessOutputCh, pg, redisDB)
	//tag.Start()

	//初始化asset
	asset := asset.New(egCtx, eg, accessOutputCh, pg, redisDB, dorisDb)
	asset.Start()

	//初始化运行统计模块
	//stati := stat.New(egCtx, eg)
	//stati.Start()

	//wg.Add(1)
	//tak, _ := task1.New(ctx, redisDB, pg, dorisDb, procureTask1, consumeTask1, namingClient)
	//tak.Start()
	//任务间的通信channel再次初始化

	//任务需要的协程在次初始化

	//go func() error {
	//	return tak.ProduceTestMessage(ctx)
	//}()
	//
	//tak.Start()
	//成员的运行方法

	err = eg.Wait()
	if err != nil {
		fmt.Println("任务结束 in runTask")
	}
	return err
	//time.Sleep(20)
}
