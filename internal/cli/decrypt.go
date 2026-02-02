package cli

import (
	"fmt"
	"os"

	"keysync/internal/config"
	"keysync/internal/crypto"

	"github.com/spf13/cobra"
)

var (
	decryptOutput   string
	decryptIdentity string
)

var decryptCmd = &cobra.Command{
	Use:    "decrypt [file]",
	Short:  "Decrypt a file using an SSH identity",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]

		// Read input file
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		// Decrypt
		identity := decryptIdentity
		if identity == "" {
			// Try to load from config
			cfg, err := config.Load()
			if err == nil && cfg != nil && cfg.IdentityFile != "" {
				identity = cfg.IdentityFile
			}
		}

		if identity == "" {
			return fmt.Errorf("identity key not specified. Use --identity or run 'keysync signup'")
		}

		decryptedData, err := crypto.Decrypt(data, identity)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		// Write output
		if decryptOutput != "" {
			if err := os.WriteFile(decryptOutput, decryptedData, 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Printf("Decrypted %s -> %s\n", inputFile, decryptOutput)
		} else {
			// Print to stdout
			fmt.Print(string(decryptedData))
		}

		return nil
	},
}

func init() {
	decryptCmd.Flags().StringVarP(&decryptOutput, "output", "o", "", "Output file path (default: stdout)")
	decryptCmd.Flags().StringVarP(&decryptIdentity, "identity", "i", "", "Path to SSH private key (optional if logged in)")

	rootCmd.AddCommand(decryptCmd)
}
