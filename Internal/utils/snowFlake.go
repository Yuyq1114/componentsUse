package utils

import (
	"fmt"
	"sync"
	"time"
)

// 雪花算法参数
const (
	epoch int64 = 1704067200000 // 2024-01-01 00:00:00 (UTC)
	// 自定义起始时间（毫秒级，避免时间戳溢出）
	machineBits  uint = 10 // 机器 ID 位数（最大支持 1024 台机器）
	sequenceBits uint = 12 // 序列号位数（每毫秒最多 4096 个 ID）

	maxMachineID   int64 = -1 ^ (-1 << machineBits)   // 最大机器 ID
	maxSequence    int64 = -1 ^ (-1 << sequenceBits)  // 最大序列号
	machineShift   uint  = sequenceBits               // 机器 ID 左移位数
	timestampShift uint  = sequenceBits + machineBits // 时间戳左移位数
)

// 雪花 ID 生成器
type Snowflake struct {
	mu        sync.Mutex
	lastTime  int64
	machineID int64
	sequence  int64
}

// 初始化雪花算法
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("机器 ID 超出范围: 0-%d", maxMachineID)
	}
	return &Snowflake{machineID: machineID}, nil
}

// 生成唯一 ID
func (s *Snowflake) GenerateID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	// 如果时间戳相同，则序列号递增
	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 当前毫秒内序列号用尽，等待下一毫秒
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	// 组合 ID
	id := ((now - epoch) << timestampShift) | (s.machineID << machineShift) | s.sequence
	return id
}
