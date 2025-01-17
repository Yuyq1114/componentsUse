package main

import (
	"test_component/cmd"
)

func main() {

	cmd.Execute()

}

// -------------------------------------------------------------------------------------------
//package main
//
//import (
//	"bytes"
//	"encoding/base64"
//	"encoding/json"
//	"fmt"
//	"github.com/bytedance/sonic"
//	"github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"
//	"io"
//	"net/http"
//	"time"
//)
//
//type StreamLoad struct {
//	url       string
//	dbName    string
//	tableName string
//	data      string
//	userName  string
//	password  string
//}
//
//// 实现Doris用户认证信息
//func auth(load StreamLoad) string {
//	s := load.userName + ":" + load.password
//	b := []byte(s)
//
//	sEnc := base64.StdEncoding.EncodeToString(b)
//	fmt.Printf("enc=[%s]\n", sEnc)
//
//	sDec, err := base64.StdEncoding.DecodeString(sEnc)
//	if err != nil {
//		fmt.Printf("base64 decode failure, error=[%v]\n", err)
//	} else {
//		fmt.Printf("dec=[%s]\n", sDec)
//	}
//	return sEnc
//}
//
//// 内存流数据，通过Stream Load导入Doris表中
//func batch_load_data(load StreamLoad, data *[]byte) {
//	client := &http.Client{
//		Timeout: 30 * time.Second,
//	}
//	//生成要访问的url
//	url := "http://117.50.85.130:8030/api/mydb/user_info/_stream_load"
//	//fmt.Formatter(.Format(url,load.dbName,l))
//
//	record := bytes.NewReader(*data)
//	//提交请求
//	reqest, err := http.NewRequest(http.MethodPut, url, record)
//
//	//增加header选项
//	reqest.Header.Add("Authorization", "basic "+auth(load))
//	reqest.Header.Add("EXPECT", "100-continue")
//	var u1 = uuid.Must(uuid.NewV4())
//	reqest.Header.Add("label", u1.String())
//	reqest.Header.Add("column_separator", ",")
//	//reqest.Header.Add("strip_outer_array", "true")
//
//	if err != nil {
//		panic(err)
//	}
//	//处理返回结果
//	fmt.Println(reqest)
//
//	response, err := client.Do(reqest)
//	if err != nil {
//		fmt.Println(err)
//	}
//	if response.StatusCode == 200 {
//		body, _ := io.ReadAll(response.Body)
//		responseBody := ResponseBody{}
//		jsonStr := string(body)
//		err := json.Unmarshal([]byte(jsonStr), &responseBody)
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//		if responseBody.Status == "Success" {
//			//如果有被过滤的数据，打印错误的URL
//			if responseBody.NumberFilteredRows > 0 {
//				fmt.Printf("Error Data : %s ", responseBody.ErrorURL)
//			} else {
//				fmt.Printf("Success import data : %d", responseBody.NumberLoadedRows)
//			}
//		} else {
//			fmt.Printf("Error Message : %s \n", responseBody.Message)
//			fmt.Printf("Error Data : %s ", responseBody.ErrorURL)
//		}
//		//fmt.Println(jsonStr)
//	}
//	defer response.Body.Close()
//}
//
//// Stream load返回消息结构体
//type ResponseBody struct {
//	TxnID                  int    `json:"TxnId"`
//	Label                  string `json:"Label"`
//	Status                 string `json:"Status"`
//	Message                string `json:"Message"`
//	NumberTotalRows        int    `json:"NumberTotalRows"`
//	NumberLoadedRows       int    `json:"NumberLoadedRows"`
//	NumberFilteredRows     int    `json:"NumberFilteredRows"`
//	NumberUnselectedRows   int    `json:"NumberUnselectedRows"`
//	LoadBytes              int    `json:"LoadBytes"`
//	LoadTimeMs             int    `json:"LoadTimeMs"`
//	BeginTxnTimeMs         int    `json:"BeginTxnTimeMs"`
//	StreamLoadPutTimeMs    int    `json:"StreamLoadPutTimeMs"`
//	ReadDataTimeMs         int    `json:"ReadDataTimeMs"`
//	WriteDataTimeMs        int    `json:"WriteDataTimeMs"`
//	CommitAndPublishTimeMs int    `json:"CommitAndPublishTimeMs"`
//	ErrorURL               string `json:"ErrorURL"`
//}
//
//// 获取BE列表返回结构体
//type Backend struct {
//	Msg  string `json:"msg"`
//	Code int    `json:"code"`
//	Data struct {
//		Backends []struct {
//			IP       string `json:"ip"`
//			HTTPPort int    `json:"http_port"`
//			IsAlive  bool   `json:"is_alive"`
//		} `json:"backends"`
//	} `json:"data"`
//	Count int `json:"count"`
//}
//type UserInfo struct {
//	UserID       int64      `json:"user_id"`       // 用户id
//	Username     string     `json:"username"`      // 用户名
//	City         *string    `json:"city"`          // 用户所在城市 (可为空)
//	Age          *int16     `json:"age"`           // 用户年龄 (可为空)
//	Sex          *int8      `json:"sex"`           // 用户性别 (可为空)
//	Phone        *int64     `json:"phone"`         // 电话 (可为空)
//	Address      *string    `json:"address"`       // 地址 (可为空)
//	RegisterTime *time.Time `json:"register_time"` // 用户注册时间 (可为空)
//}
//
//func main() {
//	var load StreamLoad
//	load.userName = "root"
//	load.password = "mypassword"
//	//auth_info := auth(load)
//	//fmt.Println(auth_info)
//	//backends := get_doris_be_list()
//	//for e := backends.Front(); e != nil; e = e.Next() {
//	// fmt.Println(e.Value)
//	//}
//	user1 := []UserInfo{
//		{
//			UserID:       10001,
//			Username:     "root",
//			City:         nil,
//			Age:          nil,
//			Sex:          nil,
//			Phone:        nil,
//			Address:      nil,
//			RegisterTime: nil,
//		},
//	}
//
//	marshal, _ := sonic.Marshal(user1)
//
//	batch_load_data(load, &marshal)
//	//batch_load_file(/load, "/Users/zhangfeng/Downloads/test.csv")
//}
