package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	genKeyEmail string
	genKeyName  string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate a new SSH key pair for usage with KeySync",
	Example: "  keysync generate --email me@example.com",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Determine key path
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		sshDir := filepath.Join(home, ".ssh")
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			return fmt.Errorf("failed to create ~/.ssh directory: %w", err)
		}

		keyName := genKeyName
		if keyName == "" {
			keyName = "id_ed25519"
		}
		keyPath := filepath.Join(sshDir, keyName)
		pubPath := keyPath + ".pub"

		// 2. Check if exists
		if _, err := os.Stat(keyPath); err == nil {
			return fmt.Errorf("key already exists at %s. To overwrite, delete it manualy or use a different name", keyPath)
		}

		// 3. Generate using ssh-keygen (safest way to ensure compatibility)
		// ssh-keygen -t ed25519 -C "email" -f path -N ""
		cmdGen := exec.Command("ssh-keygen", "-t", "ed25519", "-C", genKeyEmail, "-f", keyPath, "-N", "")
		cmdGen.Stdout = os.Stdout
		cmdGen.Stderr = os.Stderr

		fmt.Printf("  🎲  Generating new SSH key: \033[1m%s\033[0m\n", filepath.Base(keyPath))
		if err := cmdGen.Run(); err != nil {
			return fmt.Errorf("ssh-keygen failed: %w", err)
		}

		fmt.Println("\n  ✨  \033[1mSuccess!\033[0m")
		fmt.Printf("  🌍  Public Key:  %s\n", pubPath)
		fmt.Println("      (Share this key with your project owner)")

		fmt.Println("\n  To start using KeySync:")
		fmt.Printf("  keysync signup --email %s --key %s\n", genKeyEmail, pubPath)

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVar(&genKeyEmail, "email", "", "Email/Comment for the key (optional)")
	generateCmd.Flags().StringVar(&genKeyName, "name", "id_ed25519", "Filename for the key")

	rootCmd.AddCommand(generateCmd)
}
