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
%type <Selector>              traceSelector
%type <Matchers>              traceMatchers
%type <Matcher>               traceMatcher

%type <Field>                 traceField
%type <Field>                 spanField
%type <Field>                 processField
%type <Field>                 statusField

%type <Operator>              operator

%token <str>      IDENTIFIER STRING
%token <integer>  INTEGER
%token <float>    FLOAT
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE EQ NEQ RE NRE GT GTE LT LTE
                  STREAM_TYPE_SPANS STREAM_TYPE_TRACES
                  FIELD_DURATION FIELD_NAME FIELD_ATTS FIELD_EVENTS FIELD_STATUS FIELD_CODE FIELD_MSG FIELD_PROCESS FIELD_PARENT FIELD_DESCENDANT FIELD_SPAN FIELD_ROOT_SPAN

%%

expr:
      STREAM_TYPE_SPANS  spanSelector      { yylex.(*lexer).expr = newExpr(STREAM_TYPE_SPANS, $2) }
    | STREAM_TYPE_TRACES traceSelector     { yylex.(*lexer).expr = newExpr(STREAM_TYPE_TRACES, $2) }
    ;

spanSelector:
      OPEN_BRACE spanMatchers CLOSE_BRACE  { $$ = $2 }
    ;

spanMatchers:
      spanMatcher                          { $$ = []ValueMatcher{ $1 } }
    | spanMatchers COMMA spanMatcher       { $$ = append($1, $3)       }
    ;

spanMatcher:
      spanField operator STRING            { $$ = newStringMatcher($3, $2, $1) }
    | spanField operator INTEGER           { $$ = newIntMatcher($3, $2,  $1) }
    | spanField operator FLOAT             { $$ = newFloatMatcher($3, $2, $1) }
    ;

traceSelector:
      OPEN_BRACE traceMatchers CLOSE_BRACE { $$ = $2 }
    ;

traceMatchers:
      traceMatcher                         { $$ = []ValueMatcher{ $1 } }
    | traceMatchers COMMA traceMatcher     { $$ = append($1, $3)       }
    ;

traceMatcher:
      traceField operator STRING           { $$ = newStringMatcher($3, $2, $1) }
    | traceField operator INTEGER          { $$ = newIntMatcher($3, $2,  $1) }
    | traceField operator FLOAT            { $$ = newFloatMatcher($3, $2, $1) }
    ;

traceField:
      FIELD_SPAN DOT spanField             { $$ = wrapComplexField(FIELD_SPAN, $3)      }
    | FIELD_ROOT_SPAN DOT spanField        { $$ = wrapComplexField(FIELD_ROOT_SPAN, $3) }
    ;

spanField:
      FIELD_DURATION                      { $$ = newComplexField(FIELD_DURATION, "")    }
    | FIELD_NAME                          { $$ = newComplexField(FIELD_NAME, "")        }
    | FIELD_DESCENDANT DOT spanField      { $$ = wrapComplexField(FIELD_DESCENDANT, $3) }  
    | FIELD_PARENT DOT spanField          { $$ = wrapComplexField(FIELD_PARENT, $3)     }
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
