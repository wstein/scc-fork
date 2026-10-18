// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestPrologExampleRuns is an opt-in check that examples/language/prolog.p
// actually runs under a real Prolog toolchain. Prolog's complexitychecks
// list was previously a C-family template with no real Prolog keywords in
// it, and there was no example fixture at all to catch that — see
// TestFlixExampleCompiles for the same problem found in Flix's example.
//
// scc doesn't bundle a Prolog toolchain, so this is skipped unless a
// `swipl` binary is on PATH.
//
// Run manually with: go test -run TestPrologExampleRuns ./...
func TestPrologExampleRuns(t *testing.T) {
	swipl, err := exec.LookPath("swipl")
	if err != nil {
		t.Skip("no Prolog toolchain available (put `swipl` on PATH to run this check)")
	}

	src, err := os.ReadFile(filepath.Join("examples", "language", "prolog.p"))
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	workDir := t.TempDir()
	srcPath := filepath.Join(workDir, "main.p")
	if err := os.WriteFile(srcPath, src, 0644); err != nil {
		t.Fatal(err)
	}

	// --on-error=halt is required: by default swipl prints a load-time
	// exception (e.g. from `:- initialization(Goal).`) to stderr but still
	// exits 0, which would make this test pass regardless of breakage.
	cmd := exec.Command(swipl, "-q", "-g", "true", "-t", "halt", "--on-error=halt", srcPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("examples/language/prolog.p failed to load/run:\n%s", stripAnsi(string(out)))
	}
}
