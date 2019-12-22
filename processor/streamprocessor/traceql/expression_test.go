package traceql

import (
	"testing"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
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

func TestMatchesSpan(t *testing.T) {
	trace := []*streampb.Span{
		&streampb.Span{
			Name:     "rootSpan",
			Duration: 100,
			Events: map[string]*streampb.KeyValuePair{
				"test": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
			},
			Attributes: map[string]*streampb.KeyValuePair{
				"test": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
			},
			ParentIndex: -1,
			Process: &streampb.Process{
				Name: "proc1",
			},
			Status: &streampb.Status{
				Code:    12,
				Message: "status",
			},
		},
	}

	for _, tc := range []struct {
		in           string
		matchesSpans []int
	}{
		{
			in:           `spans{}`,
			matchesSpans: []int{0},
		},
		{
			in:           `spans{duration=3, name="asdf"}`,
			matchesSpans: []int{},
		},
		{
			in:           `spans{duration < 5}`,
			matchesSpans: []int{0},
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			expr, err := ParseExpr(tc.in)

			assert.Nil(t, err)
			if expr == nil {
				assert.FailNow(t, "expr is unexpectedly nil.")
			}

			for i, span := range trace {
				assert.Equalf(t, contains(tc.matchesSpans, i), expr.MatchesSpan(span, trace), "failed for span %d", i)
			}
		})
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
