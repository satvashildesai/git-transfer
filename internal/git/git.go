package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// FileStatus represents the status of a file in the Git repository
type FileStatus struct {
	Path     string
	Staged   bool
	Unstaged bool
	Status   string // e.g., "modified"
}

// IsGitRepo checks if the current directory is inside a Git repository.
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// GetHeadCommit returns the current HEAD commit hash.
func GetHeadCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD commit: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetModifiedTrackedFiles parses `git status` to find modified tracked files.
func GetModifiedTrackedFiles() ([]FileStatus, error) {
	cmd := exec.Command("git", "status", "-z", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	var files []FileStatus
	if len(out) == 0 {
		return files, nil
	}

	// -z separates entries with NUL
	entries := bytes.Split(out, []byte{0})
	for _, entry := range entries {
		if len(entry) < 4 {
			continue
		}
		
		statusCode := string(entry[0:2])
		path := string(entry[3:])

		// We are focusing on modified tracked files for Milestone 1.
		staged := statusCode[0] == 'M'
		unstaged := statusCode[1] == 'M'
		
		if staged || unstaged {
			files = append(files, FileStatus{
				Path:     path,
				Staged:   staged,
				Unstaged: unstaged,
				Status:   "modified",
			})
		}
	}

	return files, nil
}

// GetStagedContent reads the staged content of a file from the index.
func GetStagedContent(path string) ([]byte, error) {
	cmd := exec.Command("git", "show", ":"+path)
	return cmd.Output()
}

// AddFile stages a file using git add.
func AddFile(path string) error {
	cmd := exec.Command("git", "add", path)
	return cmd.Run()
}
