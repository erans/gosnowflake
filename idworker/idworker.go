package idworker

import (
	"fmt"
	"github.com/erans/gosnowflake/snowflake"
	"sync"
	"time"
)

const (
	nano = 1000 * 1000
)

const (
	WorkerIdBits       = 5
	DatacenterIdBits   = 5
	MaxWorkerId        = -1 ^ (-1 << WorkerIdBits)
	MaxDatacenterId    = -1 ^ (-1 << DatacenterIdBits)
	SequenceBits       = 12
	WorkerIdShift      = SequenceBits
	DatacenterIdShift  = SequenceBits + WorkerIdBits
	TimestampLeftShift = SequenceBits + WorkerIdBits + DatacenterIdBits
	SequenceMask       = -1 ^ (-1 << SequenceBits)
)

var (
	Epoch uint64 = 1288834974657 /* tweet poch */
)

type IdWorker struct {
	snowflake.Snowflake

	lastTimestamp uint64
	workerId      uint64
	datacenterId  uint64
	sequence      uint64
	lock          sync.Mutex
}

func timeGen() uint64 {
	return uint64(time.Now().UnixNano() / nano)
}

func timestamp() uint64 {
	return uint64(time.Now().UnixNano() / nano)
}

func tillNextMillis(ts uint64) uint64 {
	i := timestamp()
	for i <= ts {
		i = timestamp()
	}
	return i
}

func (worker *IdWorker) GetWorkerId() (int64, error) {
	return int64(worker.workerId), nil
}

func (worker *IdWorker) GetDatacenterId() (int64, error) {
	return int64(worker.datacenterId), nil
}

func (worker *IdWorker) GetTimestamp() (int64, error) {
	return int64(time.Now().UnixNano() / nano), nil
}

func (worker *IdWorker) GetId(useragent string) (r int64, err error) {
	id, err := worker.Next()
	return int64(id), err
}

func (worker *IdWorker) Next() (uint64, error) {
	worker.lock.Lock()
	defer worker.lock.Unlock()

	ts := timeGen()
	if ts < worker.lastTimestamp {
		err := fmt.Errorf("Clock is moving backwards. Rejecting requests until %d.", worker.lastTimestamp)
		return 1, err
	}

	if worker.lastTimestamp == ts {
		worker.sequence = (worker.sequence + 1) & SequenceMask
		if worker.sequence == 0 {
			ts = tillNextMillis(ts)
		}
	} else {
		worker.sequence = 0
	}

	worker.lastTimestamp = ts

	id := ((worker.lastTimestamp - Epoch) << TimestampLeftShift) |
		(worker.datacenterId << DatacenterIdShift) |
		(worker.workerId << WorkerIdShift) |
		worker.sequence

	return id, nil
}

func NewIdWorker(workerId uint64, datacenterId uint64) (*IdWorker, error) {
	if workerId > MaxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("workerId can't be greater than %d or less than 0", workerId)
	}

	if datacenterId > MaxDatacenterId || datacenterId < 0 {
		return nil, fmt.Errorf("datacenterId can't be greater than %d or less than 0", datacenterId)
	}
	return &IdWorker{workerId: workerId, datacenterId: datacenterId, lastTimestamp: 1, sequence: 0}, nil
}
