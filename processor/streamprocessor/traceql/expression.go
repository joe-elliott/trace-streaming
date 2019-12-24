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
		if !matchesField(m, m.field().id, s) {
			return false
		}
	}

	return true
}

func (e *Expr) MatchesTrace(t []*streampb.Span) bool {

	// check other filters
	for _, s := range t {
		matches := true

		for _, m := range e.matchers {
			if !matchesTraceField(m, m.field().id, s, t) {
				matches = false
				break
			}
		}

		if matches {
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
		fields := m.field().id

		for _, f := range fields {
			if f == FIELD_DESCENDANT || f == FIELD_PARENT {
				return true
			}
		}
	}

	return false
}

func matchesTraceField(m ValueMatcher, id fieldID, s *streampb.Span, t []*streampb.Span) bool {
	if len(id) == 0 {
		return false
	}

	rootID := id[0]

	switch rootID {
	case FIELD_PARENT:
		if int(s.ParentIndex) < len(t) {
			return matchesField(m, id[1:], t[s.ParentIndex])
		}

		return false
	case FIELD_DESCENDANT:
		for _, s := range t {
			parentIdx := s.ParentIndex
			for parentIdx >= 0 && int(parentIdx) < len(t) {
				s = t[parentIdx]

				if matchesField(m, id[1:], s) {
					return true
				}

				parentIdx = s.ParentIndex
			}
		}
	case FIELD_SPAN:
		return matchesField(m, id[1:], s)
	}

	return matchesField(m, id, s)
}

func matchesField(m ValueMatcher, id fieldID, s *streampb.Span) bool {
	if len(id) == 0 {
		return false
	}

	rootID := id[0]

	switch rootID {
	case FIELD_DURATION:
		return m.compareInt(int(s.Duration))
	case FIELD_NAME:
		return m.compareString(s.Name)
	case FIELD_ATTS:
		if a, ok := s.Attributes[m.field().name]; ok {
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
		if e, ok := s.Events[m.field().name]; ok {
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
		subfield := id[1]
		if subfield == FIELD_CODE {
			return m.compareInt(int(s.Status.Code))
		}
		if subfield == FIELD_MSG {
			return m.compareString(s.Status.Message)
		}
	case FIELD_PROCESS:
		// unsafe check
		subfield := id[1]
		if subfield == FIELD_NAME {
			return m.compareString(s.Process.Name)
		}
	case FIELD_IS_ROOT:
		isRoot := 0
		if len(s.ParentSpanID) == 0 {
			isRoot = 1
		}
		return m.compareInt(isRoot)
	}

	return false
}
