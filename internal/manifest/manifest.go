package manifest

import (
	"encoding/json"
	"io"
)

// Version is the current manifest format version
const Version = 1

// Manifest represents the structure of the changes.gtb metadata
type Manifest struct {
	FormatVersion int        `json:"formatVersion"`
	CreatedAt     string     `json:"createdAt"`
	Git           GitInfo    `json:"git"`
	Files         []FileInfo `json:"files"`
}

type GitInfo struct {
	Head string `json:"head"`
}

type FileInfo struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	Staged   bool   `json:"staged"`
	Unstaged bool   `json:"unstaged"`
}

// Write serializes the manifest to JSON
func (m *Manifest) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(m)
}

// Read parses the manifest from JSON
func Read(r io.Reader) (*Manifest, error) {
	var m Manifest
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
