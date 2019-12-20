// goyacc -o processor/streamprocessor/traceql/expr.y.go processor/streamprocessor/traceql/expr.y

%{
package traceql
%}

%union{
  Selector  []ValueOperator
  Matchers  []ValueOperator
  Matcher   ValueOperator
  Field     complexField
  Operator  int

  str       string
  integer   int
  float     float64
}

%start expr

%type <Selector>              selector
%type <Matchers>              matchers
%type <Matcher>               matcher
%type <Field>                 field
%type <Operator>              operator

%token <str>      IDENTIFIER STRING
%token <integer>  INTEGER
%token <float>    FLOAT
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE EQ NEQ RE NRE GT GTE LT LTE
                  STREAM_TYPE_SPANS
                  FIELD_DURATION FIELD_NAME FIELD_ATTS FIELD_EVENTS FIELD_STATUS FIELD_STATUS_CODE FIELD_STATUS_MSG

%%

expr:
      STREAM_TYPE_SPANS selector       { yylex.(*lexer).expr = newExpr(STREAM_TYPE_SPANS, $2) }
    ;

selector:
      OPEN_BRACE matchers CLOSE_BRACE  { $$ = $2 }
    ;

matchers:
      matcher                          { $$ = []ValueOperator{ $1 } }
    | matchers COMMA matcher           { $$ = append($1, $3) }
    ;

matcher:
      field operator STRING            { $$ = newStringOperator($3, $2, $1) }
    | field operator INTEGER           { $$ = newIntOperator($3, $2,  $1) }
    | field operator FLOAT             { $$ = newFloatOperator($3, $2, $1) }
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

field:
      FIELD_DURATION                      { $$ = newComplexField(FIELD_DURATION, "")    }
    | FIELD_NAME                          { $$ = newComplexField(FIELD_NAME, "")        }
    | FIELD_STATUS DOT FIELD_STATUS_CODE  { $$ = newComplexField(FIELD_STATUS_CODE, "") }
    | FIELD_STATUS DOT FIELD_STATUS_MSG   { $$ = newComplexField(FIELD_STATUS_MSG, "")  }
    | FIELD_ATTS DOT IDENTIFIER           { $$ = newComplexField(FIELD_ATTS, $3)        }
    | FIELD_EVENTS DOT IDENTIFIER         { $$ = newComplexField(FIELD_EVENTS, $3)      }
    ;

%%
