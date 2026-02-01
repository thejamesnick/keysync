package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestEncryptDecrypt(t *testing.T) {
	// 1. Generate a temporary Ed25519 key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ed25519 key: %v", err)
	}

	// 2. Convert to SSH public key format (string)
	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		t.Fatalf("Failed to create SSH public key: %v", err)
	}
	sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)
	sshPubKeyStr := string(sshPubKeyBytes)

	// 3. Convert to SSH private key format (PEM file) and write to temp file
	// We use the "OPENSSH PRIVATE KEY" format which is standard for modern keys
	privKeyBytes, err := ssh.MarshalPrivateKey(privKey, "") // "" = no passphrase
	if err != nil {
		t.Fatalf("Failed to marshal private key: %v", err)
	}

	block := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: privKeyBytes.Bytes, // ssh.MarshalPrivateKey returns a pem.Block-like struct in recent versions or we wrap it?
		// Wait, golang.org/x/crypto/ssh/MarshalPrivateKey returns *pem.Block.
	}

	tmpFile, err := os.CreateTemp("", "keysync_test_key")
	if err != nil {
		t.Fatalf("Failed to create temp key file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	if err := pem.Encode(tmpFile, block); err != nil {
		t.Fatalf("Failed to write private key to file: %v", err)
	}
	tmpFile.Close()

	// 4. Test Encryption
	originalMsg := []byte("The secret sauce is in the keys.")
	encrypted, err := Encrypt(originalMsg, []string{sshPubKeyStr})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(encrypted) == 0 {
		t.Fatal("Encrypted data is empty")
	}

	// 5. Test Decryption
	decrypted, err := Decrypt(encrypted, tmpFile.Name())
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(originalMsg) {
		t.Errorf("Decryption mismatch. Got %s, want %s", string(decrypted), string(originalMsg))
	}
}
