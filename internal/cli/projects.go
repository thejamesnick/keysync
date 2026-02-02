package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"keysync/internal/config"

	"github.com/spf13/cobra"
)

var (
	projectName string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new KeySync project in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Check if already initialized
		exists, err := config.IsProjectInitialized(cwd)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("keysync is already initialized in this directory")
		}

		// Defaults
		if projectName == "" {
			projectName = filepath.Base(cwd)
		}

		// Create project config
		proj := &config.ProjectConfig{
			Name: projectName,
			Keys: []string{},
		}

		// Optionally auto-add the current user's key if they are logged in.
		// We can check global config.
		globalCfg, err := config.Load()
		if err == nil && globalCfg != nil && globalCfg.IdentityFile != "" {
			// Try to find the associated public key for the identity file
			// This is a naive guess: private key path + ".pub"
			pubKeyPath := globalCfg.IdentityFile + ".pub"
			pubBytes, err := os.ReadFile(pubKeyPath)
			if err == nil {
				proj.Keys = append(proj.Keys, string(pubBytes))
				fmt.Printf("✨ Auto-added your public key (%s)\n", filepath.Base(pubKeyPath))
			}
		}

		if err := config.SaveProjectConfig(cwd, proj); err != nil {
			return fmt.Errorf("failed to save project config: %w", err)
		}

		// Update .gitignore
		gitignorePath := filepath.Join(cwd, ".gitignore")
		f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			// Check if .env is already ignored (basic check)
			content, _ := os.ReadFile(gitignorePath)
			if !strings.Contains(string(content), ".env") {
				if _, err := f.WriteString("\n# KeySync\n.env\n"); err != nil {
					fmt.Printf("⚠️  Failed to update .gitignore: %v\n", err)
				} else {
					fmt.Println("📝 Added .env to .gitignore")
				}
			}
		}

		fmt.Printf("\n  🚀  Initialized project \033[1m%s\033[0m\n", proj.Name)
		fmt.Printf("  📄  Config: %s\n", config.ProjectConfigFileName)
		return nil
	},
}

var addKeyCmd = &cobra.Command{
	Use:   "add-key [key-string-or-path]",
	Short: "Add an SSH public key to the project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyInput := args[0]

		// 1. Try to read as file
		var keyContent string
		if _, err := os.Stat(keyInput); err == nil {
			content, err := os.ReadFile(keyInput)
			if err != nil {
				return fmt.Errorf("failed to read key file: %w", err)
			}
			keyContent = string(content)
		} else {
			// 2. Treat as raw string
			keyContent = keyInput
		}

		// Validate key roughly (must be ssh-...)
		// In a real app we might parse it with golang.org/x/crypto/ssh
		// For now simple check
		// (omitted for brevity, we trust the user implies a key)

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		proj, err := config.LoadProjectConfig(cwd)
		if err != nil {
			return fmt.Errorf("failed to load project config: %w", err)
		}
		if proj == nil {
			return fmt.Errorf("no project found. Run 'keysync init' first")
		}

		if err := proj.AddKey(keyContent); err != nil {
			return err
		}

		if err := config.SaveProjectConfig(cwd, proj); err != nil {
			return err
		}

		fmt.Printf("  ✅  Added key: \033[90m%s...\033[0m\n", keyContent[:20])
		return nil
	},
}

var removeKeyCmd = &cobra.Command{
	Use:   "remove-key [key-string-or-part]",
	Short: "Remove an SSH public key from the project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyInput := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		proj, err := config.LoadProjectConfig(cwd)
		if err != nil {
			return fmt.Errorf("failed to load project config: %w", err)
		}
		if proj == nil {
			return fmt.Errorf("no project found. Run 'keysync init' first")
		}

		// Simple matching: exact match or contains?
		// User might paste the whole key, or just a comment part?
		// For safety, let's require strict matching or maybe simple substring if unambiguous.
		// For now, let's try exact match on the stored string (which often has newlines stripped or kept).
		// Better: iterate and match.

		// Plan implementation used "remove-key".
		// We'll trust exact content or try to find a single match.
		// Let's rely on exact match for now as implemented in config/project.go
		if err := proj.RemoveKey(keyInput); err != nil {
			return fmt.Errorf("failed to remove key: %w (ensure exact match)", err)
		}

		if err := config.SaveProjectConfig(cwd, proj); err != nil {
			return err
		}

		fmt.Println("  🗑️   Key removed from project.")
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&projectName, "name", "", "Name of the project (default: current directory name)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addKeyCmd)
	rootCmd.AddCommand(removeKeyCmd)
}
