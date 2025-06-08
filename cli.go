package main

import (
	"fmt"
	"os"
)

// CLIArgs represents parsed command line arguments
type CLIArgs struct {
	Verbose              bool
	TargetLanguage       string
	TransformInstruction string
	ShowList             bool
	ShowListModels       bool
	ListModelsProvider   string
	ShowConfig           bool
	SetConfig            []string // Key=value pairs
	InitConfig           bool
}

// parseArgs parses command line arguments and returns CLIArgs
func parseArgs() (*CLIArgs, error) {
	args := os.Args[1:]
	cliArgs := &CLIArgs{}

	// Handle verbose flag
	if len(args) > 0 && args[0] == "-v" {
		cliArgs.Verbose = true
		args = args[1:]
		log("Verbose mode enabled")
	}

	if len(args) < 1 {
		return nil, fmt.Errorf("missing required arguments")
	}

	// Handle --list options
	if args[0] == "--list" {
		cliArgs.ShowList = true
		return cliArgs, nil
	}

	if args[0] == "--list-models" {
		cliArgs.ShowListModels = true
		if len(args) > 1 {
			cliArgs.ListModelsProvider = args[1]
		}
		return cliArgs, nil
	}

	// Handle config commands
	if args[0] == "--config" {
		cliArgs.ShowConfig = true
		return cliArgs, nil
	}

	if args[0] == "--init-config" {
		cliArgs.InitConfig = true
		return cliArgs, nil
	}

	if args[0] == "--set" {
		if len(args) < 2 {
			return nil, fmt.Errorf("--set requires key=value pairs")
		}
		cliArgs.SetConfig = args[1:]
		return cliArgs, nil
	}

	// Parse target language and optional transform instruction
	cliArgs.TargetLanguage = args[0]
	if len(args) > 1 {
		cliArgs.TransformInstruction = args[1]
	}

	return cliArgs, nil
}

// showUsage displays the usage information
func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: doc [-v] <language_code> [transform_instruction]\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  cat README.md | doc ja\n")
	fmt.Fprintf(os.Stderr, "  cat README.md | doc -v ru\n")
	fmt.Fprintf(os.Stderr, "  doc --list          # Show supported language codes\n")
	fmt.Fprintf(os.Stderr, "  doc --list-models   # Show all available models\n")
	fmt.Fprintf(os.Stderr, "  doc --list-models openai # Show OpenAI models only\n")
	fmt.Fprintf(os.Stderr, "\nConfiguration Commands:\n")
	fmt.Fprintf(os.Stderr, "  doc --config        # Show current configuration\n")
	fmt.Fprintf(os.Stderr, "  doc --init-config   # Create default config file\n")
	fmt.Fprintf(os.Stderr, "  doc --set provider=openai # Set configuration value\n")
	fmt.Fprintf(os.Stderr, "  doc --set openai_api_key=sk-... # Set API key\n")
	fmt.Fprintf(os.Stderr, "\nEnvironment Variables (override config file):\n")
	fmt.Fprintf(os.Stderr, "  LLM_PROVIDER      - Provider type: claude-code, openai, anthropic (default: claude-code)\n")
	fmt.Fprintf(os.Stderr, "  OPENAI_API_KEY    - OpenAI API key (required for openai provider)\n")
	fmt.Fprintf(os.Stderr, "  ANTHROPIC_API_KEY - Anthropic API key (required for anthropic provider)\n")
	fmt.Fprintf(os.Stderr, "  OPENAI_MODEL      - OpenAI model to use (default: gpt-4o-mini)\n")
	fmt.Fprintf(os.Stderr, "  ANTHROPIC_MODEL   - Anthropic model to use (default: claude-3-5-haiku-20241022)\n")
	fmt.Fprintf(os.Stderr, "  CLAUDE_MODEL      - Claude Code model to use (default: sonnet)\n")
	fmt.Fprintf(os.Stderr, "\nConfig File: $XDG_CONFIG_HOME/bigdra50/doc/config.toml (or ~/.config/bigdra50/doc/config.toml)\n")
}

