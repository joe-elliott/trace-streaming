package traceql

import (
	"math"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

type aggregationFunc func(s *streampb.Span, reset bool) []float64

func generateAggregationFunc(agg int, f field, args []float64) aggregationFunc {
	switch agg {
	case AGG_MIN:
		min := math.MaxFloat64
		return func(s *streampb.Span, reset bool) []float64 {
			if reset {
				ret := []float64{min}
				min = math.MaxFloat64
				return ret
			}

			min = math.Min(min, f.getFloatValue(s))

			return nil
		}
	case AGG_MAX:
		max := math.SmallestNonzeroFloat64
		return func(s *streampb.Span, reset bool) []float64 {
			if reset {
				ret := []float64{max}
				max = math.SmallestNonzeroFloat64
				return ret
			}

			max = math.Max(max, f.getFloatValue(s))

			return nil
		}
	case AGG_COUNT:
		count := 0.0
		return func(s *streampb.Span, reset bool) []float64 {
			if reset {
				// no reset.  count is a counter
				return []float64{count}
			}

			count++
			return nil
		}
	case AGG_SUM:
		sum := 0.0
		return func(s *streampb.Span, reset bool) []float64 {
			if reset {
				ret := []float64{sum}
				sum = 0.0
				return ret
			}

			sum = sum + f.getFloatValue(s)
			return nil
		}
	case AGG_AVG:
		count := 0.0
		sum := 0.0
		return func(s *streampb.Span, reset bool) []float64 {
			if reset {
				var ret []float64

				if count == 0.0 {
					ret = []float64{0.0}
				} else {
					ret = []float64{sum / count}
				}

				count = 0.0
				sum = 0.0
				return ret
			}

			count++
			sum = sum + f.getFloatValue(s)
			return nil
		}
	case AGG_HIST:
		if len(args) != 3 {
			// todo: propagate error
			return nil
		}

		start := args[0]
		bucketCount := int(args[1])
		diff := args[2]

		buckets := make([]float64, bucketCount+1)

		return func(s *streampb.Span, reset bool) []float64 {
			if reset {
				// no reset.  buckets are counters
				return buckets
			}

			val := f.getFloatValue(s)
			bucket := int((val - start) / diff)

			if bucket < 0 {
				bucket = 0
			}

			if bucket > bucketCount {
				bucket = bucketCount
			}

			buckets[bucket]++

			return nil
		}
	}

	return nil
}
