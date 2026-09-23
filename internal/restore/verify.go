package restore

import (
	"fmt"

	"github.com/git-transfer/git-transfer/internal/bundle"
	"github.com/git-transfer/git-transfer/internal/git"
)

// Verify checks the bundle for compatibility with the current repository
func Verify(bundlePath string) error {
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
	
	fmt.Println("Bundle is valid.")

	if head == manifest.Git.Head {
		fmt.Println("Repository: compatible")
	} else {
		fmt.Println("Repository: incompatible")
	}

	fmt.Printf("Base commit: %s\n", manifest.Git.Head)
	fmt.Printf("Current HEAD: %s\n", head)

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
		fmt.Println("Potential conflicts: detected")
		for _, c := range conflicts {
			fmt.Printf("  %s\n", c)
		}
	} else {
		fmt.Println("Potential conflicts: none")
	}

	return nil
}

// Inspect prints out bundle contents
func Inspect(bundlePath string) error {
	r, err := bundle.NewReader(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to read bundle: %w", err)
	}
	defer r.Close()

	manifest, err := r.ReadManifest()
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}
	
	fmt.Println("Git Transfer Bundle")
	fmt.Println()
	fmt.Printf("Base commit:\n%s\n\n", manifest.Git.Head)
	
	if len(manifest.Files) > 0 {
		fmt.Println("Files:")
		for _, f := range manifest.Files {
			statusStr := ""
			if f.Staged && f.Unstaged {
				statusStr = "staged, unstaged"
			} else if f.Staged {
				statusStr = "staged"
			} else if f.Unstaged {
				statusStr = "unstaged"
			}
			fmt.Printf("  M %s (%s)\n", f.Path, statusStr)
		}
	}
	
	fmt.Printf("\nCreated: %s\n", manifest.CreatedAt)
	return nil
}
