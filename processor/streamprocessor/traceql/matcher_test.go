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

func TestIntMatcher(t *testing.T) {
	for _, tc := range []struct {
		in intMatcher
		i  intTest
		f  floatTest
		s  stringTest
	}{
		{
			in: newIntMatcher(0, EQ, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newIntMatcher(0, NEQ, complexField{}),
			i:  intTest{0, false},
			f:  floatTest{0, false},
			s:  stringTest{"0", false},
		},
		{
			in: newIntMatcher(0, GT, complexField{}),
			i:  intTest{3, true},
			f:  floatTest{3, true},
			s:  stringTest{"3", true},
		},
		{
			in: newIntMatcher(0, GTE, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newIntMatcher(0, LT, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newIntMatcher(0, LTE, complexField{}),
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

func TestFloatMatcher(t *testing.T) {
	for _, tc := range []struct {
		in floatMatcher
		i  intTest
		f  floatTest
		s  stringTest
	}{
		{
			in: newFloatMatcher(0, EQ, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newFloatMatcher(0, NEQ, complexField{}),
			i:  intTest{0, false},
			f:  floatTest{0, false},
			s:  stringTest{"0", false},
		},
		{
			in: newFloatMatcher(0, GT, complexField{}),
			i:  intTest{3, true},
			f:  floatTest{3, true},
			s:  stringTest{"3", true},
		},
		{
			in: newFloatMatcher(0, GTE, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newFloatMatcher(0, LT, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newFloatMatcher(0, LTE, complexField{}),
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

func TestStringMatcher(t *testing.T) {
	for _, tc := range []struct {
		in stringMatcher
		i  intTest
		f  floatTest
		s  stringTest
	}{
		{
			in: newStringMatcher("0", EQ, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, false},
			s:  stringTest{"0", true},
		},
		{
			in: newStringMatcher("0", NEQ, complexField{}),
			i:  intTest{0, false},
			f:  floatTest{0, true},
			s:  stringTest{"0", false},
		},
		{
			in: newStringMatcher("0", GT, complexField{}),
			i:  intTest{3, true},
			f:  floatTest{3, true},
			s:  stringTest{"3", true},
		},
		{
			in: newStringMatcher("0", GTE, complexField{}),
			i:  intTest{0, true},
			f:  floatTest{0, true},
			s:  stringTest{"0", true},
		},
		{
			in: newStringMatcher("0", LT, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newStringMatcher("0", LTE, complexField{}),
			i:  intTest{-3, true},
			f:  floatTest{-3, true},
			s:  stringTest{"-3", true},
		},
		{
			in: newStringMatcher(".*blerg", RE, complexField{}),
			i:  intTest{-3, false},
			f:  floatTest{-3, false},
			s:  stringTest{"thingblerg", true},
		},
		{
			in: newStringMatcher(".*blerg", NRE, complexField{}),
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
