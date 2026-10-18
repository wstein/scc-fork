// SPDX-License-Identifier: MIT

package processor

import (
	"os"
	"testing"
)

// examples/language/AntlrCalc.g4 is a real ANTLR 4 grammar, verified with the
// ANTLR 4.13.2 tool: it generates without error or warning, the generated Java
// compiles, and the resulting parser parses sample input — including `1-2` —
// into the expected tree. Its counts are pinned here.
//
// The complexity of 27 was hand-checked token by token against a per-line
// attribution dump. Every tick is a genuine branch point: there are no known
// false positives left in this fixture, because ANTLR's three opaque regions
// (single quoted literals, `[...]` character sets and double quoted strings
// inside action blocks) are all modelled as strings, so the operator characters
// they contain are never scanned.
func TestCountStatsAntlrExample(t *testing.T) {
	ProcessConstants()

	content, err := os.ReadFile("../examples/language/AntlrCalc.g4")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(string(content))
	CountStats(&fileJob)

	if fileJob.Lines != 93 {
		t.Errorf("Expected 93 lines got %d", fileJob.Lines)
	}
	if fileJob.Blank != 14 {
		t.Errorf("Expected 14 blank got %d", fileJob.Blank)
	}
	if fileJob.Comment != 26 {
		t.Errorf("Expected 26 comment got %d", fileJob.Comment)
	}
	if fileJob.Code != 53 {
		t.Errorf("Expected 53 code got %d", fileJob.Code)
	}
	if fileJob.Complexity != 27 {
		t.Errorf("Expected 27 complexity got %d", fileJob.Complexity)
	}
}

// The `|` alternative separator is the only branch construct a grammar has, so
// it plays the role `else` plays in a C-family language.
//
// Known limitation: TComplexity in workers.go requires the byte before a prefix
// check to not be an identifier character, so a pipe written tight against an
// identifier (`expr|ID`) is skipped and only the spaced form counts. Making `|`
// a postfix check instead would trade this for a worse failure — it would miss
// the leading-pipe continuation style that every alternative below uses, which
// is how essentially all real grammars are formatted (0 of 842 pipes across
// nine antlr/grammars-v4 files were written tight against an identifier).
func TestAntlrAlternativesCountComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`value
    : STRING
    | NUMBER
    | obj
    ;`)
	CountStats(&fileJob)

	if fileJob.Complexity != 2 {
		t.Errorf("Expected 2 complexity (two alternatives past the first), got %d", fileJob.Complexity)
	}
}

// ANTLR literals are single quoted and routinely contain the exact characters
// used as complexity operators. Without string handling a grammar for any
// C-like language would have its complexity inflated by every operator token it
// defines — 67 such literals across the nine surveyed grammars.
func TestAntlrQuotedOperatorsAreNotComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`STAR : '*'; PLUS : '+'; QUESTION : '?'; OR : '|'; ANDAND : '&&';`)
	CountStats(&fileJob)

	if fileJob.Complexity != 0 {
		t.Errorf("Expected 0 complexity (all operators are quoted literals), got %d", fileJob.Complexity)
	}
}

// A grammar that lexes comments contains `'/*'` and `'//'` as literals, and
// grammars that hand the delimiters to different rules leave them unpaired.
// An unpaired `'/*'` literal that was allowed to open a real comment would
// swallow every following line until some later rule happened to mention
// `'*/'` — here that is two whole rules.
func TestAntlrQuotedCommentDelimitersStayCode(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`OPEN : '/*';
LINE : '//' ~[\r\n]*;
ID : [a-z]+;
CLOSE : '*/';`)
	CountStats(&fileJob)

	if fileJob.Code != 4 {
		t.Errorf("Expected 4 code lines (quoted delimiters must not open a comment), got %d", fileJob.Code)
	}
	if fileJob.Comment != 0 {
		t.Errorf("Expected 0 comment lines, got %d", fileJob.Comment)
	}
}

