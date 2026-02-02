package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"keysync/internal/config"
	"keysync/internal/secrets"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project status and configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		proj, err := config.LoadProjectConfig(cwd)
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		// 1. Not Initialized State
		if proj == nil {
			fmt.Println("\n  KeySync is not initialized here.")
			fmt.Println("  Run \033[1mkeysync init\033[0m to start a project.")
			return nil
		}

		// 2. Initialized State - Apple Style Header
		fmt.Println()
		fmt.Printf("  �  \033[1m%s\033[0m  \033[90m%s\033[0m\n", proj.Name, cwd)
		fmt.Println("  ────────────────────────────────────────")

		// 3. Stats Grid
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

		// Check file status
		secretsPath := filepath.Join(cwd, config.ProjectConfigDir, "secrets.enc")
		hasSecrets := "No"
		if _, err := os.Stat(secretsPath); err == nil {
			hasSecrets = "Yes"
		}

		// Env count
		envCount := 0
		envPath := filepath.Join(cwd, ".env")
		if envMap, err := secrets.ParseEnvFile(envPath); err == nil {
			envCount = len(envMap)
		}

		fmt.Fprintf(w, "  \033[90mStatus\033[0m\tActive\n")
		fmt.Fprintf(w, "  \033[90mSecrets\033[0m\t%s\n", hasSecrets) // Simple yes/no for now
		fmt.Fprintf(w, "  \033[90mLocal\033[0m\t%d variables (.env)\n", envCount)
		fmt.Fprintf(w, "  \033[90mKeys\033[0m\t%d developers\n", len(proj.Keys))
		w.Flush()

		fmt.Println()

		// 4. Access Keys List (Clean & Subtle)
		if len(proj.Keys) > 0 {
			fmt.Println("  \033[1mAccess Keys\033[0m")
			for _, key := range proj.Keys {
				// Parse comment from key if possible (ssh-ed25519 AAAA... comment)
				parts := strings.Split(strings.TrimSpace(key), " ")
				keyType := parts[0]
				fingerprint := "..." + parts[1][len(parts[1])-10:] // Last 10 chars

				comment := ""
				if len(parts) > 2 {
					comment = strings.Join(parts[2:], " ")
				}

				// Format:   • user@machine (ed25519) ...X5d9A
				fmt.Printf("  \033[32m•\033[0m %-20s \033[90m%s %s\033[0m\n", comment, keyType, fingerprint)
			}
		} else {
			fmt.Println("  ⚠️  No keys added. Run \033[1mkeysync add-key\033[0m")
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
