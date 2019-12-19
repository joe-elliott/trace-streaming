package traceql

import "strconv"

const (
	opEQ int = iota
	opNEQ
	opRE
	opNRE
	opGT
	opGTE
	opLT
	opLTE

	streamSpans int = iota

	fieldDuration int = iota
	fieldName
	fieldTags
)

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

//
type ValueOperator interface {
	compareInt(int) bool
	compareFloat(float64) bool
	compareString(string) bool
}

type intOperator struct {
	val   int
	op    int
	field int
}

func newIntOperator(val int, op int, field int) intOperator {
	return intOperator{
		val:   val,
		op:    op,
		field: field,
	}
}

func (m intOperator) compareInt(n int) bool {
	return m.val == n
}

func (m intOperator) compareFloat(f float64) bool {
	return float64(m.val) == f
}

func (m intOperator) compareString(s string) bool {
	n, err := strconv.Atoi(s)

	if err != nil {
		return false
	}

	return m.val == n
}
