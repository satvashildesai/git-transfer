package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/git-transfer/git-transfer/internal/manifest"
)

// Writer creates a new GTB bundle
type Writer struct {
	zipWriter *zip.Writer
	file      *os.File
}

// NewWriter initializes a bundle writer
func NewWriter(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create bundle file: %w", err)
	}
	
	return &Writer{
		zipWriter: zip.NewWriter(f),
		file:      f,
	}, nil
}

// Close finalizes the bundle
func (w *Writer) Close() error {
	if err := w.zipWriter.Close(); err != nil {
		w.file.Close()
		return err
	}
	return w.file.Close()
}

// WriteManifest writes the manifest.json to the bundle
func (w *Writer) WriteManifest(m *manifest.Manifest) error {
	f, err := w.zipWriter.Create("manifest.json")
	if err != nil {
		return err
	}
	return m.Write(f)
}

// AddFile adds a file from the filesystem (unstaged) to the bundle
func (w *Writer) AddUnstagedFile(repoPath string) error {
	srcFile, err := os.Open(repoPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destPath := filepath.ToSlash(filepath.Join("files", "unstaged", repoPath))
	wFile, err := w.zipWriter.Create(destPath)
	if err != nil {
		return err
	}
	
	_, err = io.Copy(wFile, srcFile)
	return err
}

// AddStagedFile adds raw content as a staged file to the bundle
func (w *Writer) AddStagedFile(repoPath string, content []byte) error {
	destPath := filepath.ToSlash(filepath.Join("files", "staged", repoPath))
	wFile, err := w.zipWriter.Create(destPath)
	if err != nil {
		return err
	}
	
	_, err = wFile.Write(content)
	return err
}
