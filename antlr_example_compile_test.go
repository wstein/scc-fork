// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestAntlrExampleCompiles is an opt-in check that examples/language/AntlrCalc.g4
// is a grammar the real ANTLR tool accepts, following the same pattern as the
// Flix and Prolog example checks. scc only tokenizes — it never parses — so
// nothing in the normal test suite can tell a valid grammar from something that
// merely looks like one, and the fixture's pinned counts are only meaningful if
// the fixture is genuinely representative ANTLR.
//
// scc doesn't bundle the ANTLR tool (a ~2MB jar needing a JVM), so this test is
// skipped unless one is available:
//   - SCC_ANTLR_JAR pointing at an antlr-complete jar (invoked as `java -jar $SCC_ANTLR_JAR`), or
//   - an `antlr4` or `antlr` binary on PATH. Both names are probed: the pip
//     package installs `antlr4` while Homebrew's formula installs `antlr`, and
//     probing only one silently skips this check on machines that do in fact
//     have a toolchain — which is how the `1-2` defect in an earlier version of
//     the fixture survived.
//
// Run manually with: SCC_ANTLR_JAR=/path/to/antlr-4.13.2-complete.jar go test -run TestAntlrExampleCompiles ./...
func TestAntlrExampleCompiles(t *testing.T) {
	var cmdArgs []string
	if jar := os.Getenv("SCC_ANTLR_JAR"); jar != "" {
		if _, err := os.Stat(jar); err != nil {
			t.Skipf("SCC_ANTLR_JAR=%s does not exist: %v", jar, err)
		}
		javaPath, err := exec.LookPath("java")
		if err != nil {
			t.Skip("SCC_ANTLR_JAR set but no java on PATH")
		}
		cmdArgs = []string{javaPath, "-jar", jar}
	} else if antlrPath := lookPathAny("antlr4", "antlr"); antlrPath != "" {
		cmdArgs = []string{antlrPath}
	} else {
		t.Skip("no ANTLR toolchain available (set SCC_ANTLR_JAR or put `antlr4`/`antlr` on PATH to run this check)")
	}

	// ANTLR requires the file name to match the grammar name, so the fixture
	// has to keep its own base name in the work directory.
	const fixture = "AntlrCalc.g4"
	src, err := os.ReadFile(filepath.Join("examples", "language", fixture))
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, fixture), src, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(cmdArgs[0], append(cmdArgs[1:], "-o", "gen", fixture)...)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("examples/language/%s was rejected by the ANTLR tool:\n%s", fixture, out)
	}
	// The tool exits 0 on warnings, but a curated example should not produce
	// any (an ambiguous or unreachable rule means the fixture is not the clean
	// grammar it claims to be).
	if strings.Contains(string(out), "warning(") {
		t.Errorf("examples/language/%s generated ANTLR warnings:\n%s", fixture, out)
	}
}

// lookPathAny returns the resolved path of the first name found on PATH, or ""
// when none of them are present.
func lookPathAny(names ...string) string {
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}

	return ""
}
