package streamprocessor

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

type batcher struct {
	maxBatches   int
	batchTimeout time.Duration
	batches      map[uint64]*batch
	lock         sync.RWMutex
}

type batch struct {
	lastAppend time.Time
	trace      []*streampb.Span
}

func newBatcher() *batcher {
	return &batcher{
		maxBatches:   50,
		batches:      make(map[uint64]*batch),
		batchTimeout: 5 * time.Second,
	}
}

// all spans must be from the same traceid
func (b *batcher) addBatch(spans []*streampb.Span) error {
	if len(spans) == 0 {
		return nil
	}

	if len(spans) >= b.maxBatches {
		return fmt.Errorf("batcher is full")
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	traceID := traceID(spans[0].TraceID)
	existing, ok := b.batches[traceID]

	if !ok {
		existing = &batch{}
		b.batches[traceID] = existing
	}

	existing.lastAppend = time.Now()
	existing.trace = append(existing.trace, spans...)

	return nil
}

// after calling completeBatches it is the responsibility of the caller do to something with them
func (b *batcher) completeBatches() [][]*streampb.Span {
	var completed [][]*streampb.Span

	b.lock.RLock()
	defer b.lock.RUnlock()

	for id, batch := range b.batches {
		if time.Now().After(batch.lastAppend.Add(b.batchTimeout)) {
			completed = append(completed, batch.trace)

			delete(b.batches, id)
		}
	}

	return completed
}

// let's hope there's no collisions in the last 8 bytes!
//  todo:  use full byte slice.  this only works b/c we generate 64 bit traceids in jaeger
func traceID(b []byte) uint64 {
	lastBytes := b[len(b)-8:]
	return binary.BigEndian.Uint64(lastBytes)
}
