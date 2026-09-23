package restore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-transfer/git-transfer/internal/bundle"
	"github.com/git-transfer/git-transfer/internal/git"
)

// Apply restores the bundle to the current repository
func Apply(bundlePath string) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("current directory is not inside a Git repository")
	}

	head, err := git.GetHeadCommit()
	if err != nil {
		return fmt.Errorf("failed to get base commit: %w", err)
	}

	r, err := bundle.NewReader(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to read bundle: %w", err)
	}
	defer r.Close()

	manifest, err := r.ReadManifest()
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	if head != manifest.Git.Head {
		fmt.Println("WARNING: The repository HEAD does not match the bundle base commit.")
		fmt.Printf("Bundle created from: %s\n", manifest.Git.Head)
		fmt.Printf("Current repository: %s\n", head)
		fmt.Println("Aborting.")
		return fmt.Errorf("repository mismatch")
	}

	localFiles, err := git.GetModifiedTrackedFiles()
	if err != nil {
		return fmt.Errorf("failed to get local git status: %w", err)
	}

	localConflicts := make(map[string]bool)
	for _, f := range localFiles {
		localConflicts[f.Path] = true
	}

	var conflicts []string
	for _, f := range manifest.Files {
		if localConflicts[f.Path] {
			conflicts = append(conflicts, f.Path)
		}
	}

	if len(conflicts) > 0 {
		fmt.Println("CONFLICT")
		fmt.Println("The following files already contain local changes:")
		for _, c := range conflicts {
			fmt.Printf("  %s\n", c)
		}
		fmt.Println("Nothing has been changed.")
		return fmt.Errorf("conflicts detected")
	}

	fmt.Println("Applying bundle...")
	for _, f := range manifest.Files {
		if f.Staged {
			content, err := r.GetFileContent(filepath.ToSlash(filepath.Join("files", "staged", f.Path)))
			if err != nil {
				return err
			}
			if err := os.WriteFile(f.Path, content, 0644); err != nil {
				return err
			}
			if err := git.AddFile(f.Path); err != nil {
				return err
			}
		}

		if f.Unstaged {
			content, err := r.GetFileContent(filepath.ToSlash(filepath.Join("files", "unstaged", f.Path)))
			if err != nil {
				return err
			}
			if err := os.WriteFile(f.Path, content, 0644); err != nil {
				return err
			}
		}
		fmt.Printf("  ✓ %s\n", f.Path)
	}

	fmt.Println("Done.")
	return nil
}