// showProviderHelp displays provider-specific help information
func showProviderHelp(providerType string) {
	switch providerType {
	case ProviderTypeOpenAI:
		fmt.Fprintf(os.Stderr, "\nOpenAI Provider Help:\n")
		fmt.Fprintf(os.Stderr, "  Set OPENAI_API_KEY environment variable with your OpenAI API key\n")
		fmt.Fprintf(os.Stderr, "  Example: export OPENAI_API_KEY=sk-...\n")
	case ProviderTypeAnthropic:
		fmt.Fprintf(os.Stderr, "\nAnthropic Provider Help:\n")
		fmt.Fprintf(os.Stderr, "  Set ANTHROPIC_API_KEY environment variable with your Anthropic API key\n")
		fmt.Fprintf(os.Stderr, "  Example: export ANTHROPIC_API_KEY=sk-ant-...\n")
	case ProviderTypeClaude:
		fmt.Fprintf(os.Stderr, "\nClaude Code Provider Help:\n")
		fmt.Fprintf(os.Stderr, "  Ensure Claude Code CLI is installed and available in PATH\n")
		fmt.Fprintf(os.Stderr, "  Install: npm install -g @anthropic-ai/claude-code\n")
	}
}

// showAllModels displays all available models
func showAllModels() {
	fmt.Fprintf(os.Stderr, "Available Models:\n\n")

	catalog := GetModelCatalog()

	fmt.Fprintf(os.Stderr, "OpenAI Models:\n")
	for _, model := range catalog.OpenAI {
		fmt.Fprintf(os.Stderr, "  %-25s %s (tier: %s, cost: $%.2f/$%.2f per 1M tokens)\n",
			model.ID, model.Name, model.Tier, model.InputCostPer1M, model.OutputCostPer1M)
	}

	fmt.Fprintf(os.Stderr, "\nAnthropic Models:\n")
	for _, model := range catalog.Anthropic {
		fmt.Fprintf(os.Stderr, "  %-25s %s (tier: %s, cost: $%.2f/$%.2f per 1M tokens)\n",
			model.ID, model.Name, model.Tier, model.InputCostPer1M, model.OutputCostPer1M)
	}

	fmt.Fprintf(os.Stderr, "\nClaude Code Models:\n")
	fmt.Fprintf(os.Stderr, "  %-25s %s\n", "opus", "Claude Opus (high capability)")
	fmt.Fprintf(os.Stderr, "  %-25s %s\n", "sonnet", "Claude Sonnet (balanced)")
	fmt.Fprintf(os.Stderr, "  %-25s %s\n", "haiku", "Claude Haiku (fast)")
}

// showModelsForProvider displays models for a specific provider
func showModelsForProvider(provider string) {
	switch provider {
	case "openai":
		fmt.Fprintf(os.Stderr, "OpenAI Models:\n")
		for _, model := range GetModelsByProvider(ProviderTypeOpenAI) {
			fmt.Fprintf(os.Stderr, "  %-25s %s (tier: %s)\n", model.ID, model.Name, model.Tier)
			fmt.Fprintf(os.Stderr, "    Cost: $%.2f input / $%.2f output per 1M tokens\n",
				model.InputCostPer1M, model.OutputCostPer1M)
			fmt.Fprintf(os.Stderr, "    Context: %d tokens\n", model.ContextWindow)
			fmt.Fprintf(os.Stderr, "    Best for: %v\n\n", model.RecommendedFor)
		}
	case "anthropic":
		fmt.Fprintf(os.Stderr, "Anthropic Models:\n")
		for _, model := range GetModelsByProvider(ProviderTypeAnthropic) {
			fmt.Fprintf(os.Stderr, "  %-25s %s (tier: %s)\n", model.ID, model.Name, model.Tier)
			fmt.Fprintf(os.Stderr, "    Cost: $%.2f input / $%.2f output per 1M tokens\n",
				model.InputCostPer1M, model.OutputCostPer1M)
			fmt.Fprintf(os.Stderr, "    Context: %d tokens\n", model.ContextWindow)
			fmt.Fprintf(os.Stderr, "    Best for: %v\n\n", model.RecommendedFor)
		}
	case "claude-code":
		fmt.Fprintf(os.Stderr, "Claude Code Models:\n")
		fmt.Fprintf(os.Stderr, "  %-25s %s\n", "opus", "High capability, best performance")
		fmt.Fprintf(os.Stderr, "  %-25s %s\n", "sonnet", "Balanced performance and speed (default)")
		fmt.Fprintf(os.Stderr, "  %-25s %s\n", "haiku", "Fast response, lower cost")
	default:
		fmt.Fprintf(os.Stderr, "Unknown provider: %s\n", provider)
		fmt.Fprintf(os.Stderr, "Available providers: openai, anthropic, claude-code\n")
	}
}
