package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bigdra50/doc/internal/config"
)

func main() {
	// Parse command line arguments
	cliArgs, err := parseArgs()
	if err != nil {
		showUsage()
		os.Exit(1)
	}

	// Set global verbose flag
	verbose = cliArgs.Verbose

	// Handle special commands
	if handleSpecialCommands(cliArgs) {
		return
	}

	// Handle merge command
	if cliArgs.IsMergeCommand {
		if err := runMerge(cliArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Run translation
	if err := runTranslation(cliArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// handleSpecialCommands handles configuration and listing commands
func handleSpecialCommands(cliArgs *CLIArgs) bool {
	// Handle config commands
	if cliArgs.ShowConfig {
		showCurrentConfig()
		return true
	}

	if cliArgs.InitConfig {
		initConfigFile()
		return true
	}

	if len(cliArgs.SetConfig) > 0 {
		setConfigValues(cliArgs.SetConfig)
		return true
	}

	// Handle list commands
	if cliArgs.ShowList {
		showSupportedLanguages()
		return true
	}

	if cliArgs.ShowListModels {
		if cliArgs.ListModelsProvider != "" {
			showModelsForProvider(cliArgs.ListModelsProvider)
		} else {
			showAllModels()
		}
		return true
	}

	return false
}

// runTranslation performs the main translation operation
func runTranslation(cliArgs *CLIArgs) error {
	// Load configuration
	config := LoadConfig()
	config.Verbose = verbose

	if verbose {
		log("Configuration: Provider=%s, OpenAI=%s, Anthropic=%s",
			config.ProviderType,
			maskAPIKey(config.OpenAIAPIKey),
			maskAPIKey(config.AnthropicAPIKey))
	}

	// Create LLM provider
	provider, err := NewLLMProvider(config)
	if err != nil {
		showProviderHelp(config.ProviderType)
		return fmt.Errorf("failed to initialize %s provider: %w", config.ProviderType, err)
	}

	log("Using provider: %s", provider.GetProviderName())

	// Validate language code
	if err := validateLanguage(cliArgs.TargetLanguage, provider); err != nil {
		return err
	}

	log("Target language: %s", cliArgs.TargetLanguage)
	if cliArgs.TransformInstruction != "" {
		log("Custom instruction: %s", cliArgs.TransformInstruction)
	}

	// Read document from stdin
	content, err := readDocument()
	if err != nil {
		return err
	}

	// Perform translation
	result, err := performTranslation(provider, content, cliArgs.TargetLanguage, cliArgs.TransformInstruction)
	if err != nil {
		return fmt.Errorf("translation failed: %w", err)
	}

	// Output the translation result
	fmt.Print(result)
	return nil
}

// validateLanguage validates the target language code
func validateLanguage(targetLang string, provider LLMProvider) error {
	supportedLangs := provider.GetSupportedLanguages()
	if err := validateLanguageCodeWithMap(targetLang, supportedLangs); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		// Show similar codes if available
		similar := getSimilarLanguageCodesWithMap(targetLang, supportedLangs)
		if len(similar) > 0 {
			fmt.Fprintf(os.Stderr, "\nDid you mean:\n")
			for _, code := range similar {
				fmt.Fprintf(os.Stderr, "  %s - %s\n", code, supportedLangs[code])
			}
		}

		fmt.Fprintf(os.Stderr, "\nUse 'doc --list' to see all supported language codes.\n")
		return err
	}
	return nil
}

// showCurrentConfig displays the current configuration
func showCurrentConfig() {
	cfg := LoadConfig()
	fmt.Printf("Current Configuration:\n")
	fmt.Printf("Config file: %s\n", config.GetConfigPath())
	fmt.Printf("\n")
	fmt.Printf("provider = \"%s\"\n", cfg.ProviderType)
	fmt.Printf("claude_code_path = \"%s\"\n", cfg.ClaudeCodePath)
	fmt.Printf("openai_model = \"%s\"\n", cfg.OpenAIModel)
	fmt.Printf("anthropic_model = \"%s\"\n", cfg.AnthropicModel)
	fmt.Printf("claude_model = \"%s\"\n", cfg.ClaudeModel)
	fmt.Printf("openai_api_key = \"%s\"\n", maskAPIKey(cfg.OpenAIAPIKey))
	fmt.Printf("anthropic_api_key = \"%s\"\n", maskAPIKey(cfg.AnthropicAPIKey))
}

// initConfigFile creates a default configuration file
func initConfigFile() {
	configPath := config.GetConfigPath()

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration file already exists at: %s\n", configPath)
		fmt.Printf("Use 'doc --config' to view current settings\n")
		return
	}

	// Create default config
	defaultConfig := config.Config{
		ProviderType:   config.ProviderTypeClaude,
		ClaudeCodePath: "claude",
		OpenAIModel:    config.GetDefaultModel(config.ProviderTypeOpenAI),
		AnthropicModel: config.GetDefaultModel(config.ProviderTypeAnthropic),
		ClaudeModel:    config.GetDefaultModel(config.ProviderTypeClaude),
	}

	if err := config.SaveConfig(defaultConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created default configuration file at: %s\n", configPath)
	fmt.Printf("Use 'doc --config' to view settings\n")
	fmt.Printf("Use 'doc --set key=value' to modify settings\n")
}

// setConfigValues updates configuration values
func setConfigValues(keyValuePairs []string) {
	// Load current config
	currentConfig := LoadConfig()

	// Parse and apply changes
	for _, pair := range keyValuePairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Use key=value format.\n", pair)
			os.Exit(1)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "provider":
			if value != config.ProviderTypeClaude && value != config.ProviderTypeOpenAI && value != config.ProviderTypeAnthropic {
				fmt.Fprintf(os.Stderr, "Error: Invalid provider '%s'. Must be one of: claude-code, openai, anthropic\n", value)
				os.Exit(1)
			}
			currentConfig.ProviderType = value
		case "openai_api_key":
			currentConfig.OpenAIAPIKey = value
		case "anthropic_api_key":
			currentConfig.AnthropicAPIKey = value
		case "claude_code_path":
			currentConfig.ClaudeCodePath = value
		case "openai_model":
			currentConfig.OpenAIModel = value
		case "anthropic_model":
			currentConfig.AnthropicModel = value
		case "claude_model":
			currentConfig.ClaudeModel = value
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown configuration key '%s'\n", key)
			fmt.Fprintf(os.Stderr, "Valid keys: provider, openai_api_key, anthropic_api_key, claude_code_path, openai_model, anthropic_model, claude_model\n")
			os.Exit(1)
		}

		fmt.Printf("Set %s = %s\n", key, maskConfigValue(key, value))
	}

	// Save updated config
	if err := config.SaveConfig(currentConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration updated successfully\n")
}

// maskConfigValue masks sensitive configuration values for display
func maskConfigValue(key, value string) string {
	if strings.Contains(key, "api_key") && value != "" {
		return maskAPIKey(value)
	}
	return value
}
