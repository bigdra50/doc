package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	
	// Merge command fields
	IsMergeCommand       bool
	MergeDirectory       string
	MergeOutputFile      string
	MergeRecursive       bool
	MergeOrder           string
	MergeSeparator       string
	MergeIncludeMeta     bool
	MergeGenerateTOC     bool
	MergeTOCDepth        int
	MergeAdjustHeaders   bool
	MergeBaseLevel       int
	MergeIncludePatterns []string
	MergeExcludePatterns []string
	MergeDryRun          bool
}

// parseArgs parses command line arguments and returns CLIArgs
func parseArgs() (*CLIArgs, error) {
	args := os.Args[1:]
	cliArgs := &CLIArgs{
		// Set merge defaults
		MergeOrder:        "filename",
		MergeSeparator:    "\n\n---\n\n",
		MergeGenerateTOC:  true,
		MergeTOCDepth:     3,
		MergeBaseLevel:    2, // Start from H2, H1 reserved for document title
		MergeAdjustHeaders: true, // Default to true for better document structure
	}

	// Handle verbose flag
	if len(args) > 0 && args[0] == "-v" {
		cliArgs.Verbose = true
		args = args[1:]
		if verbose {
			log("Verbose mode enabled")
		}
	}

	if len(args) < 1 {
		return nil, fmt.Errorf("missing required arguments")
	}

	// Check if this is a merge command
	if args[0] == "merge" {
		cliArgs.IsMergeCommand = true
		return parseMergeArgs(cliArgs, args[1:])
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

// parseMergeArgs parses arguments for the merge command
func parseMergeArgs(cliArgs *CLIArgs, args []string) (*CLIArgs, error) {
	cliArgs.IsMergeCommand = true
	
	if len(args) < 1 {
		return nil, fmt.Errorf("merge command requires a directory argument")
	}

	// Parse non-flag arguments
	nonFlagArgs := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		
		if !strings.HasPrefix(arg, "-") {
			nonFlagArgs = append(nonFlagArgs, arg)
			continue
		}

		// Handle flags
		switch arg {
		case "-r", "--recursive":
			cliArgs.MergeRecursive = true
		case "--dry-run":
			cliArgs.MergeDryRun = true
		case "--include-meta":
			cliArgs.MergeIncludeMeta = true
		case "--no-toc":
			cliArgs.MergeGenerateTOC = false
		case "--adjust-headers":
			cliArgs.MergeAdjustHeaders = true
		case "-o", "--output":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("%s requires a value", arg)
			}
			i++
			cliArgs.MergeOutputFile = args[i]
		case "--order":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--order requires a value")
			}
			i++
			if !isValidOrder(args[i]) {
				return nil, fmt.Errorf("invalid order '%s'. Valid orders: filename, modified, size, custom", args[i])
			}
			cliArgs.MergeOrder = args[i]
		case "--separator":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--separator requires a value")
			}
			i++
			cliArgs.MergeSeparator = args[i]
		case "--toc-depth":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--toc-depth requires a value")
			}
			i++
			depth := parseIntOrError(args[i], "--toc-depth")
			if depth < 1 || depth > 6 {
				return nil, fmt.Errorf("--toc-depth must be between 1 and 6")
			}
			cliArgs.MergeTOCDepth = depth
		case "--base-level":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--base-level requires a value")
			}
			i++
			level := parseIntOrError(args[i], "--base-level")
			if level < 1 || level > 6 {
				return nil, fmt.Errorf("--base-level must be between 1 and 6")
			}
			cliArgs.MergeBaseLevel = level
		case "--include":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--include requires a pattern")
			}
			i++
			cliArgs.MergeIncludePatterns = append(cliArgs.MergeIncludePatterns, args[i])
		case "--exclude":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--exclude requires a pattern")
			}
			i++
			cliArgs.MergeExcludePatterns = append(cliArgs.MergeExcludePatterns, args[i])
		default:
			return nil, fmt.Errorf("unknown merge option: %s", arg)
		}
	}

	// Assign non-flag arguments
	if len(nonFlagArgs) < 1 {
		return nil, fmt.Errorf("merge command requires a directory argument")
	}
	
	cliArgs.MergeDirectory = nonFlagArgs[0]
	
	if len(nonFlagArgs) > 1 {
		cliArgs.MergeOutputFile = nonFlagArgs[1]
	} else if cliArgs.MergeOutputFile == "" {
		cliArgs.MergeOutputFile = "merged.md"
	}

	return cliArgs, nil
}

// isValidOrder checks if the order type is valid
func isValidOrder(order string) bool {
	validOrders := []string{"filename", "modified", "size", "custom"}
	for _, valid := range validOrders {
		if order == valid {
			return true
		}
	}
	return false
}

// parseIntOrError parses an integer or returns an error
func parseIntOrError(s, flag string) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	fmt.Fprintf(os.Stderr, "Error: %s requires a valid integer\n", flag)
	os.Exit(1)
	return 0
}

// showUsage displays the usage information
func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: \n")
	fmt.Fprintf(os.Stderr, "  doc [-v] <language_code> [transform_instruction]  # Translation\n")
	fmt.Fprintf(os.Stderr, "  doc [-v] merge <directory> [output_file] [options] # Merge markdown files\n")
	fmt.Fprintf(os.Stderr, "\nTranslation Examples:\n")
	fmt.Fprintf(os.Stderr, "  cat README.md | doc ja\n")
	fmt.Fprintf(os.Stderr, "  cat README.md | doc -v ru\n")
	fmt.Fprintf(os.Stderr, "\nMerge Examples:\n")
	fmt.Fprintf(os.Stderr, "  doc merge ./docs/                    # Merge all .md files to merged.md\n")
	fmt.Fprintf(os.Stderr, "  doc merge ./docs/ book.md            # Merge to book.md\n")
	fmt.Fprintf(os.Stderr, "  doc merge ./docs/ -r --include-meta  # Recursive with metadata\n")
	fmt.Fprintf(os.Stderr, "  doc merge ./docs/ --dry-run          # Preview without merging\n")
	fmt.Fprintf(os.Stderr, "\nMerge Options:\n")
	fmt.Fprintf(os.Stderr, "  -o, --output FILE         Output file (default: merged.md)\n")
	fmt.Fprintf(os.Stderr, "  -r, --recursive           Include subdirectories\n")
	fmt.Fprintf(os.Stderr, "  --order ORDER             Sort order: filename, modified, size, custom (default: filename)\n")
	fmt.Fprintf(os.Stderr, "  --separator STRING        File separator (default: \\n\\n---\\n\\n)\n")
	fmt.Fprintf(os.Stderr, "  --include PATTERN         Include files matching pattern\n")
	fmt.Fprintf(os.Stderr, "  --exclude PATTERN         Exclude files matching pattern\n")
	fmt.Fprintf(os.Stderr, "  --include-meta            Include metadata comments\n")
	fmt.Fprintf(os.Stderr, "  --no-toc                  Disable table of contents\n")
	fmt.Fprintf(os.Stderr, "  --toc-depth N             TOC depth (1-6, default: 3)\n")
	fmt.Fprintf(os.Stderr, "  --adjust-headers          Adjust header levels\n")
	fmt.Fprintf(os.Stderr, "  --base-level N            Base header level (1-6, default: 1)\n")
	fmt.Fprintf(os.Stderr, "  --dry-run                 Preview without writing\n")
	fmt.Fprintf(os.Stderr, "\nGeneral Commands:\n")
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
