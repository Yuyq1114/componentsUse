package model

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// User 定义 user 表的结构
type ExampleTbl struct {
	Timestamp time.Time `gorm:"column:timestamp;not null"` // 日期
	Type      int       `gorm:"column:type;not null"`      // 类型 (TINYINT)
	ErrorCode int       `gorm:"column:error_code"`         // 错误代码 (INT)
	ErrorMsg  string    `gorm:"column:error_msg;size:300"` // 错误消息 (VARCHAR)
}

// TableName 指定表名
func (ExampleTbl) TableName() string {
	return "example_tbl"
}

/*CREATE TABLE IF NOT EXISTS mydb.example_tbl
(
`timestamp` DATE NOT NULL COMMENT "['0000-01-01', '9999-12-31']",
`type` TINYINT NOT NULL COMMENT "[-128, 127]",
`error_code` INT COMMENT "[-2147483648, 2147483647]",
`error_msg` VARCHAR(300) COMMENT "[1-65533]"
)
DISTRIBUTED BY HASH(`timestamp`) BUCKETS 1
PROPERTIES (
"replication_num" = "1"
);*/

/*CREATE TABLE IF NOT EXISTS user_info
(
user_id LARGEINT NOT NULL COMMENT "用户id",
username varchar(50) NOT NULL COMMENT "用户名",
city VARCHAR(20) COMMENT "用户所在城市",
age SMALLINT COMMENT "用户年龄",
sex TINYINT COMMENT "用户性别",
phone LARGEINT COMMENT "电话",
address VARCHAR(500) COMMENT "地址",
register_time datetime COMMENT "用户注册时间"
)
Unique KEY(user_id, username)
DISTRIBUTED BY HASH(user_id) BUCKETS 3
PROPERTIES (
"replication_num" = "1"
);*/

func (u *ExampleTbl) InitTable(db *gorm.DB) {
	err := db.AutoMigrate(u)
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
}

func (u *ExampleTbl) InsertData(db *gorm.DB) error {
	return db.Create(u).Error
}

func StreamInsertData(exampleData *[]byte) error {

	// 将结构体数据序列化为 JSON

	//jsonData1 := bytes.NewReader(jsonData)
	// 构建 HTTP 请求
	url := "http://117.50.85.130:8040/api/mydb/example_tbl/_stream_load"
	dataOut := bytes.NewReader(*exampleData)
	req, err := http.NewRequest("PUT", url, dataOut)
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}
	root := "root"
	pass := "password"
	auth := base64.StdEncoding.EncodeToString([]byte(root + ":" + pass))
	// 设置请求头
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Expect", "100-continue")
	var u1 = uuid.Must(uuid.NewV4())
	req.Header.Add("label", u1.String())
	req.Header.Add("format", "json")
	req.Header.Add("strip_outer_array", "true") // 可选，设置请求超时（毫秒）

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode == http.StatusOK {
		fmt.Println("数据插入成功！")
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	return err
}
