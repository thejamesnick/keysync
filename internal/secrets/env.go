package secrets

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseEnvFile reads a file and parses it as a map of key-value pairs.
// It supports basic KEY=VALUE syntax and ignores comments (#).
func ParseEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	secrets := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Warn or skip? For strictness, let's error.
			// But for flexibility, skipping might be better.
			// Let's return error to help user debug bad .env
			return nil, fmt.Errorf("line %d: invalid format (expected KEY=VALUE)", lineNum)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Basic quote removal
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		secrets[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return secrets, nil
}

// WriteEnvFile writes map of secrets to a file in KEY=VALUE format
func WriteEnvFile(path string, secrets map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range secrets {
		// Basic escaping if needed (skip for MVP)
		_, err := fmt.Fprintf(f, "%s=%s\n", k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
