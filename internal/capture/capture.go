package capture

import (
	"fmt"
	"time"

	"github.com/git-transfer/git-transfer/internal/bundle"
	"github.com/git-transfer/git-transfer/internal/git"
	"github.com/git-transfer/git-transfer/internal/manifest"
)

// Capture bundles the current git state into a GTB file
func Capture(bundlePath string) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("current directory is not inside a Git repository")
	}

	head, err := git.GetHeadCommit()
	if err != nil {
		return fmt.Errorf("failed to get base commit: %w", err)
	}

	files, err := git.GetModifiedTrackedFiles()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	fmt.Printf("Scanning Git repository...\n")
	fmt.Printf("Found %d modified files.\n", len(files))

	if len(files) == 0 {
		fmt.Println("No changes to bundle.")
		return nil
	}

	w, err := bundle.NewWriter(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to create bundle: %w", err)
	}
	defer w.Close()

	m := &manifest.Manifest{
		FormatVersion: manifest.Version,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		Git: manifest.GitInfo{
			Head: head,
		},
		Files: make([]manifest.FileInfo, 0, len(files)),
	}

	for _, file := range files {
		fmt.Printf("Bundling %s...\n", file.Path)
		m.Files = append(m.Files, manifest.FileInfo{
			Path:     file.Path,
			Status:   file.Status,
			Staged:   file.Staged,
			Unstaged: file.Unstaged,
		})

		if file.Staged {
			content, err := git.GetStagedContent(file.Path)
			if err != nil {
				return fmt.Errorf("failed to read staged content for %s: %w", file.Path, err)
			}
			if err := w.AddStagedFile(file.Path, content); err != nil {
				return fmt.Errorf("failed to write staged file %s to bundle: %w", file.Path, err)
			}
		}

		if file.Unstaged {
			if err := w.AddUnstagedFile(file.Path); err != nil {
				return fmt.Errorf("failed to write unstaged file %s to bundle: %w", file.Path, err)
			}
		}
	}

	if err := w.WriteManifest(m); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	fmt.Println("Done.")
	return nil
}
