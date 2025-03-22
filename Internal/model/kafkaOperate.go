package model

import (
	"context"
	"github.com/segmentio/kafka-go"
	"test_component/Internal/settings"
)

func ConsumeFromCertainTopic(ctx context.Context, config *settings.KafkaConfig, topic string, ch chan *kafka.Message) (err error) {
	//初始化kafka
	return nil
}
