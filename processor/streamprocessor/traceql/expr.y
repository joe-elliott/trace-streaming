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

%type <Selector>              selector
%type <Matchers>              matchers
%type <Matcher>               matcher

%type <Field>                 field
%type <Field>                 processField
%type <Field>                 statusField

%type <Operator>              operator

%token <str>      IDENTIFIER STRING
%token <integer>  INTEGER
%token <float>    FLOAT
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE EQ NEQ RE NRE GT GTE LT LTE
                  STREAM_TYPE_SPANS STREAM_TYPE_TRACES
                  FIELD_DURATION FIELD_NAME FIELD_ATTS FIELD_EVENTS FIELD_STATUS FIELD_CODE FIELD_MSG FIELD_PROCESS FIELD_PARENT FIELD_DESCENDANT FIELD_IS_ROOT

%%

expr:
      STREAM_TYPE_SPANS  selector      { yylex.(*lexer).expr = newExpr(STREAM_TYPE_SPANS, $2) }
    | STREAM_TYPE_TRACES selector      { yylex.(*lexer).expr = newExpr(STREAM_TYPE_TRACES, $2) }
    ;

selector:
      OPEN_BRACE CLOSE_BRACE           { }
    | OPEN_BRACE matchers CLOSE_BRACE  { $$ = $2 }
    ;

matchers:
      matcher                          { $$ = []ValueMatcher{ $1 } }
    | matchers COMMA matcher           { $$ = append($1, $3)       }
    ;

matcher:
      field operator STRING            { $$ = newStringMatcher($3, $2, $1) }
    | field operator INTEGER           { $$ = newIntMatcher($3, $2,  $1) }
    | field operator FLOAT             { $$ = newFloatMatcher($3, $2, $1) }
    ;

field:
      FIELD_DURATION                      { $$ = newComplexField(FIELD_DURATION, "")    }
    | FIELD_IS_ROOT                       { $$ = newComplexField(FIELD_IS_ROOT, "")     }
    | FIELD_NAME                          { $$ = newComplexField(FIELD_NAME, "")        }
    | FIELD_DESCENDANT DOT field          { $$ = wrapComplexField(FIELD_DESCENDANT, $3) }  
    | FIELD_PARENT DOT field              { $$ = wrapComplexField(FIELD_PARENT, $3)     }
    | FIELD_PROCESS DOT processField      { $$ = wrapComplexField(FIELD_PROCESS, $3)    }
    | FIELD_STATUS DOT statusField        { $$ = wrapComplexField(FIELD_STATUS, $3)     }
    | FIELD_ATTS DOT IDENTIFIER           { $$ = newComplexField(FIELD_ATTS, $3)        }
    | FIELD_EVENTS DOT IDENTIFIER         { $$ = newComplexField(FIELD_EVENTS, $3)      }
    ;

processField:
      FIELD_NAME                          { $$ = newComplexField(FIELD_NAME, "")       }
    ;

statusField:
      FIELD_CODE                          { $$ = newComplexField(FIELD_CODE, "")  }
    | FIELD_MSG                           { $$ = newComplexField(FIELD_MSG, "")   }
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
