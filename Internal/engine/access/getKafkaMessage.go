package access

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
	"test_component/Internal/model"
	"test_component/Internal/settings"
	"test_component/Internal/utils"
	"time"
	"unicode/utf8"
)

type Access struct {
	Ctx         context.Context
	Eg          *errgroup.Group
	OutPutChan  chan *model.TLV
	KafkaConfig settings.KafkaConfig
	DorisConfig *gorm.DB
	Pg          *gorm.DB
}

var (
	buffer       []model.TLV // 缓存数据
	bufferMutex  sync.Mutex  // 互斥锁，防止并发问题
	flushTicker  *time.Ticker
	flushTimeout = 3 * time.Minute // 3分钟超时
	maxBufferLen = 5               // 每50条数据批量写入一次
)

const (
	mitmproxyTopic = "mitmproxy-topic"
)

// init 初始化的时候自动运行，没3分钟写一次到数据库
func init() {
	flushTicker = time.NewTicker(flushTimeout)

	// 启动后台定时写入任务
	go func() {
		for range flushTicker.C {
			flushToDatabase()
		}
	}()
}

func New(Ctx context.Context, eg *errgroup.Group, AccessOutputCh chan *model.TLV, KafkaRea settings.KafkaConfig, Doris *gorm.DB, Pg *gorm.DB) *Access {
	//OutPutChan := make(chan *model.TLV)
	return &Access{
		Ctx:         Ctx,
		Eg:          eg,
		OutPutChan:  AccessOutputCh,
		KafkaConfig: KafkaRea,
		DorisConfig: Doris,
		Pg:          Pg,
	}
}
func (access *Access) Start() {
	access.Eg.Go(func() error {
		return access.GetKafkaMessage()
	})
}

func (access *Access) GetKafkaMessage() error {
	//雪花算法生成ID并加入
	machine, err := utils.NewSnowflake(1)
	if err != nil {
		fmt.Println("❌ 雪花算法生成ID失败:", err)
		return err
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{access.KafkaConfig.Addr}, // 修改为你的 Kafka 地址
		Topic:   mitmproxyTopic,
		GroupID: access.KafkaConfig.GroupId,
	})

	fmt.Println("🟢 正在监听 Kafka 数据...")

	for {
		select {
		case <-access.Ctx.Done():
			fmt.Println("access  mission end")
			return nil
		default:
			for {
				msg, err := r.ReadMessage(context.Background())
				if err != nil {
					log.Fatalf("❌ 读取消息失败: %v", err)
				}

				//fmt.Printf("✅ 收到消息: %s\n", msg.Value)
				//将消息解析成结构体
				var tlv model.TLV
				err = json.Unmarshal(msg.Value, &tlv)
				if err != nil {
					log.Printf("❌ JSON 解析失败: %v, 消息内容: %s", err, msg.Value)
					return nil // 直接丢弃或放入死信队列
				}

				//生成唯一ID
				tlv.LogId = machine.GenerateID()
				//jsonData, err := json.Marshal(tlv)
				//jsonData, err := sonic.Marshal(tlv)
				//if err != nil {
				//	fmt.Println("❌ 解析 JSON 失败:", err)
				//	return err
				//}
				//if tlv.HttpsData.Content == "\"\"" {
				//	fmt.Println("Content 为空")
				//
				//}
				tlv.HttpsData.Content = processContent(tlv.HttpsData.Content)
				//jsonData, err := ConvertToJSON(tlv)
				//if err != nil {
				//	log.Fatal("JSON 转换失败:", err)
				//}

				// 打印解析结果
				fmt.Println("✅ 解析成功:")
				//fmt.Printf("Type: %s\n", tlv.Type)
				fmt.Printf("URL: %s\n", tlv.HttpsData.URL)
				//-----------------------------------下面的两个地方会不会出现数据不一致？先写数据库 后走通道就不会了--------------------------------
				//批量写入数据库
				//放入缓存中  每50条写入一次，或者每3分钟写入一次
				//err = model.StreamInsertData(&jsonData, access.DorisConf, "https_logs")
				//if err != nil {
				//	fmt.Println("数据插入失败")
				//} else {
				//	fmt.Println("数据插入成功")
				//}
				// 放入缓存
				bufferMutex.Lock()
				buffer = append(buffer, tlv)
				//fmt.Println(len(buffer))
				if len(buffer) >= maxBufferLen {
					go flushToDatabase()
				}
				bufferMutex.Unlock()

				access.OutPutChan <- &tlv

			}
		}

	}
}

// flushToDatabase  将上面的存到doris改为每有50条数据或者3分钟存到数据库中，减小读取压力
func flushToDatabase() {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	if len(buffer) == 0 {
		return
	}

	// 转换为 JSON
	jsonData, err := ConvertToJSON(buffer)
	if err != nil {
		log.Fatal("JSON 转换失败:", err)
	}

	// 批量写入数据库
	err = model.StreamInsertData(&jsonData, "https_logs")
	if err != nil {
		log.Printf("❌ 数据插入失败: %v", err)
	} else {
		fmt.Println("✅ 批量数据插入成功")
	}

	// 清空缓存
	buffer = nil
}

// ConvertToJSON 将 TLV 结构体平铺成符合doris格式的 JSON
func ConvertToJSON(tlv []model.TLV) ([]byte, error) {
	// 创建一个对象切片用于JSON编码
	jsonObjects := make([]map[string]interface{}, 0, len(tlv))

	for _, item := range tlv {
		obj := map[string]interface{}{
			"log_id":       item.LogId,
			"type":         item.Type,
			"url":          item.HttpsData.URL,
			"method":       item.HttpsData.Method,
			"status_code":  item.HttpsData.StatusCode,
			"headers_json": item.HttpsData.Headers,
			"content":      item.HttpsData.Content,
		}
		jsonObjects = append(jsonObjects, obj)
	}

	// 将整个对象切片编码为JSON数组
	return json.Marshal(jsonObjects)
}

// processContent 处理 Content：删除反斜杠 & 限制 1024 字节 & 修复 JSON去掉 `\`，限制 1024 字节，正确转义 JSON
func processContent(content string) string {
	// 删除所有 `\` 避免转义问题
	content = strings.ReplaceAll(content, "\\", "")

	// 限制最大 1024 字节
	if len(content) > 1024 {
		runes := []rune(content)
		limitedRunes := []rune{}
		totalBytes := 0
		for _, r := range runes {
			charSize := utf8.RuneLen(r)
			if totalBytes+charSize > 1024 {
				break
			}
			totalBytes += charSize
			limitedRunes = append(limitedRunes, r)
		}
		content = string(limitedRunes)
	}

	// **正确转义 JSON**
	escapedContent, _ := json.Marshal(content)
	return string(escapedContent[1 : len(escapedContent)-1]) // 去掉最外层的 `"`
}
