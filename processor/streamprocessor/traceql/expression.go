package traceql

import "github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"

type QueryType int

const (
	QueryTypeSpans QueryType = iota
	QueryTypeBatchedSpans
	QueryTypeTraces
	QueryTypeMetrics
)

type Query interface {
	MatchesSpan(*streampb.Span) bool
	MatchesSpanBatched(*streampb.Span, []*streampb.Span) bool
	MatchesTrace([]*streampb.Span) bool
	Aggregate(s *streampb.Span, reset bool) []float64

	QueryType() QueryType
}

//
type Expr struct {
	stream   int
	matchers []matcher

	aggFunc aggregationFunc
}

func newExpr(stream int, m []matcher) *Expr {
	// todo:  sort matchers by execution cost so cheapest are executed first
	return &Expr{
		stream:   stream,
		matchers: m,
	}
}

func newMetricsExpr(agg int, expr *Expr, f field, args []float64) *Expr {
	expr.aggFunc = generateAggregationFunc(agg, f, args)

	return expr
}

func (e *Expr) MatchesSpan(s *streampb.Span) bool {
	for _, m := range e.matchers {
		if !m.compare(s, s) {
			return false
		}
	}

	return true
}

func (e *Expr) MatchesSpanBatched(s *streampb.Span, t []*streampb.Span) bool {
	for _, m := range e.matchers {
		if !matchesTraceField(m, s, t) {
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
			if !matchesTraceField(m, s, t) {
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

func (e *Expr) Aggregate(s *streampb.Span, reset bool) []float64 {
	if e.aggFunc == nil {
		return []float64{0.0}
	}

	if reset {
		return e.aggFunc(nil, true)
	}

	if s != nil && e.MatchesSpan(s) {
		return e.aggFunc(s, false)
	}

	return nil
}

func (e *Expr) QueryType() QueryType {
	if e.stream == STREAM_TYPE_TRACES {
		return QueryTypeTraces
	}

	if e.aggFunc != nil {
		return QueryTypeMetrics
	}

	// if any matchers have descendant or parent fields then we require trace batching
	for _, m := range e.matchers {
		for _, field := range []field{m.lhs, m.rhs} {
			for _, f := range field.getRelationshipID() {
				if f == FIELD_DESCENDANT || f == FIELD_PARENT {
					return QueryTypeBatchedSpans
				}
			}

		}
	}

	return QueryTypeSpans
}

func matchesTraceField(m matcher, s *streampb.Span, t []*streampb.Span) bool {
	lhsSpans := getRelatedSpansForField(m.lhs, s, t)
	rhsSpans := getRelatedSpansForField(m.rhs, s, t)

	for _, lhs := range lhsSpans {
		for _, rhs := range rhsSpans {
			if m.compare(lhs, rhs) {
				return true
			}
		}
	}

	return false
}

func getRelatedSpansForField(f field, s *streampb.Span, t []*streampb.Span) []*streampb.Span {
	relID := f.getRelationshipID()
	if len(relID) == 0 {
		return []*streampb.Span{s}
	}

	if relID[0] == FIELD_PARENT {
		// todo: support multiple levels of parent
		if int(s.ParentIndex) < len(t) && s.ParentIndex >= 0 {
			return []*streampb.Span{t[s.ParentIndex]}
		}
	}

	if relID[0] == FIELD_DESCENDANT {
		// todo: improve this.  inefficient for deep traces
		descendants := []*streampb.Span{}
		parentIdx := s.ParentIndex
		for parentIdx >= 0 && int(parentIdx) < len(t) {
			s = t[parentIdx]

			descendants = append(descendants, s)
			parentIdx = s.ParentIndex
		}

		return descendants
	}

	return []*streampb.Span{}
}
