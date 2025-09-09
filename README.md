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
- `tt branch` - Create, switch, and list git branches
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

### Branch Command

The `tt branch` command provides comprehensive branch management with listing, creating, and switching capabilities.

#### List all branches and recent commit graph
```bash
tt branch
```

This will:
1. Display the current branch
2. List all local branches with the current branch marked
3. Show a paginated commit graph with recent commits

#### Create a new branch
```bash
tt branch <branch-name>

# Optionally auto-push the new branch to remote
tt branch <branch-name> --push
```

#### Switch to an existing branch
```bash
tt branch <branch-name>
```

If the branch exists, tt will switch to it (no action required). If the branch does not exist, tt will create a new branch with that name.

#### Delete a branch
```bash
tt branch delete <branch>

# Delete a remote branch (requires confirmation phrase)
tt branch delete --remote <branch>
```

This will:
1. Verify the branch exists and is not the current branch.
2. For local branches, try to delete safely with `--merged`.
3. If the branch has unmerged commits, prompt for confirmation to force delete.
4. For remote branches, require typing the exact phrase "confirm delete remote <branch>" to proceed.
5. Display success or error messages accordingly.