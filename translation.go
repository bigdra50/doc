package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// readDocument reads the document from stdin with validation
func readDocument() (string, error) {
	log("Checking if stdin is available...")
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("no document provided via stdin")
	}
	log("Stdin is available")

	progress("Reading document...")
	log("Reading from stdin...")
	
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}
	
	content := strings.Join(lines, "\n")
	log("Read %d characters from stdin", len(content))

	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("empty document provided")
	}

	return content, nil
}

// performTranslation performs the translation using the specified provider
func performTranslation(provider LLMProvider, content, targetLang, customInstruction string) (string, error) {
	options := TranslationOptions{
		TargetLanguage:    targetLang,
		CustomInstruction: customInstruction,
		PreserveFormat:    true,
		Verbose:          verbose,
	}

	providerName := provider.GetProviderName()
	spinner := NewSpinner(fmt.Sprintf("Translating with %s...", providerName))
	spinner.Start()

	ctx := context.Background()
	response, err := provider.Translate(ctx, content, options)
	if err != nil {
		spinner.Stop("Translation failed")
		return "", fmt.Errorf("%s translation failed: %w", providerName, err)
	}

	spinner.Stop("Translation completed")

	if response.Status != "success" {
		return "", fmt.Errorf("translation failed: %s (status: %s)", response.Message, response.Status)
	}

	return response.Content, nil
}