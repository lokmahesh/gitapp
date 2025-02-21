package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func runCommitLint(repo string) string {
	cmd := exec.Command("npx", "--yes", "commitlint", "--from", "HEAD~1", "--to", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Commit Lint Failed: %s", out.String())
	}
	return "Commit Lint Passed"
}
