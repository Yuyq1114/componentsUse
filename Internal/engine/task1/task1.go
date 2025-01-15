package task1

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"log"
	"time"
)

type Task1 struct {
	Ctx        context.Context
	TestString string
	Rdb        *redis.ClusterClient
	Pg         *gorm.DB
	//DorisDb    *gorm.DB
	Kw           *kafka.Writer
	Kr           *kafka.Reader
	NamingClient naming_client.INamingClient
}

func New(ctx context.Context, rdb *redis.ClusterClient, pg *gorm.DB, kafkawrite *kafka.Writer, kafkaReader *kafka.Reader, naminfClient naming_client.INamingClient) (*Task1, error) {
	task1 := &Task1{
		Ctx:        ctx,
		TestString: "hello，this is the task1 start",
		Rdb:        rdb,
		Pg:         pg,
		//DorisDb:    dorisDb,
		Kw:           kafkawrite,
		Kr:           kafkaReader,
		NamingClient: naminfClient,
	}
	return task1, nil
}

// ProduceTestMessage 生产测试消息
func (task1 *Task1) ProduceTestMessage(ctx context.Context) error {
	log.Println("task1的消息生产协程start")
	for {
		select {
		case <-ctx.Done(): //检查ctx.Done()是否被关闭
			log.Println("task1的消息生产协程结束")
			return nil
		default:
			time.Sleep(3 * time.Second)
			value, _ := time.Now().MarshalBinary()
			err := task1.Kw.WriteMessages(task1.Ctx, kafka.Message{
				Key:   []byte("key"),
				Value: value,
			})
			if err != nil {
				return err
			}
		}

	}

}

func (task1 *Task1) Start() {
	//task1.KafkaConn.WriteMessages("hello")
	fmt.Println("mission start")
	for {
		select {
		case <-task1.Ctx.Done():
			log.Println("task1的消息consume协程结束")
		default:
			m, err := task1.Kr.ReadMessage(task1.Ctx)
			fmt.Println("start consume")
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

			}

		}
	}

	fmt.Println(task1.TestString)
}
