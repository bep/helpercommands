package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fixissue <issue-id>")
		os.Exit(1)
	}

	issueID := os.Args[1]
	branch := "fix/issue-" + issueID

	// Detect the GitHub repo from the current directory.
	repo, err := detectRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to detect GitHub repo: %v\n", err)
		os.Exit(1)
	}

	issueURL := fmt.Sprintf("https://github.com/%s/issues/%s", repo, issueID)

	// Create and switch to the branch (use existing if it already exists).
	if err := run("git", "checkout", "-b", branch); err != nil {
		if err := run("git", "checkout", branch); err != nil {
			fmt.Fprintf(os.Stderr, "failed to switch to branch: %v\n", err)
			os.Exit(1)
		}
	}

	// Run claude to fix the issue.
	prompt := fmt.Sprintf(`Fix the issue described in %s.

* If there is a test case provided in the issue, use that. Else create a failing test case (or adjust an existing one) that demonstrates the issue. 
* Run the test a and make sure it fails.
* Then fix the code to make the test pass.
* Commit the changes.`, issueURL)
	// This allows that Claude has all the needed permissions to read the code and commit the changes.
	// We could use --dangerously-skip-permissions but that would be less secure.
	if err := run("claude", "-p", prompt); err != nil {
		fmt.Fprintf(os.Stderr, "claude failed: %v\n", err)
		os.Exit(1)
	}
}

func shell() string {
	if s := os.Getenv("SHELL"); s != "" {
		return s
	}
	return "bash"
}

func detectRepo() (string, error) {
	cmd := exec.Command(shell(), "-ic", "gh repo view --json nameWithOwner -q .nameWithOwner")
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("gh repo view failed: %s", strings.TrimSpace(string(ee.Stderr)))
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func run(name string, args ...string) error {
	parts := make([]string, 0, len(args)+1)
	parts = append(parts, name)
	for _, a := range args {
		parts = append(parts, shellQuote(a))
	}
	cmdStr := strings.Join(parts, " ")

	cmd := exec.Command(shell(), "-ic", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
