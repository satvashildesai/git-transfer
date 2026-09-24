# Git Transfer

> Move your work, not your Git history.

**Git Transfer** is a cross-platform CLI tool that allows you to transfer your current Git working state from one machine to another **without creating a commit on the source machine**.

## Why Git Transfer?

Sometimes you start working on a project on one machine and later want to continue the same work on another machine.

For example:

```text
Machine A
    ↓
Make changes
    ↓
Don't commit
    ↓
Create Git Transfer bundle
    ↓
Transfer the bundle
    ↓
Machine B
    ↓
Apply the bundle
    ↓
Continue working
```

Git Transfer is designed to preserve your working state, including:

* Modified files
* Deleted files
* New/untracked files
* Binary files
* Staged changes
* Unstaged changes
* Mixed staged + unstaged changes

## Installation

### Prerequisites

* [Go](https://go.dev/dl/) 1.22 or later
* [Git](https://git-scm.com/) installed and available in your PATH

### Build from source

Clone the repository and build the executable:

```bash
git clone https://github.com/satvashildesai/git-transfer.git
cd git-transfer
go build -o git-transfer.exe ./cmd/git-transfer   # Windows
go build -o git-transfer ./cmd/git-transfer        # macOS / Linux
```

Or install it directly into your `$GOPATH/bin`:

```bash
go install ./cmd/git-transfer
```

### Verify the installation

```bash
git-transfer --help
```

You should see:

```text
Git Transfer manages portable uncommitted working state.
It allows developers to transfer their current Git working state from one machine to another
without committing those changes on the source machine.

Usage:
  git-transfer [command]

Available Commands:
  apply       Validate the bundle against the current repository and restore its working state
  bundle      Capture the current Git working state and create a bundle
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  inspect     Preview bundle contents
  verify      Check bundle integrity and repository compatibility without modifying working tree

Flags:
  -h, --help   help for git-transfer

Use "git-transfer [command] --help" for more information about a command.
```

## Usage

### Create a bundle

On the source machine, navigate to your Git repository and run:

```bash
git-transfer bundle my_work.gtb
```

This creates a portable `.gtb` file containing your current working state.

The command does **not** create a normal Git commit or modify your Git history.

### Transfer the bundle

Copy the `.gtb` file to the other machine using any method (USB drive, email, cloud storage, SCP, etc.).

### Apply the bundle

On the destination machine, navigate to the **same Git repository** and run:

```bash
git-transfer apply my_work.gtb
```

Your changes will be restored into the working tree.

You can then check:

```bash
git status
```

and continue working normally.

## Command Reference

| Command | Description |
| :--- | :--- |
| `git-transfer bundle <file.gtb>` | Capture the current Git working state and create a bundle |
| `git-transfer apply <file.gtb>` | Validate the bundle and restore its working state |
| `git-transfer inspect <file.gtb>` | Preview bundle contents without applying |
| `git-transfer verify <file.gtb>` | Check bundle integrity and repository compatibility |

### Inspect a bundle

```bash
git-transfer inspect my_work.gtb
```

Shows information about the bundle without applying it.

### Verify a bundle

```bash
git-transfer verify my_work.gtb
```

Checks whether the bundle is valid and has not been corrupted.

## Quickstart

Here is a complete hands-on example you can try right now:

### 1. Build

```bash
go build -o git-transfer.exe ./cmd/git-transfer   # Windows
```

### 2. Make some uncommitted changes

```bash
echo "work in progress" > test_file.txt
git status
# test_file.txt shows as untracked
```

### 3. Create a bundle

```bash
git-transfer bundle my_work.gtb
# Creates my_work.gtb containing your uncommitted changes
```

### 4. Inspect the bundle

```bash
git-transfer inspect my_work.gtb
```

### 5. Simulate applying on another machine

```bash
# Remove the test file to simulate a clean workspace
rm test_file.txt        # macOS / Linux
del test_file.txt       # Windows

# Apply the bundle
git-transfer apply my_work.gtb

# Verify the file is restored
git status
# test_file.txt is back!
```

### 6. Continue working and commit normally

```bash
git add .
git commit -m "My changes"
```

## Safety

Git Transfer is designed to be safe by default.

Before applying a bundle, it can verify:

* Bundle integrity
* Bundle format version
* Git repository
* Base commit
* File conflicts
* Unsafe file paths

If the destination repository is incompatible or applying the changes could overwrite existing work, Git Transfer will stop instead of silently losing data.

## Design

Git Transfer uses Git's native capabilities where possible instead of implementing its own Git patch system.

The tool provides a simple portable layer around Git working-state transfer:

```text
Git Repository
      │
      ▼
Git Transfer
      │
      ▼
  my_work.gtb
      │
      ▼
Git Transfer
      │
      ▼
Git Repository
```

## Project Status

🚧 **Early development**

The project is currently being designed and implemented.

## License

License TBD.
