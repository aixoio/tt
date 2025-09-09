# tt - Git Helper Tool

A beautifully styled, 100% git compatible tool built to make git operations intuitive and enjoyable for developers.

## Features

• **Interactive prompts** for commit messages and branch names
• **Smart defaults** for common git operations
• **Beautiful styling** with clear visual feedback
• **Beautiful log view** with styled commit history and graph support
• **Auto-push** and upstream management
• **Conflict-aware** merge operations
• **Safe reset** command with confirmation prompt

## Installation

```bash
go build -o tt .
```

## Usage

Get started by running: `tt init`

### Commands

- `tt init` - Initialize a new git repository
- `tt commit` or `tt c` - Commit changes with style
- `tt reset` - Hard-reset the repository after confirmation
- `tt branch` - Create and manage branches
- `tt merge` - Merge branches with conflict resolution
- `tt push` - Push changes to remote
- `tt log` - View styled commit history

### Reset Command

The `tt reset` command performs a hard reset of the repository, discarding all uncommitted changes. It requires user confirmation before proceeding.

```bash
tt reset
```

This will:
1. Show a confirmation prompt with the commands to be executed
2. If confirmed, run `git add .` followed by `git reset --hard`
3. Display success or error messages accordingly

**Warning:** This action cannot be undone. Make sure to backup any important uncommitted changes.