package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"keysync/internal/config"

	"github.com/spf13/cobra"
)

var (
	projectName string
	addKeyMe    bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init [project-name]",
	Short:   "Initialize a new KeySync project in the current directory",
	Example: "  keysync init my-project",
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
	Use:     "add-key [key-string-or-path]",
	Short:   "Add an SSH public key to the project",
	Example: "  keysync add-key github:username\n  keysync add-key bob.pub\n  keysync add-key --me",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var keyInput string
		if len(args) > 0 {
			keyInput = args[0]
		}

		// Handle --me flag
		if addKeyMe {
			globalCfg, err := config.Load()
			if err != nil || globalCfg == nil || globalCfg.IdentityFile == "" {
				return fmt.Errorf("must be logged in to use --me. Run 'keysync signup' or 'keysync login'")
			}
			pubKeyPath := globalCfg.IdentityFile + ".pub"
			if _, err := os.Stat(pubKeyPath); err != nil {
				return fmt.Errorf("could not find your public key at %s", pubKeyPath)
			}
			fmt.Printf("  🔍  Using your identity key: \033[1m%s\033[0m\n", filepath.Base(pubKeyPath))
			keyInput = pubKeyPath
		}

		if keyInput == "" {
			return fmt.Errorf("requires a key path, github:user, or --me")
		}

		// 1. GitHub Integration (github:username)
		var keyContent string

		if strings.HasPrefix(keyInput, "github:") {
			username := strings.TrimPrefix(keyInput, "github:")
			url := fmt.Sprintf("https://github.com/%s.keys", username)
			fmt.Printf("  🔍  Fetching keys for \033[1m%s\033[0m from GitHub...\n", username)

			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to fetch keys from GitHub: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("github user '%s' not found (status %d)", username, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read github response: %w", err)
			}

			// GitHub returns multiple keys separated by newlines
			keys := strings.Split(string(body), "\n")
			addedCount := 0

			// Load project early effectively
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			proj, err := config.LoadProjectConfig(cwd)
			if err != nil {
				return err
			}
			if proj == nil {
				return fmt.Errorf("no project found. Run 'keysync init' first")
			}

			for _, k := range keys {
				k = strings.TrimSpace(k)
				if k == "" {
					continue
				}
				// Append a comment so we know it came from github
				if !strings.Contains(k, "github.com") {
					k = fmt.Sprintf("%s %s@github", k, username)
				}
				if err := proj.AddKey(k); err == nil {
					addedCount++
				}
			}

			if err := config.SaveProjectConfig(cwd, proj); err != nil {
				return err
			}

			if addedCount == 0 {
				fmt.Printf("  ⚠️  No new keys found for %s (maybe already added?)\n", username)
			} else {
				fmt.Printf("  ✅  Imported %d keys for %s\n", addedCount, username)
			}
			return nil

		} else if _, err := os.Stat(keyInput); err == nil {
			// 2. Try to read as file
			content, err := os.ReadFile(keyInput)
			if err != nil {
				return fmt.Errorf("failed to read key file: %w", err)
			}
			keyContent = string(content)
		} else {
			// 3. Treat as raw string
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
	addKeyCmd.Flags().BoolVar(&addKeyMe, "me", false, "Add your own identity key")
	rootCmd.AddCommand(addKeyCmd)
	rootCmd.AddCommand(removeKeyCmd)
}
