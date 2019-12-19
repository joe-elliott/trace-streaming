// goyacc -o processor/streamprocessor/traceql/expr.y.go processor/streamprocessor/traceql/expr.y

%{
package traceql

%}

%union{
  Selector  []ValueOperator
  Matchers  []ValueOperator
  Matcher   ValueOperator
  Field     int

  integer   int
  float     float64
}

%start expr

%type <Selector>              selector
%type <Matchers>              matchers
%type <Matcher>               matcher
%type <Field>                 field

%token <str>      IDENTIFIER STRING
%token <integer>  INTEGER
%token <float>    FLOAT
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE EQ NEQ RE NRE GT GTE LT LTE
                  STREAM_TYPE_SPANS
                  FIELD_DURATION FIELD_NAME FIELD_TAGS

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
      field EQ STRING                  { }
    | field NEQ STRING                 { }
    | field RE STRING                  { }
    | field NRE STRING                 { }
    | field EQ INTEGER                 { $$ = newIntOperator($3, EQ,  $1) }
    | field NEQ INTEGER                { $$ = newIntOperator($3, NEQ, $1) }
    | field GT INTEGER                 { $$ = newIntOperator($3, GT,  $1) }
    | field GTE INTEGER                { $$ = newIntOperator($3, GTE, $1) }
    | field LT INTEGER                 { $$ = newIntOperator($3, LT,  $1) }
    | field LTE INTEGER                { $$ = newIntOperator($3, LTE, $1) }
    | field EQ FLOAT                   { }
    | field NEQ FLOAT                  { }
    | field GT FLOAT                   { }
    | field GTE FLOAT                  { }
    | field LT FLOAT                   { }
    | field LTE FLOAT                  { }
    ;

field:
      FIELD_DURATION                   { $$ = FIELD_DURATION }
    | FIELD_NAME                       { $$ = FIELD_NAME     }
    | FIELD_TAGS DOT IDENTIFIER        { $$ = FIELD_TAGS     }
    ;

%%
