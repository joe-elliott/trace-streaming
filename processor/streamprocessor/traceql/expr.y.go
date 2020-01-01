// Code generated by goyacc -l -o processor/streamprocessor/traceql/expr.y.go processor/streamprocessor/traceql/expr.y. DO NOT EDIT.
package traceql

import __yyfmt__ "fmt"

type yySymType struct {
	yys       int
	Selector  []matcher
	Matchers  []matcher
	Matcher   matcher
	TempField field
	LHSField  field
	RHSField  field
	Operator  int

	str     string
	integer int
	float   float64
}

const IDENTIFIER = 57346
const STRING = 57347
const INTEGER = 57348
const FLOAT = 57349
const COMMA = 57350
const DOT = 57351
const OPEN_BRACE = 57352
const CLOSE_BRACE = 57353
const OPEN_BRACKET = 57354
const CLOSE_BRACKET = 57355
const EQ = 57356
const NEQ = 57357
const RE = 57358
const NRE = 57359
const GT = 57360
const GTE = 57361
const LT = 57362
const LTE = 57363
const STREAM_TYPE_SPANS = 57364
const STREAM_TYPE_TRACES = 57365
const FIELD_DURATION = 57366
const FIELD_NAME = 57367
const FIELD_ATTS = 57368
const FIELD_EVENTS = 57369
const FIELD_STATUS = 57370
const FIELD_CODE = 57371
const FIELD_MSG = 57372
const FIELD_PROCESS = 57373
const FIELD_PARENT = 57374
const FIELD_DESCENDANT = 57375
const FIELD_IS_ROOT = 57376

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENTIFIER",
	"STRING",
	"INTEGER",
	"FLOAT",
	"COMMA",
	"DOT",
	"OPEN_BRACE",
	"CLOSE_BRACE",
	"OPEN_BRACKET",
	"CLOSE_BRACKET",
	"EQ",
	"NEQ",
	"RE",
	"NRE",
	"GT",
	"GTE",
	"LT",
	"LTE",
	"STREAM_TYPE_SPANS",
	"STREAM_TYPE_TRACES",
	"FIELD_DURATION",
	"FIELD_NAME",
	"FIELD_ATTS",
	"FIELD_EVENTS",
	"FIELD_STATUS",
	"FIELD_CODE",
	"FIELD_MSG",
	"FIELD_PROCESS",
	"FIELD_PARENT",
	"FIELD_DESCENDANT",
	"FIELD_IS_ROOT",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 66

var yyAct = [...]int{

	12, 13, 14, 49, 50, 47, 7, 54, 11, 53,
	9, 12, 13, 14, 2, 3, 40, 39, 5, 15,
	17, 22, 23, 21, 38, 37, 20, 19, 18, 16,
	15, 17, 22, 23, 21, 43, 41, 20, 19, 18,
	16, 52, 36, 35, 44, 45, 27, 28, 29, 30,
	31, 32, 33, 34, 25, 51, 4, 24, 26, 48,
	6, 46, 10, 42, 8, 1,
}
var yyPact = [...]int{

	-8, -1000, 8, 8, -1000, -5, -1000, -1000, 46, -1000,
	32, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 34, 33,
	16, 15, 5, 4, -1000, 6, 6, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 6, 6, -20, -26, 50,
	36, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -4, -6, -1000, -1000,
}
var yyPgo = [...]int{

	0, 65, 56, 64, 10, 63, 62, 8, 61, 59,
	58,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 4, 6, 5,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 8, 9, 9, 10, 10, 10, 10, 10,
	10, 10, 10,
}
var yyR2 = [...]int{

	0, 2, 2, 2, 3, 1, 3, 3, 1, 1,
	1, 1, 1, 1, 1, 1, 3, 3, 3, 3,
	4, 4, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1,
}
var yyChk = [...]int{

	-1000, -1, 22, 23, -2, 10, -2, 11, -3, -4,
	-6, -7, 5, 6, 7, 24, 34, 25, 33, 32,
	31, 28, 26, 27, 11, 8, -10, 14, 15, 16,
	17, 18, 19, 20, 21, 9, 9, 9, 9, 12,
	12, -4, -5, -7, -7, -7, -8, 25, -9, 29,
	30, 5, 5, 13, 13,
}
var yyDef = [...]int{

	0, -2, 0, 0, 1, 0, 2, 3, 0, 5,
	0, 8, 10, 11, 12, 13, 14, 15, 0, 0,
	0, 0, 0, 0, 4, 0, 0, 25, 26, 27,
	28, 29, 30, 31, 32, 0, 0, 0, 0, 0,
	0, 6, 7, 9, 16, 17, 18, 22, 19, 23,
	24, 0, 0, 20, 21,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yylex.(*lexer).expr = newExpr(STREAM_TYPE_SPANS, yyDollar[2].Selector)
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yylex.(*lexer).expr = newExpr(STREAM_TYPE_TRACES, yyDollar[2].Selector)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Selector = yyDollar[2].Matchers
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Matchers = []matcher{yyDollar[1].Matcher}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Matchers = append(yyDollar[1].Matchers, yyDollar[3].Matcher)
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Matcher = newMatcher(yyDollar[1].RHSField, yyDollar[2].Operator, yyDollar[3].LHSField)
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.RHSField = yyDollar[1].TempField
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.LHSField = yyDollar[1].TempField
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newStringField(yyDollar[1].str)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newIntField(yyDollar[1].integer)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newFloatField(yyDollar[1].float)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_DURATION, "")
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_IS_ROOT, "")
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_NAME, "")
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.TempField = wrapRelationshipField(FIELD_DESCENDANT, yyDollar[3].TempField)
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.TempField = wrapRelationshipField(FIELD_PARENT, yyDollar[3].TempField)
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.TempField = wrapDynamicField(FIELD_PROCESS, yyDollar[3].TempField)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.TempField = wrapDynamicField(FIELD_STATUS, yyDollar[3].TempField)
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_ATTS, yyDollar[3].str)
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_EVENTS, yyDollar[3].str)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_NAME, "")
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_CODE, "")
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.TempField = newDynamicField(FIELD_MSG, "")
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = EQ
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = NEQ
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = RE
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = NRE
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = GT
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = GTE
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = LT
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Operator = LTE
		}
	}
	goto yystack /* stack new state and value */
}
