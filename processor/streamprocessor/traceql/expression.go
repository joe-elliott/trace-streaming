package traceql

import "github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"

type Query interface {
	MatchesSpan(*streampb.Span) bool
	MatchesTrace([]*streampb.Span) bool
	RequiresTraceBatching() bool
}

//
type Expr struct {
	stream   int
	matchers []ValueMatcher
}

func newExpr(stream int, m []ValueMatcher) *Expr {
	// todo:  sort matchers by execution cost
	return &Expr{
		stream:   stream,
		matchers: m,
	}
}

func (e *Expr) MatchesSpan(s *streampb.Span) bool {
	return false
}

// jpe - change shape of "trace" and make recursive?
func (e *Expr) MatchesTrace(t []*streampb.Span) bool {
	for _, s := range t {
		if e.MatchesSpan(s) {
			return true
		}
	}

	return false
}

// RequiresTraceBatching indicates if this expression expects/requires an entire trace at a time to be evaluated or if it
//  can be done a span at a time
func (e *Expr) RequiresTraceBatching() bool {
	if e.stream == STREAM_TYPE_TRACES {
		return true
	}

	// if any matchers have descendant or parent fields then we require trace batching
	for _, m := range e.matchers {
		var fields []int

		switch v := m.(type) {
		case intMatcher:
			fields = v.field.fieldID
		case floatMatcher:
			fields = v.field.fieldID
		case stringMatcher:
			fields = v.field.fieldID
		}

		for _, f := range fields {
			if f == FIELD_DESCENDANT || f == FIELD_PARENT {
				return true
			}
		}
	}

	return false
}
