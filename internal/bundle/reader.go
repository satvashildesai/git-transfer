package bundle

import (
	"archive/zip"
	"fmt"
	"io"

	"github.com/git-transfer/git-transfer/internal/manifest"
)

// Reader reads a GTB bundle
type Reader struct {
	zipReader *zip.ReadCloser
}

// NewReader opens a bundle for reading
func NewReader(path string) (*Reader, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open bundle: %w", err)
	}
	return &Reader{zipReader: zr}, nil
}

// Close closes the bundle reader
func (r *Reader) Close() error {
	return r.zipReader.Close()
}

// ReadManifest extracts and parses the manifest.json
func (r *Reader) ReadManifest() (*manifest.Manifest, error) {
	for _, f := range r.zipReader.File {
		if f.Name == "manifest.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return manifest.Read(rc)
		}
	}
	return nil, fmt.Errorf("manifest.json not found in bundle")
}

// GetFileContent retrieves the content of a file from the bundle
func (r *Reader) GetFileContent(internalPath string) ([]byte, error) {
	for _, f := range r.zipReader.File {
		if f.Name == internalPath {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("file %s not found in bundle", internalPath)
}
