package traceql

import "github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"

type Evaluator interface {
	MatchesSpan(*streampb.Span) bool
	MatchesTrace([]*streampb.Span) bool
}

//
type Expr struct {
	stream   int
	matchers []ValueMatcher
}

func newExpr(stream int, m []ValueMatcher) *Expr {
	return &Expr{
		stream:   stream,
		matchers: m,
	}
}

func (e *Expr) MatchesSpan(s *streampb.Span) bool {

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
