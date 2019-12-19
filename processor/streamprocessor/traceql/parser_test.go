package traceql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newString(s string) *string {
	return &s
}

func TestParse(t *testing.T) {
	t.Skip("TODO")

	for _, tc := range []struct {
		in  string
		exp string
		err error
	}{
		{
			in:  `spans{duration=3, name="asdf"}`,
			exp: `spans{foo="bar"}`,
		},
		{
			in:  `spans{duration=3, tags.test="blerg"}`,
			exp: `spans{foo="bar"}`,
		},
		{
			in: `spans{foo="bar"} foo`,
			err: ParseError{
				msg:  "syntax error: unexpected IDENTIFIER, expecting != or !~ or |~ or |=",
				line: 1,
				col:  13,
			},
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			ast, err := ParseExpr(tc.in)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.exp, ast)
		})
	}
}
