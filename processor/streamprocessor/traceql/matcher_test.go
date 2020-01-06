package traceql

import (
	"testing"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	span := &streampb.Span{
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
		ParentIndex:  -1,
		ParentSpanID: []byte{},
		Process: &streampb.Process{
			Name: "proc1",
		},
		Status: &streampb.Status{
			Code:    12,
			Message: "status",
		},
	}

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
		{
			in:       newMatcher(intField(2), GTE, wrapDynamicField(FIELD_STATUS, newDynamicField(FIELD_CODE, ""))),
			expected: false,
		},
		{
			in:       newMatcher(intField(2), LT, wrapDynamicField(FIELD_STATUS, newDynamicField(FIELD_CODE, ""))),
			expected: true,
		},
		{
			in:       newMatcher(stringField("status"), NEQ, wrapDynamicField(FIELD_STATUS, newDynamicField(FIELD_MSG, ""))),
			expected: false,
		},
		{
			in:       newMatcher(stringField("status"), LTE, wrapDynamicField(FIELD_STATUS, newDynamicField(FIELD_MSG, ""))),
			expected: true,
		},
		{
			in:       newMatcher(floatField(3), LT, newDynamicField(FIELD_EVENTS, "testFloat")),
			expected: true,
		},
		{
			in:       newMatcher(floatField(3), GT, newDynamicField(FIELD_EVENTS, "testFloat")),
			expected: false,
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.in.compare(span, span))
		})
	}
}
