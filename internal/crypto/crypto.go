package crypto

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"filippo.io/age/agessh"
)

// Encrypt encrypts the given data for the list of SSH public keys (reipients).
// It returns the encrypted binary blob.
func Encrypt(data []byte, sshPublicKeys []string) ([]byte, error) {
	var recipients []age.Recipient

	for _, pubKey := range sshPublicKeys {
		r, err := agessh.ParseRecipient(pubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key '%s': %w", pubKey, err)
		}
		recipients = append(recipients, r)
	}

	// Create a buffer to write the encrypted data to
	out := &bytes.Buffer{}

	// Create the encryption writer
	w, err := age.Encrypt(out, recipients...)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryption writer: %w", err)
	}

	// Write the data
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	// Close to finish encryption
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close encryption writer: %w", err)
	}

	return out.Bytes(), nil
}

// Decrypt decrypts the given data using the SSH private key at the specified path.
func Decrypt(encryptedData []byte, privateKeyPath string) ([]byte, error) {
	// Read the private key file
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse the SSH identity
	// Note: agessh.ParseIdentity parses a PEM-encoded private key.
	// It handles both RSA and Ed25519 if they are in the correct format.
	// If the key is encrypted (has a passphrase), agessh.ParseEncryptedIdentity might be needed,
	// checking against the complexity, we'll start with unencrypted identity first or let ParseIdentity handle it if it supports prompts (it usually doesn't).
	// For MVP, we assume unencrypted keys or handling standard key formats.
	identity, err := agessh.ParseIdentity(keyBytes)
	if err != nil {
		// Attempt to see if it's an encrypted key issue or format issue
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create the decryption reader
	r, err := age.Decrypt(bytes.NewReader(encryptedData), identity)
	if err != nil {
		return nil, fmt.Errorf("failed to create decryption reader: %w", err)
	}

	// Read all decrypted data
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, fmt.Errorf("failed to read decrypted data: %w", err)
	}

	return out.Bytes(), nil
}

// SSHKey represents a found public key
type SSHKey struct {
	Path    string
	Content string
}

// FindSSHKeys looks for standard public keys in ~/.ssh
func FindSSHKeys() ([]SSHKey, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshDir := filepath.Join(home, ".ssh")
	files, err := os.ReadDir(sshDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var keys []SSHKey
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".pub") {
			path := filepath.Join(sshDir, f.Name())
			content, err := os.ReadFile(path)
			if err == nil {
				keys = append(keys, SSHKey{
					Path:    path,
					Content: string(content),
				})
			}
		}
	}
	return keys, nil
}
