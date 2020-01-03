package traceql

import (
	"testing"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
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
		in             field
		expectedType   int
		expectedInt    int
		expectedFloat  float64
		expectedString string
		expectedRelID  fieldID
	}{
		{
			in:            intField(3),
			expectedType:  fieldTypeInt,
			expectedInt:   3,
			expectedFloat: 3,
		},
		{
			in:            floatField(3),
			expectedType:  fieldTypeFloat,
			expectedFloat: 3,
		},
		{
			in:             stringField("3"),
			expectedType:   fieldTypeString,
			expectedString: "3",
		},
		{
			in:             newDynamicField(FIELD_NAME, ""),
			expectedType:   fieldTypeString,
			expectedString: "rootSpan",
		},
		{
			in:            newDynamicField(FIELD_DURATION, ""),
			expectedType:  fieldTypeInt,
			expectedInt:   100,
			expectedFloat: 100,
		},
		{
			in:             newDynamicField(FIELD_EVENTS, "testString"),
			expectedType:   fieldTypeString,
			expectedString: "test2",
		},
		{
			in:           newDynamicField(FIELD_EVENTS, "doesntexist"),
			expectedType: fieldTypeUnknown,
		},
		{
			in:           newDynamicField(FIELD_ATTS, "testBool"),
			expectedType: fieldTypeInt,
			expectedInt:  1,
		},
		{
			in:             wrapDynamicField(FIELD_PROCESS, newDynamicField(FIELD_NAME, "")),
			expectedType:   fieldTypeString,
			expectedString: "proc1",
		},
		{
			in:             wrapRelationshipField(FIELD_PARENT, wrapDynamicField(FIELD_PROCESS, newDynamicField(FIELD_NAME, ""))),
			expectedType:   fieldTypeString,
			expectedString: "proc1",
			expectedRelID:  []int{FIELD_PARENT},
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tc.expectedType, tc.in.getNativeType(span))
			assert.Equal(t, tc.expectedInt, tc.in.getIntValue(span))
			assert.Equal(t, tc.expectedString, tc.in.getStringValue(span))
			assert.Equal(t, tc.expectedFloat, tc.in.getFloatValue(span))
			assert.Equal(t, tc.expectedRelID, tc.in.getRelationshipID())
		})
	}
}
