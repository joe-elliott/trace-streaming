package traceql

import (
	"text/scanner"
)

var tokens = map[string]int{
	",":        COMMA,
	".":        DOT,
	"{":        OPEN_BRACE,
	"}":        CLOSE_BRACE,
	"=":        EQ,
	"!=":       NEQ,
	"=~":       RE,
	"!~":       NRE,
	">":        GT,
	">=":       GTE,
	"<":        LT,
	"<=":       LTE,
	"spans":    STREAM_TYPE_SPANS,
	"duration": FIELD_DURATION,
	"name":     FIELD_NAME,
	"tags":     FIELD_TAGS,
}

type lexer struct {
	scanner.Scanner
	expr   string
	errs   []ParseError
	parser *yyParserImpl
}

func (l *lexer) Lex(lval *yySymType) int {
	r := l.Scan()
	switch r {
	case scanner.EOF:
		return 0

	case scanner.String:
		return STRING

	case scanner.Int:
		return INTEGER

	case scanner.Float:
		return FLOAT
	}

	if tok, ok := tokens[l.TokenText()+string(l.Peek())]; ok {
		l.Next()
		return tok
	}

	if tok, ok := tokens[l.TokenText()]; ok {
		return tok
	}

	return IDENTIFIER
}

func (l *lexer) Error(msg string) {
	l.errs = append(l.errs, newParseError(msg, l.Line, l.Column))
}
