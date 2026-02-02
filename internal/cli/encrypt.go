package cli

import (
	"fmt"
	"os"
	"strings"

	"keysync/internal/crypto"

	"keysync/internal/config"

	"github.com/spf13/cobra"
)

var (
	encryptOutput     string
	encryptRecipients []string
)

var encryptCmd = &cobra.Command{
	Use:    "encrypt [file]",
	Short:  "Encrypt a file for one or more SSH recipients (defaults to project keys)",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]

		// Read input file
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		// Prepare recipients
		var finalRecipients []string

		// 1. Check flags
		if len(encryptRecipients) > 0 {
			for _, r := range encryptRecipients {
				if strings.HasPrefix(r, "ssh-") {
					finalRecipients = append(finalRecipients, r)
				} else {
					keyBytes, err := os.ReadFile(r)
					if err != nil {
						return fmt.Errorf("failed to read recipient key file '%s': %w", r, err)
					}
					finalRecipients = append(finalRecipients, string(keyBytes))
				}
			}
		} else {
			// 2. Fallback to project config
			cwd, err := os.Getwd()
			if err == nil {
				proj, err := config.LoadProjectConfig(cwd)
				if err == nil && proj != nil && len(proj.Keys) > 0 {
					finalRecipients = append(finalRecipients, proj.Keys...)
					fmt.Printf("🔒 Using %d keys from project '%s'\n", len(proj.Keys), proj.Name)
				}
			}
		}

		if len(finalRecipients) == 0 {
			return fmt.Errorf("no recipients provided and no project keys found (--recipient is required if not in a project)")
		}

		// Encrypt
		encryptedData, err := crypto.Encrypt(data, finalRecipients)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		// Determine output file
		outputFile := encryptOutput
		if outputFile == "" {
			outputFile = inputFile + ".age"
		}

		// Write output
		if err := os.WriteFile(outputFile, encryptedData, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Encrypted %s -> %s\n", inputFile, outputFile)
		return nil
	},
}

func init() {
	encryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "Output file path (default: <input>.age)")
	encryptCmd.Flags().StringSliceVarP(&encryptRecipients, "recipient", "r", nil, "SSH public key or path to public key file (required)")
	encryptCmd.MarkFlagRequired("recipient")

	rootCmd.AddCommand(encryptCmd)
}
