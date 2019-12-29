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

// jpe - add field tests
// jpe - add dynamic field matcher tests
func TestMatcher(t *testing.T) {
	for _, tc := range []struct {
		in       matcher
		expected bool
	}{
		{
			in:       newMatcher(intField(3), EQ, intField(3)),
			expected: true,
		},
		{
			in:       newMatcher(floatField(3), EQ, intField(3)),
			expected: false,
		},
		{
			in:       newMatcher(floatField(3), EQ, floatField(3)),
			expected: true,
		},
		{
			in:       newMatcher(stringField("asdf"), RE, stringField(".*")),
			expected: true,
		},
		{
			in:       newMatcher(intField(2), LT, intField(3)),
			expected: true,
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.in.compare(nil, nil))
		})
	}
}
