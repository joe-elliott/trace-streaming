package traceql

import (
	"regexp"
	"strconv"
)

type intCompareFunc func(int) bool
type floatCompareFunc func(float64) bool
type stringCompareFunc func(string) bool

// complexField
type complexField struct {
	fieldID   int
	fieldName string // only valid if fieldID = FIELD_TAGS
}

func newComplexField(id int, name string) complexField {
	return complexField{
		fieldID:   id,
		fieldName: name,
	}
}

//
type ValueOperator interface {
	compareInt(int) bool
	compareFloat(float64) bool
	compareString(string) bool
}

// int operator
type intOperator struct {
	field   complexField
	compare intCompareFunc
}

func newIntOperator(val int, op int, field complexField) intOperator {
	var compare intCompareFunc

	switch op {
	case EQ:
		compare = func(n int) bool { return n == val }
	case NEQ:
		compare = func(n int) bool { return n != val }
	case GT:
		compare = func(n int) bool { return n > val }
	case GTE:
		compare = func(n int) bool { return n >= val }
	case LT:
		compare = func(n int) bool { return n < val }
	case LTE:
		compare = func(n int) bool { return n <= val }
	default:
		compare = func(n int) bool { return false }
	}

	return intOperator{
		field:   field,
		compare: compare,
	}
}

func (o intOperator) compareInt(n int) bool {
	return o.compare(n)
}

func (o intOperator) compareFloat(f float64) bool {
	// comparing a float to an int is just going to floor the float
	return o.compare(int(f))
}

func (o intOperator) compareString(s string) bool {
	n, err := strconv.Atoi(s)

	if err != nil {
		return false
	}

	return o.compare(n)
}

// float operator
type floatOperator struct {
	field   complexField
	compare floatCompareFunc
}

func newFloatOperator(val float64, op int, field complexField) floatOperator {
	var compare floatCompareFunc

	switch op {
	case EQ:
		compare = func(n float64) bool { return n == val }
	case NEQ:
		compare = func(n float64) bool { return n != val }
	case GT:
		compare = func(n float64) bool { return n > val }
	case GTE:
		compare = func(n float64) bool { return n >= val }
	case LT:
		compare = func(n float64) bool { return n < val }
	case LTE:
		compare = func(n float64) bool { return n <= val }
	default:
		compare = func(n float64) bool { return false }
	}

	return floatOperator{
		field:   field,
		compare: compare,
	}
}

func (o floatOperator) compareInt(n int) bool {
	return o.compare(float64(n))
}

func (o floatOperator) compareFloat(f float64) bool {
	return o.compare(f)
}

func (o floatOperator) compareString(s string) bool {
	f, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return false
	}

	return o.compare(f)
}

// string operator
type stringOperator struct {
	field   complexField
	compare stringCompareFunc
}

func newStringOperator(val string, op int, field complexField) stringOperator {
	var compare stringCompareFunc
	// if op is a regex, let's build the regex now

	switch op {
	case EQ:
		compare = func(n string) bool { return n == val }
	case NEQ:
		compare = func(n string) bool { return n != val }
	case RE:
		regex := regexp.MustCompile(val)
		compare = func(n string) bool { return regex.Match([]byte(n)) } // jpe - consider performance.  should we use strings at all or only []byte?
	case NRE:
		regex := regexp.MustCompile(val)
		compare = func(n string) bool { return !regex.Match([]byte(n)) }
	case GT:
		compare = func(n string) bool { return n > val }
	case GTE:
		compare = func(n string) bool { return n >= val }
	case LT:
		compare = func(n string) bool { return n < val }
	case LTE:
		compare = func(n string) bool { return n <= val }
	default:
		compare = func(n string) bool { return false }
	}

	return stringOperator{
		field:   field,
		compare: compare,
	}
}

func (o stringOperator) compareInt(n int) bool {
	return o.compare(strconv.FormatInt(int64(n), 10))
}

func (o stringOperator) compareFloat(f float64) bool {
	return o.compare(strconv.FormatFloat(f, 'E', 1, 64))
}

func (o stringOperator) compareString(s string) bool {
	return o.compare(s)
}
