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
				"testString": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
				"testInt": &streampb.KeyValuePair{
					Type:     streampb.KeyValuePair_INT,
					IntValue: 3,
				},
				"testFloat": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_DOUBLE,
					DoubleValue: 3.14,
				},
				"testBool": &streampb.KeyValuePair{
					Type:      streampb.KeyValuePair_BOOL,
					BoolValue: true,
				},
			},
			Attributes: map[string]*streampb.KeyValuePair{
				"testString": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
				"testInt": &streampb.KeyValuePair{
					Type:     streampb.KeyValuePair_INT,
					IntValue: 3,
				},
				"testFloat": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_DOUBLE,
					DoubleValue: 3.14,
				},
				"testBool": &streampb.KeyValuePair{
					Type:      streampb.KeyValuePair_BOOL,
					BoolValue: true,
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
		&streampb.Span{
			Name:     "childSpan",
			Duration: 100,
			Events: map[string]*streampb.KeyValuePair{
				"testString": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
			},
			Attributes: map[string]*streampb.KeyValuePair{
				"testString": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
			},
			ParentIndex: 0,
			Process: &streampb.Process{
				Name: "proc1",
			},
			Status: &streampb.Status{
				Code:    12,
				Message: "status",
			},
		},
		&streampb.Span{
			Name:     "child2",
			Duration: 100,
			Events: map[string]*streampb.KeyValuePair{
				"testString": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
			},
			Attributes: map[string]*streampb.KeyValuePair{
				"testString": &streampb.KeyValuePair{
					Type:        streampb.KeyValuePair_STRING,
					StringValue: "test2",
				},
			},
			ParentIndex: 1,
			Process: &streampb.Process{
				Name: "proc2",
			},
			Status: &streampb.Status{
				Code:    12,
				Message: "status",
			},
		},
		&streampb.Span{
			Name:        "noparent",
			Duration:    100,
			Events:      map[string]*streampb.KeyValuePair{},
			Attributes:  map[string]*streampb.KeyValuePair{},
			ParentIndex: 1000,
			Process:     &streampb.Process{},
			Status:      &streampb.Status{},
		},
		&streampb.Span{
			Name:        "noparent2",
			Duration:    3,
			Events:      map[string]*streampb.KeyValuePair{},
			Attributes:  map[string]*streampb.KeyValuePair{},
			ParentIndex: -1,
			Process:     &streampb.Process{},
			Status:      &streampb.Status{},
		},
	}

	for _, tc := range []struct {
		in           string
		matchesSpans []int
	}{
		{
			in:           `spans{}`,
			matchesSpans: []int{0, 1, 2, 3, 4},
		},
		{
			in:           `spans{duration=3}`,
			matchesSpans: []int{4},
		},
		{
			in:           `spans{duration > 5}`,
			matchesSpans: []int{0, 1, 2, 3},
		},
		{
			in:           `spans{duration < 5}`,
			matchesSpans: []int{4},
		},
		{
			in:           `spans{name = "rootSpan"}`,
			matchesSpans: []int{0},
		},
		{
			in:           `spans{name != "rootSpan"}`,
			matchesSpans: []int{1, 2, 3, 4},
		},
		{
			in:           `spans{name =~ ".*Span"}`,
			matchesSpans: []int{0, 1},
		},
		{
			in:           `spans{name !~ ".*Span"}`,
			matchesSpans: []int{2, 3, 4},
		},
		{
			in:           `spans{atts.testString = "test2"}`,
			matchesSpans: []int{0, 1, 2},
		},
		{
			in:           `spans{atts.testString = 0}`,
			matchesSpans: []int{},
		},
		{
			in:           `spans{atts.testFloat = 3.14}`,
			matchesSpans: []int{0},
		},
		{
			in:           `spans{atts.testBool = 1}`,
			matchesSpans: []int{0},
		},
		{
			in:           `spans{atts.testInt = 1}`,
			matchesSpans: []int{},
		},
		{
			in:           `spans{atts.blerg = 0}`,
			matchesSpans: []int{},
		},
		{
			in:           `spans{events.testString > "abc"}`,
			matchesSpans: []int{0, 1, 2},
		},
		{
			in:           `spans{events.testFloat = 3.14}`,
			matchesSpans: []int{0},
		},
		{
			in:           `spans{events.testBool = 1}`,
			matchesSpans: []int{0},
		},
		{
			in:           `spans{events.testInt = "1"}`,
			matchesSpans: []int{},
		},
		{
			in:           `spans{events.blerg = 0}`,
			matchesSpans: []int{},
		},
		{
			in:           `spans{status.code != 12}`,
			matchesSpans: []int{3, 4},
		},
		{
			in:           `spans{status.message = "status"}`,
			matchesSpans: []int{0, 1, 2},
		},
		{
			in:           `spans{process.name = "proc1"}`,
			matchesSpans: []int{0, 1},
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
