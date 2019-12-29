package traceql

import (
	"regexp"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

type matcher struct {
	lhs field
	rhs field
	op  int
}

func newMatcher(lhs field, op int, rhs field) matcher {
	return matcher{
		lhs: lhs,
		rhs: rhs,
		op:  op,
	}
}

func (m matcher) compare(lhsSpan *streampb.Span, rhsSpan *streampb.Span) bool {
	t1 := m.lhs.getNativeType(lhsSpan)
	t2 := m.rhs.getNativeType(rhsSpan)

	if t1 != t2 {
		// todo: bubble this error up/metric
		return false
	}

	switch t1 {
	case fieldTypeInt:
		return m.compareInt(lhsSpan, rhsSpan)
	case fieldTypeFloat:
		return m.compareFloat(lhsSpan, rhsSpan)
	case fieldTypeString:
		return m.compareString(lhsSpan, rhsSpan)
	}

	return false
}

func (m matcher) compareInt(lhsSpan *streampb.Span, rhsSpan *streampb.Span) bool {
	n1 := m.lhs.getIntValue(lhsSpan)
	n2 := m.rhs.getIntValue(rhsSpan)

	switch m.op {
	case EQ:
		return n1 == n2
	case NEQ:
		return n1 != n2
	case GT:
		return n1 > n2
	case GTE:
		return n1 >= n2
	case LT:
		return n1 < n2
	case LTE:
		return n1 <= n2
	default:
		return false
	}
}

func (m matcher) compareFloat(lhsSpan *streampb.Span, rhsSpan *streampb.Span) bool {
	f1 := m.lhs.getFloatValue(lhsSpan)
	f2 := m.rhs.getFloatValue(rhsSpan)

	switch m.op {
	case EQ:
		return f1 == f2
	case NEQ:
		return f1 != f2
	case GT:
		return f1 > f2
	case GTE:
		return f1 >= f2
	case LT:
		return f1 < f2
	case LTE:
		return f1 <= f2
	default:
		return false
	}
}

func (m matcher) compareString(lhsSpan *streampb.Span, rhsSpan *streampb.Span) bool {
	s1 := m.lhs.getStringValue(lhsSpan)
	s2 := m.rhs.getStringValue(rhsSpan)

	switch m.op {
	case EQ:
		return s1 == s2
	case NEQ:
		return s1 != s2
	case RE:
		// todo: fix serious performance issues here.  compiling regex everytime. string -> byte conversion
		regex := regexp.MustCompile(s2)
		return regex.Match([]byte(s1))
	case NRE:
		regex := regexp.MustCompile(s2)
		return !regex.Match([]byte(s1))
	case GT:
		return s1 > s2
	case GTE:
		return s1 >= s2
	case LT:
		return s1 < s2
	case LTE:
		return s1 <= s2
	default:
		return false
	}
}
