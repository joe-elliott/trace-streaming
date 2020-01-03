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
		{` spans{ "test" !~ name } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, STRING, NRE, FIELD_NAME, CLOSE_BRACE}},
		{` spans{ name <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ duration <= 1.21 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DURATION, LTE, FLOAT, CLOSE_BRACE}},
		{` spans{ atts["thing"] <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_ATTS, OPEN_BRACKET, STRING, CLOSE_BRACKET, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ events["thing"] <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_EVENTS, OPEN_BRACKET, STRING, CLOSE_BRACKET, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ 13 <= status.code } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, INTEGER, LTE, FIELD_STATUS, DOT, FIELD_CODE, CLOSE_BRACE}},
		{` spans{ status.message <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_STATUS, DOT, FIELD_MSG, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ parent*.duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DESCENDANT, DOT, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ parent.duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_PARENT, DOT, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ parent.parent.duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_PARENT, DOT, FIELD_PARENT, DOT, FIELD_DURATION, LTE, INTEGER, CLOSE_BRACE}},
		{` spans{ process.name <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_PROCESS, DOT, FIELD_NAME, LTE, INTEGER, CLOSE_BRACE}},
		{` traces{ process.name <= 13 } `, []int{STREAM_TYPE_TRACES, OPEN_BRACE, FIELD_PROCESS, DOT, FIELD_NAME, LTE, INTEGER, CLOSE_BRACE}},
		{` traces{ isRoot = 1 } `, []int{STREAM_TYPE_TRACES, OPEN_BRACE, FIELD_IS_ROOT, EQ, INTEGER, CLOSE_BRACE}},
		{` traces{ parent*.duration =~ events.blerg } `, []int{STREAM_TYPE_TRACES, OPEN_BRACE, FIELD_DESCENDANT, DOT, FIELD_DURATION, RE, FIELD_EVENTS, DOT, IDENTIFIER, CLOSE_BRACE}},
		{`count(spans{ parent*.duration =~ events.blerg })`, []int{AGG_COUNT, OPEN_PARENS, STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DESCENDANT, DOT, FIELD_DURATION, RE, FIELD_EVENTS, DOT, IDENTIFIER, CLOSE_BRACE, CLOSE_PARENS}},
		{`max(spans{})`, []int{AGG_MAX, OPEN_PARENS, STREAM_TYPE_SPANS, OPEN_BRACE, CLOSE_BRACE, CLOSE_PARENS}},
		{`avg(spans{})`, []int{AGG_AVG, OPEN_PARENS, STREAM_TYPE_SPANS, OPEN_BRACE, CLOSE_BRACE, CLOSE_PARENS}},
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
