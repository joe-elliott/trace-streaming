package traceql

import (
	"log"
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
}

type lexer struct {
	scanner.Scanner
	expr string
}

func (l *lexer) Lex(lval *yySymType) int {
	r := l.Scan()
	switch r {
	case scanner.EOF:
		return 0

	case scanner.String:
		return STRING

	case scanner.Int:
		return NUMBER

	case scanner.Float:
		return NUMBER
	}

	if tok, ok := tokens[l.TokenText()+string(l.Peek())]; ok {
		l.Next()
		return tok
	}

	if tok, ok := tokens[l.TokenText()]; ok {
		return tok
	}

	return 0
}

func (l *lexer) Error(msg string) {
	log.Fatalf("oops %v", msg)
}
