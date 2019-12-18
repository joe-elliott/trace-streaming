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
		{` spans{ name <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_NAME, LTE, NUMBER, CLOSE_BRACE}},
		{` spans{ duration <= 13 } `, []int{STREAM_TYPE_SPANS, OPEN_BRACE, FIELD_DURATION, LTE, NUMBER, CLOSE_BRACE}},
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
