// goyacc -o processor/streamprocessor/traceql/expr.y.go processor/streamprocessor/traceql/expr.y

%{
package traceql

%}

%union{
  Matchers  []string
  Matcher   string
}

%start root

%type <Matchers>              matchers
%type <Matcher>               matcher
%type <Field>                 field

%token <str>      STRING
%token <num>      NUMBER
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE EQ NEQ RE NRE GT GTE LT LTE OPEN_PARENTHESIS CLOSE_PARENTHESIS
                  STREAM_TYPE_SPANS
                  FIELD_DURATION FIELD_NAME

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
    | field EQ NUMBER                  { }
    | field NEQ NUMBER                 { }
    | field GT NUMBER                  { }
    | field GTE NUMBER                 { }
    | field LT NUMBER                  { }
    | field LTE NUMBER                 { }
    ;

field:
      FIELD_DURATION                   { }
    | FIELD_NAME                       { }
    ;
%%
