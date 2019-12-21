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
	// todo:  sort matchers by execution cost so cheapest are executed first
	return &Expr{
		stream:   stream,
		matchers: m,
	}
}

func (e *Expr) MatchesSpan(s *streampb.Span) bool {
	for _, m := range e.matchers {
		if !matchesField(m, s) {
			return false
		}
	}

	return true
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
		fields := m.field().fieldID

		for _, f := range fields {
			if f == FIELD_DESCENDANT || f == FIELD_PARENT {
				return true
			}
		}
	}

	return false
}

func matchesField(m ValueMatcher, s *streampb.Span) bool {
	f := m.field()

	if len(f.fieldID) == 0 {
		return false
	}

	id := f.fieldID[0]

	switch id {
	case FIELD_DURATION:
		return m.compareInt(int(s.Duration))
	case FIELD_NAME:
		return m.compareString(s.OperationName)
	case FIELD_ATTS:
		// todo
	case FIELD_EVENTS:
		// todo
	case FIELD_STATUS:
		// todo
	case FIELD_CODE:
		// todo
	case FIELD_MSG:
		// todo
	case FIELD_PROCESS:
		// todo
	case FIELD_PARENT:
		// todo
	case FIELD_DESCENDANT:
		// todo
	case FIELD_SPAN:
		// todo
	case FIELD_ROOT_SPAN:
		// todo
	}

	return false
}
