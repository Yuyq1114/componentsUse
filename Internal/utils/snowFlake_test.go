package utils

import (
	"fmt"
	"testing"
	//"utils"
)

func TestGenerateID(t *testing.T) {
	sf, err := NewSnowflake(1) // 机器 ID 1
	if err != nil {
		t.Fatalf("初始化 Snowflake 失败: %v", err)
	}

	// 存储生成的 ID 进行唯一性检查
	//idSet := make(map[int64]bool)
	//var lastID int64

	for i := 0; i < 5; i++ {
		id := sf.GenerateID()
		fmt.Println(id)
		//// 断言 ID 为正数
		//if id <= 0 {
		//	t.Errorf("生成的 ID 不是正数: %d", id)
		//}
		//
		//// 断言 ID 不重复
		//if idSet[int64(id)] {
		//	t.Errorf("生成的 ID 发生重复: %d", id)
		//}
		//idSet[int64(id)] = true
		//
		//// 断言 ID 递增
		//if lastID > 0 && int64(id) <= lastID {
		//	t.Errorf("生成的 ID 没有递增: 上一个 %d, 当前 %d", lastID, id)
		//}
		//
		//lastID = int64(id)
		//t.Logf("生成的 ID: %d", id)
	}
}
