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

func (e *Expr) MatchesSpan(s *streampb.Span, trace []*streampb.Span) bool {
	for _, m := range e.matchers {
		if !matchesField(m, m.field().fieldID, s, trace) {
			return false
		}
	}

	return true
}

func (e *Expr) MatchesTrace(t []*streampb.Span) bool {
	for _, s := range t {
		if e.MatchesSpan(s, t) {
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

func matchesField(m ValueMatcher, fieldID []int, s *streampb.Span, trace []*streampb.Span) bool {
	if len(fieldID) == 0 {
		return false
	}

	f := m.field()
	id := f.fieldID[0]

	switch id {
	case FIELD_DURATION:
		return m.compareInt(int(s.Duration))
	case FIELD_NAME:
		return m.compareString(s.Name)
	case FIELD_ATTS:
		if a, ok := s.Attributes[m.field().fieldName]; ok {
			switch a.Type {
			case streampb.KeyValuePair_DOUBLE:
				return m.compareFloat(a.DoubleValue)
			case streampb.KeyValuePair_INT:
				return m.compareInt(int(a.IntValue))
			case streampb.KeyValuePair_STRING:
				return m.compareString(a.StringValue)
			case streampb.KeyValuePair_BOOL:
				if a.BoolValue {
					return m.compareInt(1)
				} else {
					return m.compareInt(0)
				}
			}
		}
	case FIELD_EVENTS:
		if e, ok := s.Events[m.field().fieldName]; ok {
			switch e.Type {
			case streampb.KeyValuePair_DOUBLE:
				return m.compareFloat(e.DoubleValue)
			case streampb.KeyValuePair_INT:
				return m.compareInt(int(e.IntValue))
			case streampb.KeyValuePair_STRING:
				return m.compareString(e.StringValue)
			case streampb.KeyValuePair_BOOL:
				if e.BoolValue {
					return m.compareInt(1)
				} else {
					return m.compareInt(0)
				}
			}
		}
	case FIELD_STATUS:
		// unsafe check for code/msg
		subfield := fieldID[1]
		if subfield == FIELD_CODE {
			return m.compareInt(int(s.Status.Code))
		}
		if subfield == FIELD_MSG {
			return m.compareString(s.Status.Message)
		}
	case FIELD_PROCESS:
		// unsafe check
		subfield := fieldID[1]
		if subfield == FIELD_NAME {
			return m.compareString(s.Process.Name)
		}
	case FIELD_PARENT:
		if int(s.ParentIndex) < len(trace) {
			return matchesField(m, fieldID[1:], trace[s.ParentIndex], trace)
		}
	case FIELD_DESCENDANT:
		parentIdx := s.ParentIndex
		for parentIdx >= 0 && int(parentIdx) < len(trace) {
			s = trace[parentIdx]

			if matchesField(m, fieldID[1:], s, trace) {
				return true
			}

			parentIdx = s.ParentIndex
		}
	case FIELD_SPAN:
		return matchesField(m, fieldID[1:], s, trace)
	case FIELD_ROOT_SPAN:
		return matchesField(m, fieldID[1:], trace[0], trace) // assumes trace[0] is the rootspan
	}

	return false
}
