package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"keysync/internal/config"
	"keysync/internal/crypto"

	"github.com/spf13/cobra"
)

var (
	signupEmail string
	signupKey   string
	signupMe    bool
)

var signupCmd = &cobra.Command{
	Use:     "signup",
	Short:   "Create a new KeySync account (local config for now)",
	Example: "  keysync signup --email me@example.com --me\n  keysync signup --email me@example.com --key ~/.ssh/id_ed25519.pub",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle --me flag logic
		if signupMe {
			keys, err := crypto.FindSSHKeys()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				return fmt.Errorf("no keys found in ~/.ssh/. Generate one with 'keysync generate'")
			}
			// Use the first found key (usually id_ed25519.pub if sorted or standard)
			// Ideally we prioritize id_ed25519.pub
			var selectedKey string
			for _, k := range keys {
				if strings.Contains(k.Path, "id_ed25519.pub") {
					selectedKey = k.Path
					break
				}
			}
			if selectedKey == "" {
				selectedKey = keys[0].Path // Fallback to first
			}
			signupKey = selectedKey
			fmt.Printf("  🔍  Found identity: \033[1m%s\033[0m\n", filepath.Base(signupKey))
		}

		// validate inputs
		if signupEmail == "" || signupKey == "" {
			return fmt.Errorf("email and key (or --me) are required")
		}

		// expand home dir in key path if needed
		if strings.HasPrefix(signupKey, "~/") {
			home, _ := os.UserHomeDir()
			signupKey = filepath.Join(home, signupKey[2:])
		}

		keyPath, err := filepath.Abs(signupKey)
		if err != nil {
			return err
		}

		// check if public key exists (we expect public key for signup identification)
		// Usually signup takes pubkey, but config needs private key for decryption?
		// "Keysync uses SSH keys as identity... strictly tied to SSH keys."
		// For SIGNUP, we provide the public key to the server.
		// For CONFIG, we need to know where the PRIVATE key is to decrypt stuff later.
		// Usually they are side-by-side (id_ed25519 and id_ed25519.pub).
		// Let's assume the user points to their PUBLIC key for signup.
		// We infer the identity file (private key) is the same path without .pub?
		// Or we ask for identity file separately?
		// The prompt says: "signup --email ... --key ~/.ssh/id_ed25519.pub"
		// The 'login' command validates we own it.

		// For local config, we want to store the Identity File (private key path) so 'decrypt' works automatically.
		// Let's try to infer private key from public key path.
		identityFile := keyPath
		if strings.HasSuffix(keyPath, ".pub") {
			identityFile = strings.TrimSuffix(keyPath, ".pub")
		}

		// Check if identity file exists
		if _, err := os.Stat(identityFile); os.IsNotExist(err) {
			return fmt.Errorf("could not find private key at %s", identityFile)
		}

		cfg := &config.Config{
			Email:        signupEmail,
			IdentityFile: identityFile,
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("  ✅  Account created for \033[1m%s\033[0m\n", signupEmail)
		fmt.Printf("  🔑  Identity: %s.pub\n", identityFile)
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to KeySync",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if cfg == nil {
			return fmt.Errorf("no account found. Run 'keysync signup' first")
		}

		fmt.Printf("  ✨  Logged in as \033[1m%s\033[0m\n", cfg.Email)
		// Future: server auth challenge
		return nil
	},
}

func init() {
	signupCmd.Flags().StringVar(&signupEmail, "email", "", "Your email address")
	signupCmd.Flags().StringVar(&signupKey, "key", "", "Path to your SSH public key")
	signupCmd.Flags().BoolVar(&signupMe, "me", false, "Use your local identity key automatically")
	signupCmd.MarkFlagRequired("email")
	// signupCmd.MarkFlagRequired("key") // Removed so --me can replace it

	rootCmd.AddCommand(signupCmd)
	rootCmd.AddCommand(loginCmd)
}
