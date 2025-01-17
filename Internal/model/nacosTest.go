package model

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

func RegisterService(client naming_client.INamingClient, serviceName string, groupName string, ip string, port uint64) {
	success, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		Healthy:     true,
		Enable:      true,
		Weight:      1,
		Ephemeral:   true,
	})
	if err != nil {
		log.Println(err)
	}
	if success {
		log.Printf("Service %s registered successfully at %s:%d", serviceName, ip, port)
	}
}

func GetIPAndPort(client naming_client.INamingClient, serviceName string, groupName string) (string, uint64) {
	instance, err := client.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName,
	})
	if err != nil {
		panic(err)
	}
	return instance.Ip, instance.Port
}
