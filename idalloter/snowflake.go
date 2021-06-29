// Package idalloter 雪花算法
// +----------------------------------------------------------------------+
// | 1 bit ｜ 41 Bit lastStamp | 10 Bit machineID | 12 Bit sequenceID |
// +----------------------------------------------------------------------+
// 1  bit，这里最高位为标识位预留出来，默认为0
// 41 bit，毫秒级时间戳，存储时间戳差值（当前时间戳-开始时间戳）；开始时间戳可任意指定，一般为系统上线时间，使用该库时建议调整，最多可使用 69 年
// 10 bit，机器标识，支持最多部署 1024 个节点
// 12 bit，计数序列号，序列号为一系列自增 id，支持每个节点每毫秒最多产生 4096 个序号
package idalloter

import (
	"errors"
	"sync"
	"time"
)

const (
	lastStampBit = 41 // 上一次时间戳占用 bit
	sequenceBit  = 12 // 序列号占用 bit
	machineBit   = 10 // 机器识别号占用 bit

	maxSequenceID        = 1<<sequenceBit - 1 // 序列号，12 bit，最大值为 2^12-1
	maxMachineID         = 1<<machineBit - 1  // 机器标识号，10 bit，最大值为 2^10-1
	startStamp    uint64 = 1622476800000      // 开始时间戳，2021-06-01 00:00:00 000
)

// SnowFlake 雪花算法结构体，64字节长度的 id 共由三部分组成
//           lastStamp     42bit    上一次时间戳，单位毫秒
//           sequenceID    12bit    序列号，使同一毫秒内不出现重复
//           machineID     10bit    机器标识id，使分布式系统内不同机器间不出现重复
type SnowFlake struct {
	lastStamp  uint64
	sequenceID uint32
	machineID  uint32
	lock       sync.Mutex
}

// SnowFlakeInterface 接口整理
type SnowFlakeInterface interface {
	// NewSnowFlake 获取 id 分配器对象
	NewSnowFlake(machineId uint32) (*SnowFlake, error)
	// GenerateID 获取唯一id
	GenerateID() (uint64, error)
}

// NewSnowFlake 获取 id 分配器对象
func NewSnowFlake(machineId uint32) (*SnowFlake, error) {
	if machineId < 0 || machineId > maxMachineID {
		return nil, errors.New("invalid machineId")
	}

	return &SnowFlake{machineID: machineId}, nil
}

// getNextMillisecond 获取下一个毫秒的时间戳
func (sf *SnowFlake) getNextMillisecond(ntime uint64) uint64 {
	for ntime == sf.lastStamp {
		time.Sleep(time.Microsecond * 100)
	}

	return getStamp()
}

// getStamp 获取当前时间戳差值
func getStamp() uint64 {
	return uint64(time.Now().UnixNano()/int64(time.Millisecond)) - startStamp
}

// getId 根据记录的结构，得到唯一id
func (sf *SnowFlake) getId() uint64 {
	return sf.lastStamp<<(machineBit+sequenceBit) | uint64(sf.machineID)<<machineBit | uint64(sf.sequenceID)
}

// GenerateID 获取唯一id
func (sf *SnowFlake) GenerateID() (uint64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	nowTime := getStamp()
	if nowTime < sf.lastStamp {
		return 0, errors.New("invalid system time")
	} else if nowTime > sf.lastStamp {
		sf.sequenceID = 0
	} else {
		sf.sequenceID = (sf.sequenceID + 1) & maxSequenceID
		if sf.sequenceID == 0 {
			nowTime = sf.getNextMillisecond(nowTime)
		}
	}
	sf.lastStamp = nowTime

	return sf.getId(), nil
}
