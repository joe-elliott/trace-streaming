package traceql

import (
	"regexp"
	"strconv"
)

type fieldID []int

func (f fieldID) isRoot() bool {
	if len(f) == 0 {
		return false
	}

	return f[0] == FIELD_ROOT_SPAN
}

type intCompareFunc func(int) bool
type floatCompareFunc func(float64) bool
type stringCompareFunc func(string) bool

// complexField
type complexField struct {
	id   fieldID
	name string // only valid if id = FIELD_ATTS or FIELD_EVENTS
}

func newComplexField(id int, name string) complexField {
	return complexField{
		id:   []int{id},
		name: name,
	}
}

func wrapComplexField(id int, c complexField) complexField {
	return complexField{
		id:   append([]int{id}, c.id...),
		name: c.name,
	}
}

//
type ValueMatcher interface {
	compareInt(int) bool
	compareFloat(float64) bool
	compareString(string) bool

	field() complexField
}

// int Matcher
type intMatcher struct {
	f       complexField
	compare intCompareFunc
}

func newIntMatcher(val int, op int, field complexField) intMatcher {
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

	return intMatcher{
		f:       field,
		compare: compare,
	}
}

func (o intMatcher) compareInt(n int) bool {
	return o.compare(n)
}

func (o intMatcher) compareFloat(f float64) bool {
	// comparing a float to an int is just going to floor the float
	return o.compare(int(f))
}

func (o intMatcher) compareString(s string) bool {
	n, err := strconv.Atoi(s)

	if err != nil {
		return false
	}

	return o.compare(n)
}

func (o intMatcher) field() complexField {
	return o.f
}

// float Matcher
type floatMatcher struct {
	f       complexField
	compare floatCompareFunc
}

func newFloatMatcher(val float64, op int, field complexField) floatMatcher {
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

	return floatMatcher{
		f:       field,
		compare: compare,
	}
}

func (o floatMatcher) compareInt(n int) bool {
	return o.compare(float64(n))
}

func (o floatMatcher) compareFloat(f float64) bool {
	return o.compare(f)
}

func (o floatMatcher) compareString(s string) bool {
	f, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return false
	}

	return o.compare(f)
}

func (o floatMatcher) field() complexField {
	return o.f
}

// string Matcher
type stringMatcher struct {
	f       complexField
	compare stringCompareFunc
}

func newStringMatcher(val string, op int, field complexField) stringMatcher {
	var compare stringCompareFunc
	// if op is a regex, let's build the regex now

	switch op {
	case EQ:
		compare = func(n string) bool { return n == val }
	case NEQ:
		compare = func(n string) bool { return n != val }
	case RE:
		regex := regexp.MustCompile(val)
		compare = func(n string) bool { return regex.Match([]byte(n)) } // todo - consider performance.  should we use strings at all or only []byte?
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

	return stringMatcher{
		f:       field,
		compare: compare,
	}
}

func (o stringMatcher) compareInt(n int) bool {
	return o.compare(strconv.FormatInt(int64(n), 10))
}

func (o stringMatcher) compareFloat(f float64) bool {
	return o.compare(strconv.FormatFloat(f, 'E', 1, 64))
}

func (o stringMatcher) compareString(s string) bool {
	return o.compare(s)
}

func (o stringMatcher) field() complexField {
	return o.f
}
