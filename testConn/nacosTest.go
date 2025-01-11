package testConn

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

func InitNacos() {
	// 配置 Nacos 客户端
	clientConfig := constant.NewClientConfig(
		constant.WithTimeoutMs(5000),
		constant.WithNamespaceId("public"), // 默认为 public
		constant.WithAccessKey("your-access-key"),
		constant.WithSecretKey("your-secret-key"),
	)
	// 创建 Nacos 服务发现客户端
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}
	// 获取服务列表
	resp, err := client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		PageSize: 10,
		PageNo:   1,
	})
	
	if err != nil {
		log.Fatalf("获取服务列表失败: %v", err)
	}

	// 输出服务列表
	fmt.Println("服务列表：", resp)
}
