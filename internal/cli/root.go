package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "keysync",
	Short: "KeySync: The SSH-Native Secret Manager",
	Long: `Sync your secrets securely, SSH-style.
KeySync uses SSH keys to encrypt and manage secrets for your team.
Zero knowledge, local-first.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
