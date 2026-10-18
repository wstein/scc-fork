/**
 * A small but complete ANTLR 4 grammar for a calculator language.
 *
 * It deliberately exercises the constructs scc has to get right for `.g4`:
 * line, block and doc comments; single quoted literals that contain the very
 * characters used as complexity operators ('*', '+', '?', '|'); character
 * sets, which contain those characters as members rather than operators and
 * may also contain a bare quote; action blocks holding target-language code
 * whose strings are double quoted and may contain an apostrophe; the EBNF
 * suffixes ?, * and + including their non-greedy forms; and the `+=` list
 * label operator, which is PLUS_ASSIGN in ANTLR's own grammar and not a
 * repetition operator at all.
 */
grammar AntlrCalc;

@members {
    // Target-language code, not grammar. The apostrophe in "don't" below sits
    // inside a double quoted Java string; treating it as a grammar literal
    // would swallow every following line.
    public static final String NOTE = "don't read me as a literal";
}

/* A parser rule with several alternatives. Each `|` is a branch point, in the
   same way an `else` is in a C-family language. */
prog
    : stat+ EOF
    ;

stat
    : expr NEWLINE
    | ID '=' expr NEWLINE
    | NEWLINE
    ;

// `args` collects a list, so `+=` is a label assignment and not a loop.
call
    : ID '(' (args += expr (',' args += expr)*)? ')'
    ;

expr
    : '-' expr
    | expr op = ('*' | '/') expr
    | expr op = ('+' | '-') expr
    | expr '?' expr ':' expr // a conditional, quoted so it is not an operator
    | call
    | NUMBER
    | STRING
    | ID
    | '(' expr ')'
    ;

/**
 * A lexer rule with an optional fraction and exponent. The sign is deliberately
 * left to the parser's unary `'-' expr` alternative: folding it in here as
 * `'-'? DIGIT+` would let maximal munch lex `1-2` as NUMBER NUMBER, making the
 * binary minus alternative unreachable for digit-adjacent input.
 */
NUMBER
    : DIGIT+ ('.' DIGIT+)? ([eE] [+-]? DIGIT+)?
    ;

ID
    : [a-zA-Z_] [a-zA-Z_0-9]*
    ;

// A character set holding a bare apostrophe, plus a non-greedy suffix.
STRING
    : '"' ~["\r\n]*? '"'
    ;

CHAR
    : '\'' ~['\r\n] '\''
    ;

NEWLINE
    : '\r'? '\n'
    ;

BLOCK_COMMENT
    : '/*' .*? '*/' -> skip
    ;

LINE_COMMENT
    : '//' ~[\r\n]* -> skip
    ;

WS
    : [ \t]+ -> skip
    ;

fragment DIGIT
    : [0-9]
    ;
