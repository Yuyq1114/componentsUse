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
	Ctx          context.Context
	TestString   string
	Rdb          *redis.Client
	Pg           *gorm.DB
	DorisDb      *gorm.DB
	Kw           *kafka.Writer
	Kr           *kafka.Reader
	km           *kafka.Message
	NamingClient naming_client.INamingClient
}

func New(ctx context.Context, rdb *redis.Client, pg *gorm.DB, dorisDb *gorm.DB, kafkawrite *kafka.Writer, kafkaReader *kafka.Reader, naminfClient naming_client.INamingClient) (*Task1, error) {
	task1 := &Task1{
		Ctx:        ctx,
		TestString: "hello，this is the task1 start",
		Rdb:        rdb,
		Pg:         pg,
		DorisDb:    dorisDb,
		Kw:         kafkawrite,
		Kr:         kafkaReader,
		//km:           kafkaMessage,
		NamingClient: naminfClient,
	}
	return task1, nil
}

//func (task1 *Task1) Run() error {
//	go task1.ProduceTestMessage(task1.Ctx)
//	go task1.Start(c)
//	return nil
//}

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
			value := time.Now().String()
			err := task1.Kw.WriteMessages(task1.Ctx, kafka.Message{
				Key:   []byte("key"),
				Value: []byte(value),
			})
			if err != nil {
				return err
			}
		}

	}

}

func (task1 *Task1) Start(ch chan kafka.Message) {
	//task1.KafkaConn.WriteMessages("hello")
	fmt.Println("mission start")
	for {
		select {
		case <-task1.Ctx.Done():
			log.Println("task1的消息consume协程结束")
		default:
			fmt.Println("start consume")
			m, err := task1.Kr.ReadMessage(task1.Ctx)

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

			}

		}
	}

	fmt.Println(task1.TestString)
}

//nacos的服务调用，微服务架构中类似注册gin.Get(
/*func (task1 *Task1) NacosMission() {
	model.RegisterService(task1.NamingClient, "add", "test_component", "117.50.85.130", 8848)

	port, u := model.GetIPAndPort(task1.NamingClient, "add", "test_component")

	serviceURL := fmt.Sprintf("http://%s:%d", port, u)

	data := map[string]string{
		"num1": "1",
		"num2": "2",
	}
	jsonData, err := json.Marshal(data)

	resp, err := http.Post(serviceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		log.Println("调用成功")
	} else {
		log.Println("调用失败")
	}
}
*/
