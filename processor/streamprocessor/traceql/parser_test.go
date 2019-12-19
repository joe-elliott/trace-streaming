package traceql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		in         string
		stream     int
		fieldIds   []int
		fieldNames []string
		err        error
	}{
		{
			in:         `spans{duration=3, name="asdf"}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []int{FIELD_DURATION, FIELD_NAME},
			fieldNames: []string{"", ""},
		},
		{
			in:         `spans{duration=3, atts.test="blerg"}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []int{FIELD_DURATION, FIELD_ATTS},
			fieldNames: []string{"", "test"},
		},
		{
			in:         `spans{duration=3, atts.test="blerg", status.message=~".*blerg", status.code=400}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []int{FIELD_DURATION, FIELD_ATTS, FIELD_STATUS_MSG, FIELD_STATUS_CODE},
			fieldNames: []string{"", "test", "", ""},
		},
		{
			in: `spans{foo="bar"}`,
			err: ParseError{
				msg:  "syntax error: unexpected IDENTIFIER",
				line: 1,
				col:  7,
			},
		},
		{
			in: `blerg{foo="bar"}`,
			err: ParseError{
				msg:  "syntax error: unexpected IDENTIFIER, expecting STREAM_TYPE_SPANS",
				line: 1,
				col:  1,
			},
		},
		{
			in: `spans{status.grub="bar"}`,
			err: ParseError{
				msg:  "syntax error: unexpected IDENTIFIER, expecting FIELD_STATUS_CODE or FIELD_STATUS_MSG",
				line: 1,
				col:  14,
			},
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			expr, err := ParseExpr(tc.in)

			assert.Equal(t, tc.err, err)

			if tc.stream != 0 {
				assert.Equal(t, tc.stream, expr.stream)

				for i, o := range expr.operators {
					var fieldID int
					var fieldName string

					switch v := o.(type) {
					case intOperator:
						fieldID = v.field.fieldID
						fieldName = v.field.fieldName
					case floatOperator:
						fieldID = v.field.fieldID
						fieldName = v.field.fieldName
					case stringOperator:
						fieldID = v.field.fieldID
						fieldName = v.field.fieldName
					default:
						assert.Failf(t, "", "Unkown type %T", v)
					}

					assert.Equal(t, tc.fieldIds[i], fieldID)
					assert.Equal(t, tc.fieldNames[i], fieldName)
				}
			}
		})
	}
}
