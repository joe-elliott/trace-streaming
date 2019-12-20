// goyacc -o processor/streamprocessor/traceql/expr.y.go processor/streamprocessor/traceql/expr.y

%{
package traceql
%}

%union{
  Selector  []ValueMatcher
  Matchers  []ValueMatcher
  Matcher   ValueMatcher
  Field     complexField
  Operator  int

  str       string
  integer   int
  float     float64
}

%start expr

%type <Selector>              spanSelector
%type <Matchers>              spanMatchers
%type <Matcher>               spanMatcher
%type <Field>                 spanField

%type <Operator>              operator

%token <str>      IDENTIFIER STRING
%token <integer>  INTEGER
%token <float>    FLOAT
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE EQ NEQ RE NRE GT GTE LT LTE
                  STREAM_TYPE_SPANS
                  FIELD_DURATION FIELD_NAME FIELD_ATTS FIELD_EVENTS FIELD_STATUS FIELD_STATUS_CODE FIELD_STATUS_MSG

%%

expr:
      STREAM_TYPE_SPANS spanSelector       { yylex.(*lexer).expr = newExpr(STREAM_TYPE_SPANS, $2) }
    ;

spanSelector:
      OPEN_BRACE spanMatchers CLOSE_BRACE  { $$ = $2 }
    ;

spanMatchers:
      spanMatcher                          { $$ = []ValueMatcher{ $1 } }
    | spanMatchers COMMA spanMatcher           { $$ = append($1, $3) }
    ;

spanMatcher:
      spanField operator STRING            { $$ = newStringMatcher($3, $2, $1) }
    | spanField operator INTEGER           { $$ = newIntMatcher($3, $2,  $1) }
    | spanField operator FLOAT             { $$ = newFloatMatcher($3, $2, $1) }
    ;

spanField:
      FIELD_DURATION                      { $$ = newComplexField(FIELD_DURATION, "")    }
    | FIELD_NAME                          { $$ = newComplexField(FIELD_NAME, "")        }
    | FIELD_STATUS DOT FIELD_STATUS_CODE  { $$ = newComplexField(FIELD_STATUS_CODE, "") }
    | FIELD_STATUS DOT FIELD_STATUS_MSG   { $$ = newComplexField(FIELD_STATUS_MSG, "")  }
    | FIELD_ATTS DOT IDENTIFIER           { $$ = newComplexField(FIELD_ATTS, $3)        }
    | FIELD_EVENTS DOT IDENTIFIER         { $$ = newComplexField(FIELD_EVENTS, $3)      }
    ;

operator:
      EQ                               { $$ =  EQ }
    | NEQ                              { $$ = NEQ }
    | RE                               { $$ =  RE }
    | NRE                              { $$ = NRE }
    | GT                               { $$ =  GT }
    | GTE                              { $$ = GTE }
    | LT                               { $$ =  LT }
    | LTE                              { $$ = LTE }
    ;

%%
