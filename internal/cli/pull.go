package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"keysync/internal/config"
	"keysync/internal/crypto"
	"keysync/internal/secrets"

	"github.com/spf13/cobra"
)

var (
	pullTargetFile string
	pullLocal      bool
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Decrypt and update local secrets from the project",
	Long:  `Reads the encrypted secrets blob, decrypts it using your identity, and writes to .env.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// 1. Identify user identity
		globalCfg, err := config.Load()
		if err != nil || globalCfg == nil || globalCfg.IdentityFile == "" {
			return fmt.Errorf("you must be logged in to pull secrets (run 'keysync signup' or 'keysync login')")
		}

		// 2. Locate encrypted blob
		// For local MVP, look in .keysync/secrets.enc
		secretsPath := filepath.Join(cwd, config.ProjectConfigDir, "secrets.enc")
		if _, err := os.Stat(secretsPath); os.IsNotExist(err) {
			return fmt.Errorf("no secrets found at %s. Run 'keysync push' first", secretsPath)
		}

		encryptedData, err := os.ReadFile(secretsPath)
		if err != nil {
			return fmt.Errorf("failed to read secrets file: %w", err)
		}

		// 3. Decrypt

		decryptedData, err := crypto.Decrypt(encryptedData, globalCfg.IdentityFile)
		if err != nil {
			// Friendly error for common failure
			return fmt.Errorf("decryption failed: %w (Are you authorized for this project?)", err)
		}

		// 4. Parse Blob
		blob, err := secrets.Unmarshal(decryptedData)
		if err != nil {
			return fmt.Errorf("invalid secret format: %w", err)
		}

		// 5. Write .env
		targetPath := filepath.Join(cwd, pullTargetFile)
		if err := secrets.WriteEnvFile(targetPath, blob.Secrets); err != nil {
			return fmt.Errorf("failed to write .env file: %w", err)
		}

		fmt.Printf("  🔓  Decrypted with \033[90m%s\033[0m\n", filepath.Base(globalCfg.IdentityFile))
		fmt.Printf("  ✅  Pulled \033[1m%d secrets\033[0m to %s\n", len(blob.Secrets), targetPath)
		fmt.Printf("      \033[90mUpdated by %s at %s\033[0m\n", blob.Author, blob.Timestamp.Format("15:04:05"))
		return nil
	},
}

func init() {
	pullCmd.Flags().StringVarP(&pullTargetFile, "output", "o", ".env", "File to write decrypted secrets to")
	pullCmd.Flags().BoolVar(&pullLocal, "local", true, "Perform local pull only (default for MVP)")

	rootCmd.AddCommand(pullCmd)
}
