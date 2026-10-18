// SPDX-License-Identifier: MIT

package processor

import (
	"os"
	"testing"
)

// examples/language/flix.flix is advertised as a working Flix program, so its
// counts are pinned here against the file actually compiling under the real
// Flix toolchain (verified out-of-band with flix.jar). A prior version of
// this fixture used `while` and `for (x <- l) yield ...`, neither of which
// are real Flix constructs (confirmed against the compiler's own lexer
// keyword table), and failed to compile with 17 distinct errors. `while` was
// also a fabricated entry in Flix's complexitychecks list, and Datalog rule
// heads (`:-`) were not counted at all despite being genuine branch points.
// `case ` was also dropped: it's used for both enum-variant declarations and
// match arms, and scc's trie scanner can't tell those apart, so counting it
// inflated complexity by ~55% (measured against two independent corpora)
// purely from declaring algebraic data types, not from branching logic.
func TestCountStatsFlixExample(t *testing.T) {
	ProcessConstants()

	content, err := os.ReadFile("../examples/language/flix.flix")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	fileJob := FileJob{Language: "Flix"}
	fileJob.SetContent(string(content))
	CountStats(&fileJob)

	if fileJob.Lines != 66 {
		t.Errorf("Expected 66 lines got %d", fileJob.Lines)
	}
	if fileJob.Blank != 9 {
		t.Errorf("Expected 9 blank got %d", fileJob.Blank)
	}
	if fileJob.Comment != 11 {
		t.Errorf("Expected 11 comment got %d", fileJob.Comment)
	}
	if fileJob.Code != 46 {
		t.Errorf("Expected 46 code got %d", fileJob.Code)
	}
	if fileJob.Complexity != 10 {
		t.Errorf("Expected 10 complexity got %d", fileJob.Complexity)
	}
}

// while is not a Flix keyword (confirmed: "Undefined name 'while'" from the
// real compiler) and must never contribute complexity, even though it looks
// exactly like a C-family keyword.
func TestFlixWhileIsNotComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Flix"}
	fileJob.SetContent("def loop(n: Int32): Unit = while (n > 0) { n }")
	CountStats(&fileJob)

	if fileJob.Complexity != 0 {
		t.Errorf("Expected 0 complexity (while is not a real Flix keyword), got %d", fileJob.Complexity)
	}
}

// Datalog rule heads (`:-`) are a genuine branch point (a rule only fires
// when its body holds) and must be counted, matching real Datalog-heavy
// Flix programs such as the official examples/datalog/is-cyclic-.flix and
// graph-reachability.flix, which previously scored zero complexity despite
// containing recursive transitive-closure logic.
func TestFlixDatalogRuleHeadCountsComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Flix"}
	fileJob.SetContent(`def reachable(): Vector[(Int32, Int32)] =
    let pr = #{
        Edge(1, 2).
        Path(x, y) :- Edge(x, y).
        Path(x, z) :- Path(x, y), Edge(y, z).
    };
    query pr select (x, y) from Path(x, y)
`)
	CountStats(&fileJob)

	if fileJob.Complexity != 2 {
		t.Errorf("Expected 2 complexity (one per rule head), got %d", fileJob.Complexity)
	}
}

// case is used both for enum-variant declarations (no branching, just a type
// definition) and match arms (a real branch). scc's trie scanner has no
// parsing context to tell them apart, so `case` is deliberately excluded
// from Flix's complexitychecks rather than over-counting one at the expense
// of the other. `match ` alone still ticks once per match expression,
// consistent with how Java/Kotlin/Go count `switch`/`when` once regardless
// of arm count.
func TestFlixCaseIsNotComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Flix"}
	fileJob.SetContent(`enum Shape {
    case Circle(Float64),
    case Square(Float64)
}

def area(s: Shape): Float64 = match s {
    case Shape.Circle(r) => 3.14159 * r * r
    case Shape.Square(x) => x * x
}
`)
	CountStats(&fileJob)

	if fileJob.Complexity != 1 {
		t.Errorf("Expected 1 complexity (one `match`, no per-case ticks), got %d", fileJob.Complexity)
	}
}
