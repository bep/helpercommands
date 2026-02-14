Automatically fix a GitHub issue by creating a branch, writing a test, fixing the code, and committing — all driven by Claude Code.

## Install

```
go install github.com/bep/helpercommands/cmd/fixissue@latest
```

## Usage

```
cd /path/to/your/repo
fixissue 123
```

This will:

1. Detect the GitHub repo from the current directory (requires `gh` CLI).
2. Create and switch to a `fix/issue-123` branch.
3. Run `claude -p` with a prompt that instructs Claude to write a failing test, fix the code, and commit.

## Claude Code Permissions

Since `fixissue` runs `claude -p` (headless/non-interactive mode)[^1], Claude cannot prompt you to approve tool use. You need to pre-approve the required permissions in your` .claude/settings.json`:

```json
{
  "permissions": {
    "allow": [
      "Bash(git checkout:*)",
      "Bash(git add:*)",
      "Bash(git commit:*)",
      "Bash(git diff:*)",
      "Bash(git log:*)",
      "Bash(git status:*)",
      "Bash(go test:*)",
      "Bash(make:*)",
      "Edit",
      "Read",
      "Glob",
      "Grep",
      "Write",
      "WebFetch"
    ]
  }
}
```

[^1]: We deliberately do not run `claude -p --dangerously-skip-permissions`.