// A `[...]` character set is an opaque leaf token, not a branch, and it is the
// one place a bare apostrophe legitimately appears in grammar text. Modelling
// sets as strings is what stops `[']` — which every grammar that lexes quoted
// strings writes — from opening a literal that runs to the next apostrophe,
// swallowing the intervening lines. This regressed exactly that way before
// character sets were modelled: both comment lines below were counted as code
// and the `+` branch point was lost.
func TestAntlrCharacterSetIsOpaque(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`APOS : ['] ;
// this is a real comment line
/* and a real block comment */
ID : [a-z]+ ;`)
	CountStats(&fileJob)

	if fileJob.Code != 2 {
		t.Errorf("Expected 2 code lines (a quote inside a character set must not open a literal), got %d", fileJob.Code)
	}
	if fileJob.Comment != 2 {
		t.Errorf("Expected 2 comment lines, got %d", fileJob.Comment)
	}
	if fileJob.Complexity != 1 {
		t.Errorf("Expected 1 complexity (only the `+` on [a-z]+), got %d", fileJob.Complexity)
	}
}

// Operator characters inside a character set are set members, not EBNF
// suffixes: `[+-]` matches a plus or a minus, it does not repeat anything.
func TestAntlrCharacterSetMembersAreNotComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`EXP : [eE] [+-] [0-9] ;`)
	CountStats(&fileJob)

	if fileJob.Complexity != 0 {
		t.Errorf("Expected 0 complexity (`+` and `-` are set members), got %d", fileJob.Complexity)
	}
}

// Action blocks hold target-language code, whose strings are double quoted and
// routinely contain apostrophes. With only single quotes modelled, the `'` in
// `don't` opened a grammar literal that never closed, so every following line
// was swallowed — comments counted as code and all later branch points lost.
// `@header`/`@members`/inline `{...}` actions are common in real grammars.
func TestAntlrActionBlockStringDoesNotSwallowFile(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`grammar A;
@members {
public void warn() { System.out.println("can't parse"); }
}
// a real comment line
prog : stat+ EOF ;
stat : expr | ID ;`)
	CountStats(&fileJob)

	if fileJob.Comment != 1 {
		t.Errorf("Expected 1 comment line (the apostrophe must stay inside the Java string), got %d", fileJob.Comment)
	}
	if fileJob.Code != 6 {
		t.Errorf("Expected 6 code lines, got %d", fileJob.Code)
	}
	if fileJob.Complexity != 2 {
		t.Errorf("Expected 2 complexity (the `+` suffix and the `|` alternative), got %d", fileJob.Complexity)
	}
}

// `+=` is PLUS_ASSIGN in ANTLR's own grammar (`labeledElement : identifier
// (ASSIGN | PLUS_ASSIGN) (atom | block)`), a list-label assignment rather than
// the PLUS repetition operator, so it is not a branch point.
func TestAntlrListLabelIsNotComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`call : ID '(' args += expr (',' args += expr) ')' ;`)
	CountStats(&fileJob)

	if fileJob.Complexity != 0 {
		t.Errorf("Expected 0 complexity (`+=` is a label assignment, not a loop), got %d", fileJob.Complexity)
	}
}

// ANTLR's ebnfSuffix is `QUESTION QUESTION? | STAR QUESTION? | PLUS QUESTION?`,
// so a non-greedy suffix such as `.*?` is one operator spelled with two
// characters and must score one, not two.
func TestAntlrNonGreedySuffixCountsOnce(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`BLOCK : '<!--' .*? '-->' -> skip;`)
	CountStats(&fileJob)

	if fileJob.Complexity != 1 {
		t.Errorf("Expected 1 complexity (`.*?` is a single EBNF suffix), got %d", fileJob.Complexity)
	}
}

// The greedy suffixes each score one on their own, which is what makes the
// non-greedy case above consistent rather than a special case.
func TestAntlrEbnfSuffixesCountComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "ANTLR"}
	fileJob.SetContent(`prog : stat+ decl* name? EOF ;`)
	CountStats(&fileJob)

	if fileJob.Complexity != 3 {
		t.Errorf("Expected 3 complexity (one per EBNF suffix), got %d", fileJob.Complexity)
	}
}
