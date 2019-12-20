package traceql

import (
	"fmt"
	"strings"
	"text/scanner"
)

func init() {
	// Improve the error messages coming out of yacc.
	yyErrorVerbose = true
}

// ParseExpr parses a string and returns an Expr.
func ParseExpr(input string) (expr *Expr, err error) {
	defer func() {
		r := recover()
		if r != nil {
			var ok bool
			if err, ok = r.(error); ok {
				return
			}
		}
	}()
	l := lexer{
		parser: yyNewParser().(*yyParserImpl),
	}
	l.Init(strings.NewReader(input))
	l.Scanner.Error = func(_ *scanner.Scanner, msg string) {
		l.Error(msg)
	}
	e := l.parser.Parse(&l)
	if e != 0 || len(l.errs) > 0 {
		return nil, l.errs[0]
	}
	return l.expr, nil
}

// ParseError is what is returned when we failed to parse.
type ParseError struct {
	msg       string
	line, col int
}

func (p ParseError) Error() string {
	if p.col == 0 && p.line == 0 {
		return fmt.Sprintf("parse error : %s", p.msg)
	}
	return fmt.Sprintf("parse error at line %d, col %d: %s", p.line, p.col, p.msg)
}

func newParseError(msg string, line, col int) ParseError {
	return ParseError{
		msg:  msg,
		line: line,
		col:  col,
	}
}
