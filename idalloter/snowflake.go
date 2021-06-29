// Package idalloter 雪花算法
// +----------------------------------------------------------+
// | 42Bit lastTimestamp | 10Bit sequenceID | 12Bit machineID |
// +----------------------------------------------------------+
package idalloter

import (
	"errors"
	"sync"
	"time"
)

// maxSequenceID 为同一毫秒内产生的 id 分配 12 位大小区间，最大值为 2^12-1
const maxSequenceID = 1<<12 - 1

// maxMachineID 为机器编号分配 10 位大小区间，最大值为 2^10-1
const maxMachineID = 1<<10 - 1

// SnowFlake 雪花算法结构体，64字节长度的 id 共由三部分组成
// 			 lastTimestamp 42bit    记录最近一次的时间戳，单位粒度毫秒
// 		     sequenceID    12bit	记录同一毫秒内区分的序列号，解决同一毫秒请求的碰撞问题
// 			 machineID     10bit	记录机器id，支持分布式
type SnowFlake struct {
	lastTimestamp uint64
	sequenceID    uint32
	machineID     uint32
	lock          sync.Mutex
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
		return nil, errors.New("invalid machine id")
	}

	return &SnowFlake{machineID: machineId}, nil
}

// nextMillisecond 获取下一个毫秒的时间戳
func (sf *SnowFlake) nextMillisecond(ntime uint64) uint64 {
	for ntime == sf.lastTimestamp {
		time.Sleep(time.Microsecond * 100)
	}

	return getTimestamp()
}

// getTimestamp 获取当前时间戳，单位毫秒
func getTimestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

// getId 根据记录的结构，得到唯一id
func (sf *SnowFlake) getId() uint64 {
	// 时间戳为高位，先将 lastTimestamp 向左移动 12 + 10 个bit，预留出 sequenceID 和 machineID 的位置
	// 再将 sequenceID 向左移动 10 个bit，为 machineID 预留这 10 个bit
	// 最后将 machineID 填补到这 10 个bit
	return sf.lastTimestamp<<(12+10) | uint64(sf.sequenceID)<<10 | uint64(sf.machineID)
}

// GenerateID 获取唯一id
func (sf *SnowFlake) GenerateID() (uint64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	// 判断当前系统得到的时间戳是否小于记录的最小时间戳，判断系统是否异常
	nowTime := getTimestamp()
	if nowTime < sf.lastTimestamp {
		return 0, errors.New("invalid system time")
	}

	// 若系统时间戳和记录的最后时间戳相等，说明在该毫秒内产生了碰撞，利用 sequenceID 解决碰撞
	if nowTime == sf.lastTimestamp {
		sf.sequenceID = (sf.sequenceID + 1) & maxSequenceID
		if sf.sequenceID == 0 {
			nowTime = sf.nextMillisecond(nowTime)
		}
	} else {
		sf.sequenceID = 0
	}

	// 记录本次分配 id 的时间戳
	sf.lastTimestamp = nowTime

	return sf.getId(), nil
}
