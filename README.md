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
- `tt add` - Stage files for commit
- `tt c` or `tt commit` - Commit changes with style
- `tt reset` - Hard-reset the repository after confirmation
- `tt branch` - Create, switch, and list git branches
- `tt merge` - Merge branches with intelligent conflict handling
- `tt push` - Push changes to remote repository
- `tt pull` - Pull changes from remote repository
- `tt clone` - Clone a repository into a new directory
- `tt log` - Show commit history with beautiful formatting
- `tt stash` - Stash changes with style
- `tt status` - Show git repository status
- `tt tag` - Create and manage git tags
- `tt revert` - Revert a commit by creating a new commit that undoes the changes
- `tt diff` - Show styled git diff with optional AI overview
- `tt aic` - Generate AI-powered commit messages
- `tt ap` - Generate AI commit message and push changes
- `tt get` - Get the current configuration values
- `tt set` - Set configuration values

### Diff Command

The `tt diff` command displays git changes with enhanced styling and optional AI-powered overview.

```bash
tt diff
```

#### With AI Overview

Generate an AI summary of your changes:

```bash
tt diff --ai
# or
tt diff -a
```

#### Show stat summary

```bash
tt diff --stat
# or
tt diff -s
```

#### Show only file names

```bash
tt diff --name-only
# or
tt diff -n
```

This will:
1. Display changes with color-coded additions (green) and deletions (red)
2. Style diff headers with bold blue
3. Optionally generate AI summary explaining what the changes achieve
4. Support all standard git diff arguments (e.g., `tt diff HEAD~1`, `tt diff main..feature`)

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

### Stash Command

The `tt stash` command provides an easy way to stash your changes, always including untracked files for simplicity.

#### Stash changes

```bash
tt stash "WIP: login fix"
```

Or without message for interactive prompt:

```bash
tt stash
```

This will:
1. Show a preview of files to be stashed
2. Prompt for a message if not provided
3. Stash all changes including untracked files (no need to remember `--include-untracked`)

#### List stashes

```bash
tt stash list
```

Shows a simplified list of your stashes with dates and messages.

#### Apply latest stash

```bash
tt stash pop
```

Applies the latest stash with confirmation and warns about potential conflicts.

### Add Command

The `tt add` command stages files for commit, mirroring `git add` with styled output.

#### Stage all changes

```bash
tt add -A
# or
tt add --all
```

#### Stage specific files or directories

```bash
tt add file.txt
tt add src/
tt add -p file1.txt -p file2.txt
```

This will:
1. Display a header and progress indicator
2. Execute `git add` with the specified files
3. Show success or error messages with styling
