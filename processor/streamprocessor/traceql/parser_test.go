package traceql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		in         string
		stream     int
		fieldIds   []fieldID
		fieldNames []string
		err        error
	}{
		{
			in:         `spans{}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []fieldID{},
			fieldNames: []string{},
		},
		{
			in:         `spans{duration=3, name="asdf"}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []fieldID{[]int{FIELD_DURATION}, []int{FIELD_NAME}},
			fieldNames: []string{"", ""},
		},
		{
			in:         `spans{duration=3, atts.test="blerg"}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []fieldID{[]int{FIELD_DURATION}, []int{FIELD_ATTS}},
			fieldNames: []string{"", "test"},
		},
		{
			in:         `spans{duration=3, atts.test="blerg", status.message=~".*blerg", status.code=400}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []fieldID{[]int{FIELD_DURATION}, []int{FIELD_ATTS}, []int{FIELD_STATUS, FIELD_MSG}, []int{FIELD_STATUS, FIELD_CODE}},
			fieldNames: []string{"", "test", "", ""},
		},
		{
			in:         `spans{parent*.duration=3}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []fieldID{[]int{FIELD_DESCENDANT, FIELD_DURATION}},
			fieldNames: []string{""},
		},
		{
			in:         `spans{parent.parent.duration=3}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []fieldID{[]int{FIELD_PARENT, FIELD_PARENT, FIELD_DURATION}},
			fieldNames: []string{""},
		},
		{
			in:         `spans{parent.parent.atts.test=3}`,
			stream:     STREAM_TYPE_SPANS,
			fieldIds:   []fieldID{[]int{FIELD_PARENT, FIELD_PARENT, FIELD_ATTS}},
			fieldNames: []string{"test"},
		},
		{
			in:         `traces{}`,
			stream:     STREAM_TYPE_TRACES,
			fieldIds:   []fieldID{},
			fieldNames: []string{},
		},
		{
			in:         `traces{duration = 3}`,
			stream:     STREAM_TYPE_TRACES,
			fieldIds:   []fieldID{[]int{FIELD_DURATION}},
			fieldNames: []string{""},
		},
		{
			in:         `traces{isRoot = 1}`,
			stream:     STREAM_TYPE_TRACES,
			fieldIds:   []fieldID{[]int{FIELD_IS_ROOT}},
			fieldNames: []string{""},
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
				msg:  "syntax error: unexpected IDENTIFIER, expecting STREAM_TYPE_SPANS or STREAM_TYPE_TRACES",
				line: 1,
				col:  1,
			},
		},
		{
			in: `spans{status.grub="bar"}`,
			err: ParseError{
				msg:  "syntax error: unexpected IDENTIFIER, expecting FIELD_CODE or FIELD_MSG",
				line: 1,
				col:  14,
			},
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			expr, err := ParseExpr(tc.in)

			assert.Equal(t, tc.err, err)

			// if we were expecting an error, just bail out at this point
			if tc.err != nil {
				return
			}

			if expr == nil {
				assert.FailNow(t, "expr is unexpectedly nil.")
			}

			assert.Equal(t, tc.stream, expr.stream)

			for i, o := range expr.matchers {
				f := o.field()

				assert.Equalf(t, tc.fieldIds[i], f.id, "actual %v", o)
				assert.Equalf(t, tc.fieldNames[i], f.name, "actual %v", o)
			}
		})
	}
}
