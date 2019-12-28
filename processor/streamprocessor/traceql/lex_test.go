package traceql

import (
	"strings"
	"testing"
	"text/scanner"

	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected []int
	}{
		{`spans{name="test"}`, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, EQ, STRING, CLOSE_BRACE}},
		{` spans{ name = "test" } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, EQ, STRING, CLOSE_BRACE}},
		{` spans{ name > "test" } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, GT, STRING, CLOSE_BRACE}},
		{` spans{ name >= "test" } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, GTE, STRING, CLOSE_BRACE}},
		{` spans{ name < "test" } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, LT, STRING, CLOSE_BRACE}},
		{` spans{ name <= "test" } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, LTE, STRING, CLOSE_BRACE}},
		{` spans{ name =~ "test" } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, RE, STRING, CLOSE_BRACE}},
		{` spans{ name !~ "test" } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, NRE, STRING, CLOSE_BRACE}},
		{` spans{ name <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ duration <= 1.21 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DURATION, LTE, FLOAT, CLOSE_BRACE}},
		{` spans{ atts.thing <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_ATTS, DOT, IDENTIFIER, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ events.thing <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_EVENTS, DOT, IDENTIFIER, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ status.code <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_STATUS, DOT, FIELD_CODE, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ status.message <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_STATUS, DOT, FIELD_MSG, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ parent*.duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DESCENDANT, DOT, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ parent.duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_PARENT, DOT, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ parent.parent.duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_PARENT, DOT, FIELD_PARENT, DOT, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ process.name <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_PROCESS, DOT, FIELD_NAME, LTE, INTEGER, CLOSE_BRACE}},
		{` traces{ process.name <= 13 } `, []int{STREAM_TYPE_TRACES, OPEN_BRACE, FIELD_PROCESS, DOT, FIELD_NAME, LTE, INTEGER, CLOSE_BRACE}},
		{` traces{ isRoot = 1 } `, []int{STREAM_TYPE_TRACES, OPEN_BRACE, FIELD_IS_ROOT, EQ, INTEGER, CLOSE_BRACE}},
	} {
		t.Run(tc.input, func(t *testing.T) {
			actual := []int{}
			l := lexer{
				Scanner: scanner.Scanner{
					Mode: scanner.SkipComments | scanner.ScanStrings,
				},
			}
			l.Init(strings.NewReader(tc.input))
			var lval yySymType
			for {
				tok := l.Lex(&lval)
				if tok == 0 {
					break
				}
				actual = append(actual, tok)
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}
