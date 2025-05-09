package model

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

// User 定义 user 表的结构
// ************此处的json对应解析的json的key名称，千万注意大小写*****************
type ExampleTbl struct {
	Timestamp time.Time `json:"timestamp"`
	Type      int       `json:"type"`
	ErrorCode int       `json:"error_code"`
	ErrorMsg  string    `json:"error_msg"`
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

func StreamInsertData(exampleData *[]byte, table string) error {

	// 将结构体数据序列化为 JSON
	//dataSource: "root:mypassword@tcp(117.50.85.130:9030)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	//FEIP: "117.50.85.130"
	//FEPORT: "8040"
	//FEDB: "mydb"
	//jsonData1 := bytes.NewReader(jsonData)
	// 构建 HTTP 请求，注意端口
	url := fmt.Sprintf("http://117.50.85.130:8040/api/mydb/%s/_stream_load", table)
	//url := "http://117.50.85.130:8040/api/mydb/example_tbl/_stream_load"
	dataOut := bytes.NewReader(*exampleData)
	req, err := http.NewRequest("PUT", url, dataOut)
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}
	root := "root"
	pass := "mypassword"
	auth := base64.StdEncoding.EncodeToString([]byte(root + ":" + pass))
	// 设置请求头
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Expect", "100-continue")
	var u1 = uuid.Must(uuid.NewV4())
	req.Header.Add("label", u1.String())
	req.Header.Add("format", "json")
	req.Header.Add("strip_outer_array", "True") // DorisDB可能期望的是一个JSON数组

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = sonic.Unmarshal(all, &result)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)

	//处理响应
	if result["Status"] == "Success" {
		return nil
	} else {
		return errors.New("return stauts is  false")
	}
}

//位于entrance的runTask中的测试doris的数据
//--------------------------------------------------------------
//example := model.ExampleTbl{
//	Timestamp: time.Now(),
//	Type:      1,
//	ErrorCode: 404,
//	ErrorMsg:  "Not Found",
//}
//example.InsertData(dorisDb)
//
//fmt.Println(" 插入完成")
//-------------------------------------------------------------
//用于测试doris的数据
//******************* 注意必须是这样的，应为可以有两个[]中阔号
/*exampleData := []model.ExampleTbl{
	{
		Timestamp: time.Now(),
		Type:      1,
		ErrorCode: 404,
		ErrorMsg:  "Not Found",
	},
}

jsonData, err := sonic.Marshal(exampleData)
fmt.Println(string(jsonData))
err = model.StreamInsertData(&jsonData, config.Doris, "example_tbl")
if err != nil {
	fmt.Println("数据插入失败")
} else {
	fmt.Println("数据插入成功")
}
fmt.Println("任务开始")*/
//------------------------------------------------------------------
