package traceql

import (
	"strconv"
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
	"traces":   STREAM_TYPE_TRACES,
	"duration": FIELD_DURATION,
	"name":     FIELD_NAME,
	"atts":     FIELD_ATTS,
	"events":   FIELD_EVENTS,
	"status":   FIELD_STATUS,
	"code":     FIELD_CODE,
	"message":  FIELD_MSG,
	"process":  FIELD_PROCESS,
	"parent":   FIELD_PARENT,
	"parent*":  FIELD_DESCENDANT,
	"span":     FIELD_SPAN,
	"rootSpan": FIELD_ROOT_SPAN,
}

type lexer struct {
	scanner.Scanner
	expr   *Expr
	errs   []ParseError
	parser *yyParserImpl
}

func (l *lexer) Lex(lval *yySymType) int {
	r := l.Scan()
	switch r {
	case scanner.EOF:
		return 0

	case scanner.String:
		var err error
		lval.str, err = strconv.Unquote(l.TokenText())
		if err != nil {
			l.Error(err.Error())
			return 0
		}
		return STRING

	case scanner.Int:
		var err error
		lval.integer, err = strconv.Atoi(l.TokenText())
		if err != nil {
			l.Error(err.Error())
			return 0
		}
		return INTEGER

	case scanner.Float:
		var err error
		lval.float, err = strconv.ParseFloat(l.TokenText(), 64)
		if err != nil {
			l.Error(err.Error())
			return 0
		}
		return FLOAT
	}

	if tok, ok := tokens[l.TokenText()+string(l.Peek())]; ok {
		l.Next()
		return tok
	}

	if tok, ok := tokens[l.TokenText()]; ok {
		return tok
	}

	lval.str = l.TokenText()
	return IDENTIFIER
}

func (l *lexer) Error(msg string) {
	l.errs = append(l.errs, newParseError(msg, l.Line, l.Column))
}
