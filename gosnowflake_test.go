package main

import (
	"fmt"
	"github.com/erans/gosnowflake/idworker"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

func TestNewIdWorker(t *testing.T) {
	worker, err := idworker.NewIdWorker(1, 1)
	assert.Nil(t, err)

	id1, err := worker.Next()
	assert.Nil(t, err)

	id2, err := worker.Next()
	assert.Nil(t, err)

	if id1 > id2 {
		t.Errorf("id2 %v is smaller then previous one %v", id2, id1)
	}
}

func TestDuplicateKeys(t *testing.T) {
	worker, err := idworker.NewIdWorker(1, 1)
	assert.Nil(t, err)

	const MAXIDS = 100000

	var ids [MAXIDS]uint64

	start := time.Now()

	for i := 0; i < MAXIDS; i++ {
		id, err := worker.Next()
		assert.Nil(t, err)

		ids[i] = id
	}

	for i := 1; i < MAXIDS; i++ {
		if ids[i-1] > ids[i] {
			t.Errorf("next id %v is smaller then previous one %v", ids[i], ids[i-1])
		}
	}

	end := time.Now()

	delta := float64(end.UnixNano()-start.UnixNano()) / float64(1000*1000)
	t.Logf("Execution time: %fms\n", delta)
	t.Logf("Ids/sec: %f\n", float64(MAXIDS)/delta)
}

func BenchmarkIdWorker(b *testing.B) {
	worker, err := idworker.NewIdWorker(1, 1)
	if err == nil {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id, err := worker.Next()
			if err == nil {
				fmt.Sprintf("%d", id)
			} else {
				b.Errorf("Failed to create id")
			}
		}
	} else {
		b.Errorf("Failed to create IdWorker")
	}
}
