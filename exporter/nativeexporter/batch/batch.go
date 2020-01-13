package batch

import (
	"encoding/binary"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/willf/bloom"

	"github.com/golang/protobuf/proto"
)

type Batch interface {
	BloomFilter() ([]byte, error)
	Index() ([]byte, error)
	Traces() ([]byte, error)
}

type batch struct {
	bloom    *bloom.BloomFilter
	traces   []byte
	mappings []traceMapping
}

type traceMapping struct {
	id     []byte
	start  uint64
	length uint32
}

func (b *batch) BloomFilter() ([]byte, error) {
	return b.bloom.GobEncode()
}

func (b *batch) Index() ([]byte, error) {
	// sort and serialize.  let's do the dumb thing to start and just write the traceid, start and length.
	//  this will be large but will allow for fast lookups on a sorted index
	// 28 bytes per entry
	// 128 (trace id) + 64 (start) + 32 (length) = 1612800000 bits ~= 200MB for 1000 traces/second cut every 2 hours

	// todo:  consider compression (columnar diff).  index must clearly indicate if a trace id is in a set
	bytes := make([]byte, 0, len(b.mappings)*28)
	buff := bytes

	for _, m := range b.mappings {
		if len(m.id) != 16 {
			return nil, fmt.Errorf("Trace Ids must be 128 bit")
		}

		// todo: does golang have memcpy?
		for i, idByte := range m.id {
			buff[i] = idByte
		}
		binary.LittleEndian.PutUint64(buff, m.start)
		binary.LittleEndian.PutUint32(buff, m.length)

		buff = buff[28:]
	}

	return bytes, nil
}

func (b *batch) Traces() ([]byte, error) {
	return b.traces, nil
}

func newBatch(traceCount uint) *batch {
	// todo: choose a bloom filter size to hit a configurable false positive rate
	//   the following assumes 1000 traces/second cut every 2 hours.  n = 7200000 with a target fp rate of 1%
	//   https://hur.st/bloomfilter/ m = 69069274, k = 7
	return &batch{
		bloom: bloom.New(69069274, 7),
	}
}

func (b *batch) addTraceData(td []consumerdata.TraceData) error {
	if len(td) == 0 || len(td[0].Spans) == 0 {
		return nil
	}

	traceID := td[0].Spans[0].TraceId
	start := len(b.traces) - 1

	// todo: spans from a process maybe distributed amonst multiple tracedata's.  consolidate here?

	// serialize and append to traces
	varIntBuf := make([]byte, binary.MaxVarintLen64)

	b.traces = append(b.traces, varIntBuf[:binary.PutUvarint(varIntBuf, uint64(len(td)))]...)
	for _, t := range td {
		bytes, err := proto.Marshal(t.Node)
		if err != nil {
			b.traces = []byte{}
			return err
		}
		b.traces = append(b.traces, bytes...)

		b.traces = append(b.traces, varIntBuf[:binary.PutUvarint(varIntBuf, uint64(len(t.Spans)))]...)
		for _, s := range t.Spans {
			bytes, err = proto.Marshal(s)
			if err != nil {
				b.traces = []byte{}
				return err
			}

			b.traces = append(b.traces, bytes...)
		}
	}

	// append to traceIDs
	b.mappings = append(b.mappings, traceMapping{
		id:     traceID,
		start:  uint64(start),
		length: uint32(len(b.traces) - 1 - start),
	})

	// add to bloom filter
	b.bloom.Add(traceID)

	return nil
}
