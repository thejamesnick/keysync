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
	pushEnvFile string
	pushLocal   bool
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Encrypt and sync local secrets to the project",
	Long:  `Reads the local .env file, encrypts it for all authorized project keys, and saves the encrypted blob.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// 1. Load project config
		proj, err := config.LoadProjectConfig(cwd)
		if err != nil {
			return fmt.Errorf("failed to load project config: %w (try 'keysync init')", err)
		}
		if len(proj.Keys) == 0 {
			return fmt.Errorf("no keys found in project. Add one with 'keysync add-key'")
		}

		// 2. Read and parse .env file
		envPath := filepath.Join(cwd, pushEnvFile)
		envMap, err := secrets.ParseEnvFile(envPath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", pushEnvFile, err)
		}

		// 3. Create blob
		// Retrieve current user email (author) from global config if possible
		author := "unknown"
		if globalCfg, err := config.Load(); err == nil && globalCfg != nil {
			author = globalCfg.Email
		}

		blob := secrets.NewBlob(envMap, author)
		blobBytes, err := blob.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal secrets: %w", err)
		}

		// 4. Encrypt blob
		encryptedBytes, err := crypto.Encrypt(blobBytes, proj.Keys)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		// 5. Save to disk (simulating "push")
		// In the future this will upload to server. For now, it saves to .keysync/secrets.enc
		secretsPath := filepath.Join(cwd, config.ProjectConfigDir, "secrets.enc") // .keysync/secrets.enc
		if err := os.MkdirAll(filepath.Dir(secretsPath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(secretsPath, encryptedBytes, 0644); err != nil {
			return fmt.Errorf("failed to save encrypted secrets: %w", err)
		}

		fmt.Printf("  🔒  Encrypted \033[1m%d secrets\033[0m for %d recipients\n", len(envMap), len(proj.Keys))
		fmt.Printf("  💾  Saved to \033[90m%s\033[0m\n", secretsPath)
		return nil
	},
}

func init() {
	pushCmd.Flags().StringVarP(&pushEnvFile, "file", "f", ".env", "Path to the .env file to push")
	pushCmd.Flags().BoolVar(&pushLocal, "local", true, "Perform local push only (default for MVP)")

	rootCmd.AddCommand(pushCmd)
}
