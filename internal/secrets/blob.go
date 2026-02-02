package secrets

import (
	"encoding/json"
	"time"
)

// Blob represents the unencrypted structure of the secret data
// containing the actual secrets and metadata.
type Blob struct {
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Author    string            `json:"author"` // Email or ID of the user who created this
	Secrets   map[string]string `json:"secrets"`
}

// NewBlob creates a new Blob from a map of secrets and author
func NewBlob(secrets map[string]string, author string) *Blob {
	return &Blob{
		Version:   "v1",
		Timestamp: time.Now(),
		Author:    author,
		Secrets:   secrets,
	}
}

// Marshal converts the blob to JSON bytes ready for encryption
func (b *Blob) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

// Unmarshal parses the JSON bytes back into a Blob
func Unmarshal(data []byte) (*Blob, error) {
	var b Blob
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
