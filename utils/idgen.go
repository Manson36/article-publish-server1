package utils

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// id format:
// timestampBits(41) | workerBits(10) | sequenceBits(13)
// 64位, 1位符号位 + 1位粒度类型(时间戳单位是s还是10ms) + 机器id占10位 + 时间戳占32位 + 自增序列占20位 = 64

const (
	workerIDBits       = uint64(10)                              // 机器id占的位数
	maxWorkerID        = int64(-1) ^ (int64(-1) << workerIDBits) // 最大机器id 1023
	timestampBits      = uint64(32)                              // 时间戳占32位
	sequenceBits       = uint64(20)                              // 自增序列占20位
	timestampLeftShift = sequenceBits                            // 时间戳左移20位
	workerIDLeftShift  = timestampBits + sequenceBits            // 机器id左移52
	sequenceMax        = int64(-1) ^ (int64(-1) << sequenceBits) // 自增序列号最大值

	// 粒度暂定为s级
	// 2019-12-04T00:00:00.000Z
	epochTime = int64(1575388800) // 起始日期
)

//Flags
var (
	workerID int64      // workr id 0 <= workerID <= maxWorkerID
	lastTs   int64 = -1 //the last timestamp in milliseconds
)

var (
	mu  sync.Mutex
	seq int64
)

func init() {
	if workerID < 0 || workerID > maxWorkerID {
		log.Fatalf("worker id must be between 0 and %d", maxWorkerID)
	}
}

func GetInt64ID() (int64, error) {
	mu.Lock()
	defer mu.Unlock()
	ts := time.Now().Unix() //获取当前时间

	switch {
	case ts < lastTs:
		return 0, fmt.Errorf("time is moving backwords, waiting util %d", lastTs)
	case ts == lastTs:
		seq = (seq + 1) & sequenceMax
		if seq == 0 {
			for ts <= lastTs {
				time.Sleep(time.Second)
				ts = time.Now().Unix()
			}
		}
	default:
		seq = 0
	}

	lastTs = ts
	return (workerID << workerIDLeftShift) | ((ts - epochTime) << timestampLeftShift) | seq, nil
}
