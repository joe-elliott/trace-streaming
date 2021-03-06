// goyacc -o processor/streamprocessor/traceql/expr.y.go processor/streamprocessor/traceql/expr.y

%{
package traceql
%}

%union{
  TempExpr  *Expr

  Selector  []matcher
  Matchers  []matcher
  Matcher   matcher
  TempField field
  LHSField  field
  RHSField  field

  AggregateArgs []float64
  AggregateArg  float64

  Operator  int
  AggregateFunc int

  str       string
  integer   int
  float     float64
}

%start expr

%type <TempExpr>              spanExpr
%type <TempExpr>              metricsExpr

%type <AggregateArgs>         aggregateArgs

%type <Selector>              selector
%type <Matchers>              matchers
%type <Matcher>               matcher

%type <LHSField>              rhsField
%type <RHSField>              lhsField
%type <TempField>             field
%type <TempField>             processField
%type <TempField>             statusField

%type <Operator>              operator
%type <AggregateFunc>         aggregateFunc

%token <str>      IDENTIFIER STRING
%token <integer>  INTEGER
%token <float>    FLOAT
%token <val>      COMMA DOT OPEN_BRACE CLOSE_BRACE OPEN_BRACKET CLOSE_BRACKET OPEN_PARENS CLOSE_PARENS
                  EQ NEQ RE NRE GT GTE LT LTE
                  STREAM_TYPE_SPANS STREAM_TYPE_TRACES
                  AGG_COUNT AGG_MAX AGG_MIN AGG_SUM AGG_AVG AGG_HIST
                  FIELD_DURATION FIELD_NAME FIELD_ATTS FIELD_EVENTS FIELD_STATUS FIELD_CODE FIELD_MSG FIELD_PROCESS FIELD_PARENT FIELD_DESCENDANT FIELD_IS_ROOT

%%

expr:
      spanExpr                         { yylex.(*lexer).expr = $1 }
    | metricsExpr                      { yylex.(*lexer).expr = $1 }
    | STREAM_TYPE_TRACES selector      { yylex.(*lexer).expr = newExpr(STREAM_TYPE_TRACES, $2) }
    ;

metricsExpr:
      AGG_COUNT OPEN_PARENS spanExpr CLOSE_PARENS                              { $$ = newMetricsExpr(AGG_COUNT, $3, nil, nil) }
    | AGG_HIST OPEN_PARENS spanExpr DOT field COMMA aggregateArgs CLOSE_PARENS { $$ = newMetricsExpr(AGG_HIST, $3, $5, $7)    }
    | aggregateFunc OPEN_PARENS spanExpr DOT field CLOSE_PARENS                { $$ = newMetricsExpr($1, $3, $5, nil)         }
    ;

spanExpr:
      STREAM_TYPE_SPANS  selector      { $$ = newExpr(STREAM_TYPE_SPANS, $2) }
    ;

aggregateArgs:
      FLOAT                          { $$ = []float64{ $1 } }
    | aggregateArgs COMMA FLOAT      { $$ = append($1, $3)  }
    ;

aggregateFunc:
      AGG_MAX                         { $$ = AGG_MAX  }
    | AGG_MIN                         { $$ = AGG_MIN  }
    | AGG_SUM                         { $$ = AGG_SUM  }
    | AGG_AVG                         { $$ = AGG_AVG  }
    ;

selector:
      OPEN_BRACE CLOSE_BRACE           { }
    | OPEN_BRACE matchers CLOSE_BRACE  { $$ = $2 }
    ;

matchers:
      matcher                          { $$ = []matcher{ $1 } }
    | matchers COMMA matcher           { $$ = append($1, $3)       }
    ;

matcher:
      lhsField operator rhsField          { $$ = newMatcher($1, $2, $3)  }
    ;

lhsField:
      field                               { $$ = $1 }
    ;

rhsField:
      field                               { $$ = $1 }
    ;

field:
      STRING                              { $$ = newStringField($1)                     }
    | INTEGER                             { $$ = newIntField($1)                        }
    | FLOAT                               { $$ = newFloatField($1)                      }
    | FIELD_DURATION                      { $$ = newDynamicField(FIELD_DURATION, "")    }
    | FIELD_IS_ROOT                       { $$ = newDynamicField(FIELD_IS_ROOT, "")     }
    | FIELD_NAME                          { $$ = newDynamicField(FIELD_NAME, "")        }
    | FIELD_DESCENDANT DOT field          { $$ = wrapRelationshipField(FIELD_DESCENDANT, $3) }  
    | FIELD_PARENT DOT field              { $$ = wrapRelationshipField(FIELD_PARENT, $3)     }
    | FIELD_PROCESS DOT processField      { $$ = wrapDynamicField(FIELD_PROCESS, $3)    }
    | FIELD_STATUS DOT statusField        { $$ = wrapDynamicField(FIELD_STATUS, $3)     }
    | FIELD_ATTS OPEN_BRACKET STRING CLOSE_BRACKET   { $$ = newDynamicField(FIELD_ATTS, $3)        }
    | FIELD_EVENTS OPEN_BRACKET STRING CLOSE_BRACKET { $$ = newDynamicField(FIELD_EVENTS, $3)      }
    ;

processField:
      FIELD_NAME                          { $$ = newDynamicField(FIELD_NAME, "")  }
    ;

statusField:
      FIELD_CODE                          { $$ = newDynamicField(FIELD_CODE, "")  }
    | FIELD_MSG                           { $$ = newDynamicField(FIELD_MSG, "")   }
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
