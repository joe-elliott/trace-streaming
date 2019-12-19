package traceql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type intTest struct {
	compare  int
	expected bool
}
type floatTest struct {
	compare  float64
	expected bool
}
type stringTest struct {
	compare  string
	expected bool
}

func TestIntOperator(t *testing.T) {
	for _, tc := range []struct {
		in intOperator
		i  intTest
		f  floatTest
		s  stringTest
	}{
		{
			in: newIntOperator(0, EQ, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newIntOperator(0, NEQ, complexField{}),
			i:  intTest{0, false},
			f:  floatTest{0, false},
			s:  stringTest{"0", false},
		},
		{
			in: newIntOperator(0, GT, complexField{}),
			i:  intTest{3, true},
			f:  floatTest{3, true},
			s:  stringTest{"3", true},
		},
		{
			in: newIntOperator(0, GTE, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newIntOperator(0, LT, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newIntOperator(0, LTE, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.Equalf(t, tc.i.expected, tc.in.compareInt(tc.i.compare), "int test")
			assert.Equalf(t, tc.f.expected, tc.in.compareFloat(tc.f.compare), "float test")
			assert.Equalf(t, tc.s.expected, tc.in.compareString(tc.s.compare), "string test")
		})
	}

}

func TestFloatOperator(t *testing.T) {
	for _, tc := range []struct {
		in floatOperator
		i  intTest
		f  floatTest
		s  stringTest
	}{
		{
			in: newFloatOperator(0, EQ, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newFloatOperator(0, NEQ, complexField{}),
			i:  intTest{0, false},
			f:  floatTest{0, false},
			s:  stringTest{"0", false},
		},
		{
			in: newFloatOperator(0, GT, complexField{}),
			i:  intTest{3, true},
			f:  floatTest{3, true},
			s:  stringTest{"3", true},
		},
		{
			in: newFloatOperator(0, GTE, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newFloatOperator(0, LT, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newFloatOperator(0, LTE, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.Equalf(t, tc.i.expected, tc.in.compareInt(tc.i.compare), "int test")
			assert.Equalf(t, tc.f.expected, tc.in.compareFloat(tc.f.compare), "float test")
			assert.Equalf(t, tc.s.expected, tc.in.compareString(tc.s.compare), "string test")
		})
	}
}

func TestStringOperator(t *testing.T) {
	for _, tc := range []struct {
		in stringOperator
		i  intTest
		f  floatTest
		s  stringTest
	}{
		{
			in: newStringOperator("0", EQ, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, false},
			s:  stringTest{"0", true},
		},
		{
			in: newStringOperator("0", NEQ, complexField{}),
			i:  intTest{0, false},
			f:  floatTest{0, true},
			s:  stringTest{"0", false},
		},
		{
			in: newStringOperator("0", GT, complexField{}),
			i:  intTest{3, true},
			f:  floatTest{3, true},
			s:  stringTest{"3", true},
		},
		{
			in: newStringOperator("0", GTE, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newStringOperator("0", LT, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newStringOperator("0", LTE, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newStringOperator(".*blerg", RE, complexField{}),
			i:  intTest{-3, false},
			f:  floatTest{-3, false},
			s:  stringTest{"thingblerg", true},
		},
		{
			in: newStringOperator(".*blerg", NRE, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"thingblerg", false},
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.Equalf(t, tc.i.expected, tc.in.compareInt(tc.i.compare), "int test")
			assert.Equalf(t, tc.f.expected, tc.in.compareFloat(tc.f.compare), "float test")
			assert.Equalf(t, tc.s.expected, tc.in.compareString(tc.s.compare), "string test")
		})
	}

}
