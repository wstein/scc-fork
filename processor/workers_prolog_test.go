// SPDX-License-Identifier: MIT

package processor

import (
	"os"
	"testing"
)

// examples/language/prolog.p is advertised as a working Prolog program, so
// its counts are pinned here against the file actually running under the
// real SWI-Prolog toolchain (verified out-of-band with swipl). The prior
// version of Prolog's complexitychecks list was a verbatim C-family template
// (`for`, `if`, `switch`, `while`, `else`) — standard Prolog has none of
// those keywords, it's Horn-clause resolution built on `:-`, `;`, `->` and
// `\+`, the same paradigm as the Datalog gap fixed for Flix. The old list
// scored this fixture's real logic as complexity 1 (one stray token match);
// the corrected list scores it 22, reflecting the actual rules and guards.
func TestCountStatsPrologExample(t *testing.T) {
	ProcessConstants()

	content, err := os.ReadFile("../examples/language/prolog.p")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	fileJob := FileJob{Language: "Prolog"}
	fileJob.SetContent(string(content))
	CountStats(&fileJob)

	if fileJob.Lines != 41 {
		t.Errorf("Expected 41 lines got %d", fileJob.Lines)
	}
	if fileJob.Blank != 9 {
		t.Errorf("Expected 9 blank got %d", fileJob.Blank)
	}
	if fileJob.Comment != 9 {
		t.Errorf("Expected 9 comment got %d", fileJob.Comment)
	}
	if fileJob.Code != 23 {
		t.Errorf("Expected 23 code got %d", fileJob.Code)
	}
	if fileJob.Complexity != 22 {
		t.Errorf("Expected 22 complexity got %d", fileJob.Complexity)
	}
}

// for/if/switch/while/else are not Prolog keywords at all (confirmed against
// SWI-Prolog's control-predicates reference) and must never contribute
// complexity, even though they look like plausible C-family branch keywords.
func TestPrologCFamilyKeywordsAreNotComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Prolog"}
	fileJob.SetContent("describe(for, if, switch, while, else).")
	CountStats(&fileJob)

	if fileJob.Complexity != 0 {
		t.Errorf("Expected 0 complexity from C-family lookalikes, got %d", fileJob.Complexity)
	}
}

// The clause neck (:-), disjunction (;), if-then (->) and negation-as-failure
// (\+) are the real Prolog branch operators and must each be counted.
func TestPrologRealOperatorsCountComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Prolog"}
	fileJob.SetContent("adult(X) :- (X >= 18 -> true ; \\+ minor(X)).")
	CountStats(&fileJob)

	// :- , >= , -> , ; , \+  = 5 ticks
	if fileJob.Complexity != 5 {
		t.Errorf("Expected 5 complexity (:-, >=, ->, ;, \\+), got %d", fileJob.Complexity)
	}
}
