package traceql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// jpe - add rhs tests
// jpe - add static field tests
func TestParse(t *testing.T) {
	for _, tc := range []struct {
		in            string
		stream        int
		lhsFieldIds   []fieldID
		lhsFieldNames []string
		rhsFieldIds   []fieldID
		rhsFieldNames []string
		err           error
	}{
		{
			in:            `spans{}`,
			stream:        STREAM_TYPE_SPANS,
			lhsFieldIds:   []fieldID{},
			lhsFieldNames: []string{},
			rhsFieldIds:   []fieldID{},
			rhsFieldNames: []string{},
		},
		{
			in:            `spans{duration=3, name="asdf"}`,
			stream:        STREAM_TYPE_SPANS,
			lhsFieldIds:   []fieldID{[]int{FIELD_DURATION}, []int{FIELD_NAME}},
			lhsFieldNames: []string{"", ""},
			rhsFieldIds:   []fieldID{[]int{}, []int{}},
			rhsFieldNames: []string{"", ""},
		},
		{
			in:            `spans{duration=3, atts.test="blerg"}`,
			stream:        STREAM_TYPE_SPANS,
			lhsFieldIds:   []fieldID{[]int{FIELD_DURATION}, []int{FIELD_ATTS}},
			lhsFieldNames: []string{"", "test"},
			rhsFieldIds:   []fieldID{[]int{}, []int{}},
			rhsFieldNames: []string{"", ""},
		},
		{
			in:            `spans{duration=3, atts.test="blerg", status.message=~".*blerg", status.code=400}`,
			stream:        STREAM_TYPE_SPANS,
			lhsFieldIds:   []fieldID{[]int{FIELD_DURATION}, []int{FIELD_ATTS}, []int{FIELD_STATUS, FIELD_MSG}, []int{FIELD_STATUS, FIELD_CODE}},
			lhsFieldNames: []string{"", "test", "", ""},
			rhsFieldIds:   []fieldID{[]int{}, []int{}, []int{}, []int{}},
			rhsFieldNames: []string{"", "", "", ""},
		},
		{
			in:            `spans{parent*.duration=3}`,
			stream:        STREAM_TYPE_SPANS,
			lhsFieldIds:   []fieldID{[]int{FIELD_DESCENDANT, FIELD_DURATION}},
			lhsFieldNames: []string{""},
			rhsFieldIds:   []fieldID{[]int{}},
			rhsFieldNames: []string{""},
		},
		{
			in:            `spans{parent.parent.duration=3}`,
			stream:        STREAM_TYPE_SPANS,
			lhsFieldIds:   []fieldID{[]int{FIELD_PARENT, FIELD_PARENT, FIELD_DURATION}},
			lhsFieldNames: []string{""},
			rhsFieldIds:   []fieldID{[]int{}},
			rhsFieldNames: []string{""},
		},
		{
			in:            `spans{parent.parent.atts.test=3}`,
			stream:        STREAM_TYPE_SPANS,
			lhsFieldIds:   []fieldID{[]int{FIELD_PARENT, FIELD_PARENT, FIELD_ATTS}},
			lhsFieldNames: []string{"test"},
			rhsFieldIds:   []fieldID{[]int{}},
			rhsFieldNames: []string{""},
		},
		{
			in:            `traces{}`,
			stream:        STREAM_TYPE_TRACES,
			lhsFieldIds:   []fieldID{},
			lhsFieldNames: []string{},
			rhsFieldIds:   []fieldID{},
			rhsFieldNames: []string{},
		},
		{
			in:            `traces{duration = 3}`,
			stream:        STREAM_TYPE_TRACES,
			lhsFieldIds:   []fieldID{[]int{FIELD_DURATION}},
			lhsFieldNames: []string{""},
			rhsFieldIds:   []fieldID{[]int{}},
			rhsFieldNames: []string{""},
		},
		{
			in:            `traces{isRoot = 1}`,
			stream:        STREAM_TYPE_TRACES,
			lhsFieldIds:   []fieldID{[]int{FIELD_IS_ROOT}},
			lhsFieldNames: []string{""},
			rhsFieldIds:   []fieldID{[]int{}},
			rhsFieldNames: []string{""},
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
				fRHS := o.rhs
				fLHS := o.lhs

				dLHS, ok := fLHS.(dynamicField)
				if ok {
					assert.Equalf(t, tc.lhsFieldIds[i], append(dLHS.relID, dLHS.id...), "lhs actual %v", o)
					assert.Equalf(t, tc.lhsFieldNames[i], dLHS.name, "lhs actual %v", o)
				}

				dRHS, ok := fRHS.(dynamicField)
				if ok {
					assert.Equalf(t, tc.rhsFieldIds[i], append(dRHS.relID, dRHS.id...), "rhs actual %v", o)
					assert.Equalf(t, tc.rhsFieldNames[i], dRHS.name, "rhs actual %v", o)
				}
			}
		})
	}
}
