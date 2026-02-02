package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"keysync/internal/crypto"

	"github.com/spf13/cobra"
)

var identifyCmd = &cobra.Command{
	Use:     "identify",
	Aliases: []string{"whoami"},
	Short:   "Show your SSH public keys (easy copy-paste)",
	Example: "  keysync identify",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Find keys in default location
		keys, err := crypto.FindSSHKeys()
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			fmt.Println("⚠️  No SSH keys found in ~/.ssh/")
			fmt.Println("   Run 'ssh-keygen -t ed25519' to generate one.")
			return nil
		}

		fmt.Println("\n  🔑  \033[1mYour Public Keys\033[0m")
		fmt.Println("  ────────────────────────────────────────")

		for _, k := range keys {
			fmt.Printf("  \033[90m%s\033[0m\n", filepath.Base(k.Path))
			// Print the key content for easy copying
			fmt.Printf("  %s\n\n", strings.TrimSpace(k.Content))
		}

		fmt.Println("  👉  Copy a key above and send it to your Project Owner.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(identifyCmd)
}
