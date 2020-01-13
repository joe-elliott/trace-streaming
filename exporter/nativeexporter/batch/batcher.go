package batch

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
)

type Batcher interface {
	Store(td consumerdata.TraceData)
	Cut() (Batch, error)
}

type batcher struct {
	lock         sync.RWMutex
	traces       map[uint64]*batchedTrace
	traceTimeout time.Duration
}

type batchedTrace struct {
	lastAppend time.Time
	traceData  []consumerdata.TraceData
}

func NewBatcher() Batcher {
	return &batcher{
		traces:       make(map[uint64]*batchedTrace),
		traceTimeout: 10 * time.Second,
	}
}

func (b *batcher) Store(td consumerdata.TraceData) {
	if len(td.Spans) == 0 {
		return
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	traceID := traceID(td.Spans[0].TraceId)
	existing, ok := b.traces[traceID]

	if !ok {
		existing = &batchedTrace{}
		b.traces[traceID] = existing
	}

	existing.lastAppend = time.Now()
	existing.traceData = append(existing.traceData, td)
}

func (b *batcher) Cut() (Batch, error) {
	// iterate through and remove all "complete" traces
	batch := &batch{}

	b.lock.RLock()
	defer b.lock.RUnlock()

	for id, t := range b.traces {
		if time.Now().After(t.lastAppend.Add(b.traceTimeout)) {
			if err := batch.addTraceData(t.traceData); err != nil {
				return nil, err
			}
			delete(b.traces, id)
		}
	}

	return batch, nil
}

// let's hope there's no collisions in the last 8 bytes!
//  todo:  use full byte slice.  this only works b/c we generate 64 bit traceids in jaeger
func traceID(b []byte) uint64 {
	lastBytes := b[len(b)-8:]
	return binary.BigEndian.Uint64(lastBytes)
}
