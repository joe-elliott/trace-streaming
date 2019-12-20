package traceql

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
