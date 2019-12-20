package traceql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiresTraceBatching(t *testing.T) {
	for _, tc := range []struct {
		in       string
		expected bool
	}{
		{
			in:       `spans{}`,
			expected: false,
		},
		{
			in:       `spans{duration=3, name="asdf"}`,
			expected: false,
		},
		{
			in:       `spans{duration=3, atts.test="blerg"}`,
			expected: false,
		},
		{
			in:       `spans{duration=3, atts.test="blerg", status.message=~".*blerg", status.code=400}`,
			expected: false,
		},
		{
			in:       `spans{parent*.duration=3}`,
			expected: true,
		},
		{
			in:       `spans{parent.parent.duration=3}`,
			expected: true,
		},
		{
			in:       `spans{parent.atts.test=3}`,
			expected: true,
		},
		{
			in:       `traces{}`,
			expected: true,
		},
		{
			in:       `traces{span.duration = 3}`,
			expected: true,
		},
		{
			in:       `traces{rootSpan.parent.duration = 3}`,
			expected: true,
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			expr, err := ParseExpr(tc.in)

			assert.Nil(t, err)
			if expr == nil {
				assert.FailNow(t, "expr is unexpectedly nil.")
			}

			assert.Equal(t, tc.expected, expr.RequiresTraceBatching())
		})
	}
}
