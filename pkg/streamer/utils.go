package streamer

import (
	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

func spanToSpan(in *tracepb.Span) *blergpb.Span {

	return &blergpb.Span{
		OperationName: in.Name.Value,
		StartTime:     in.StartTime.Seconds,
		Duration:      in.EndTime.Seconds - in.StartTime.Seconds,
	}
}
