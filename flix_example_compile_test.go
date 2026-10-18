// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestFlixExampleCompiles is an opt-in check that examples/language/flix.flix
// actually compiles under a real Flix toolchain. A prior version of that file
// looked like plausible Flix but failed with 17 compiler errors (invalid
// `for`/`while`/`ref` syntax, an unqualified exception type, and a Datalog
// type mismatch) — nothing in this repo's normal test suite could have caught
// that, since scc only tokenizes, it never parses or type-checks.
//
// scc doesn't bundle a Flix toolchain (that's a ~30MB JVM-based compiler, not
// something every contributor or CI run should need), so this test is skipped
// unless one is available:
//   - SCC_FLIX_JAR pointing at a flix.jar (invoked as `java -jar $SCC_FLIX_JAR check`), or
//   - a `flix` binary on PATH (invoked as `flix check`).
//
// Run manually with: SCC_FLIX_JAR=/path/to/flix.jar go test -run TestFlixExampleCompiles ./...
func TestFlixExampleCompiles(t *testing.T) {
	var cmdArgs []string
	if jar := os.Getenv("SCC_FLIX_JAR"); jar != "" {
		if _, err := os.Stat(jar); err != nil {
			t.Skipf("SCC_FLIX_JAR=%s does not exist: %v", jar, err)
		}
		javaPath, err := exec.LookPath("java")
		if err != nil {
			t.Skip("SCC_FLIX_JAR set but no java on PATH")
		}
		cmdArgs = []string{javaPath, "-jar", jar, "check"}
	} else if flixPath, err := exec.LookPath("flix"); err == nil {
		cmdArgs = []string{flixPath, "check"}
	} else {
		t.Skip("no Flix toolchain available (set SCC_FLIX_JAR or put `flix` on PATH to run this check)")
	}

	src, err := os.ReadFile(filepath.Join("examples", "language", "flix.flix"))
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	workDir := t.TempDir()
	srcDir := filepath.Join(workDir, "src")
	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "Main.flix"), src, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("examples/language/flix.flix failed to compile:\n%s", stripAnsi(string(out)))
	}
}

func stripAnsi(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if r == '\x1b' {
			inEscape = true
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
