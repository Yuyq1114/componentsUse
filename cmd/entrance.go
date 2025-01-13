package cmd

import (
	"context"
	"fmt"
	"test_component/Internal/engine/task1"
	"test_component/Internal/model"
	"test_component/Internal/settings"
)

type Arg struct {
	ConfigMap string
}

func Entrance(arg Arg) (err error) {
	// 初始化执行文件
	config, err := settings.InitConfig(arg.ConfigMap)
	//初始化日志
	//主流程
	runTask(&config)

	return nil
}

func runTask(config *settings.Config) {
	ctx := context.Background()

	//初始化各种组件

	//如redis
	fmt.Println(config.Redis)
	redisDB, _ := model.InitRedis(ctx, config.Redis)

	// 初始化主任务成员
	tak, _ := task1.New(ctx, *redisDB)
	tak.Start()
	//成员的运行方法

}
