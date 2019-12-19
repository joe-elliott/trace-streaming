package traceql

//
type Expr struct {
	stream    int
	operators []ValueOperator
}

func newExpr(stream int, o []ValueOperator) *Expr {
	return &Expr{
		stream:    stream,
		operators: o,
	}
}
