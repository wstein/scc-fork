// SPDX-License-Identifier: MIT

package processor

import "testing"

// Bosque has no `case` keyword at all (confirmed against the real parser:
// github.com/BosqueLanguage/BosqueCore src/frontend/parser.ts and
// parser_kw.ts — switch/match/dispatch arms are `guard => block`, never
// `case guard => block`). It previously counted `case ` as complexity, a
// pure fabrication that could never fire on real Bosque code. `dispatch` is
// a real third arm-dispatch construct (parseDispatchStatement) that wasn't
// tracked at all, and `elif`/`else` are real if-elif-else keywords
// (parseIfElseStatement) that were also missing.
func TestBosqueCaseIsNotComplexityButDispatchIs(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Bosque"}
	fileJob.SetContent("case = 1; dispatch (x) { T => a() }")
	CountStats(&fileJob)

	if fileJob.Complexity != 1 {
		t.Errorf("expected 1 complexity (dispatch fires, fabricated `case` does not), got %d", fileJob.Complexity)
	}
}

func TestBosqueElifElseCountComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Bosque"}
	fileJob.SetContent("if (x) { a() } elif (y) { b() } else { c() }")
	CountStats(&fileJob)

	if fileJob.Complexity != 3 {
		t.Errorf("expected 3 complexity (if, elif, else), got %d", fileJob.Complexity)
	}
}
