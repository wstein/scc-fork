// SPDX-License-Identifier: MIT

package processor

import "testing"

// TypeScript, TypeScript Typings, ArkTs, C3, C Header and Templ previously
// counted both `switch` and every `case` arm, inflating complexity for any
// multi-arm switch relative to Java/Kotlin/Go/Scala, which count the
// dispatch keyword once regardless of arm count. A structurally identical
// 5-case switch scored 1 in Java and 6 in TypeScript before this fix.
func TestSwitchCaseCountsOnceNotPerArm(t *testing.T) {
	ProcessConstants()

	src := "switch (x) { case 1: a(); case 2: b(); case 3: c(); case 4: d(); case 5: e(); default: f(); }"

	for _, lang := range []string{"TypeScript", "TypeScript Typings", "ArkTs", "C3", "C Header"} {
		fileJob := FileJob{Language: lang}
		fileJob.SetContent(src)
		CountStats(&fileJob)

		if fileJob.Complexity != 1 {
			t.Errorf("%s: expected 1 complexity (one `switch`, no per-case ticks), got %d", lang, fileJob.Complexity)
		}
	}
}

func TestTemplSwitchCaseCountsOnceNotPerArm(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "Templ"}
	fileJob.SetContent("switch x { case 1: a() case 2: b() case 3: c() }")
	CountStats(&fileJob)

	if fileJob.Complexity != 1 {
		t.Errorf("Templ: expected 1 complexity (one `switch`, no per-case ticks), got %d", fileJob.Complexity)
	}
}

// SystemVerilog has no `switch` keyword at all (confirmed against IEEE 1800
// case-statement references) and no `casexz` variant (only case/casex/casez
// are real) — both were previously listed as complexitychecks and never
// fire on real code. The real case/casex/casez keywords already count once
// per statement (SystemVerilog doesn't repeat the keyword per arm the way
// C-family switch/case does), so they're untouched.
func TestSystemVerilogFabricatedKeywordsAreNotComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "SystemVerilog"}
	fileJob.SetContent("initial begin switch = 1; casexz = 2; end")
	CountStats(&fileJob)

	if fileJob.Complexity != 0 {
		t.Errorf("expected 0 complexity from fabricated switch/casexz tokens, got %d", fileJob.Complexity)
	}
}

func TestSystemVerilogRealCaseVariantsCountComplexity(t *testing.T) {
	ProcessConstants()

	fileJob := FileJob{Language: "SystemVerilog"}
	fileJob.SetContent(`initial begin
    case (a) 1: x = 1; endcase
    casex (b) 1: y = 1; endcase
    casez (c) 1: z = 1; endcase
end`)
	CountStats(&fileJob)

	if fileJob.Complexity != 3 {
		t.Errorf("expected 3 complexity (one per case/casex/casez statement), got %d", fileJob.Complexity)
	}
}
