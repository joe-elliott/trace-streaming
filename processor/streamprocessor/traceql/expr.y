// goyacc -o processor/streamprocessor/traceql/expr.y.go processor/streamprocessor/traceql/expr.y

%{
package traceql

%}

%union{
  Matchers  []string
  Matcher   string
  Field     int
}

%start root

%type <Matchers>              matchers
%type <Matcher>               matcher
%type <Field>                 field

%token <str>      IDENTIFIER STRING
%token <int>      INTEGER
%token <float>    FLOAT
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE EQ NEQ RE NRE GT GTE LT LTE
                  STREAM_TYPE_SPANS
                  FIELD_DURATION FIELD_NAME FIELD_TAGS

%%

root:
      STREAM_TYPE_SPANS selector
    ;

selector:
      OPEN_BRACE matchers CLOSE_BRACE
    ;

matchers:
      matcher                          { $$ = []string{ $1 } }
    | matchers COMMA matcher           { $$ = append($1, $3) }
    ;

matcher:
      field EQ STRING                  { }
    | field NEQ STRING                 { }
    | field RE STRING                  { }
    | field NRE STRING                 { }
    | field EQ INTEGER                 { }
    | field NEQ INTEGER                { }
    | field GT INTEGER                 { }
    | field GTE INTEGER                { }
    | field LT INTEGER                 { }
    | field LTE INTEGER                { }
    | field EQ FLOAT                   { }
    | field NEQ FLOAT                  { }
    | field GT FLOAT                   { }
    | field GTE FLOAT                  { }
    | field LT FLOAT                   { }
    | field LTE FLOAT                  { }
    ;

field:
      FIELD_DURATION                   { }
    | FIELD_NAME                       { }
    | FIELD_TAGS DOT IDENTIFIER        { }
    ;

%%
