# doc - Document Ordering Clerk (D.O.C.)

[![Go](https://github.com/bigdra50/doc/actions/workflows/go.yml/badge.svg)](https://github.com/bigdra50/doc/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bigdra50/doc)](https://goreportcard.com/report/github.com/bigdra50/doc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple command-line tool for translating documents while preserving their original format using multiple LLM providers (Claude Code, OpenAI, Anthropic) with intelligent response handling.

## Project Overview

This is a Go-based CLI tool that follows the UNIX philosophy of "do one thing well" - translate documents in any format while maintaining their structure perfectly. Features structured JSON response processing for robust error handling and intelligent translation status detection.

## Installation

```bash
# Install from GitHub
go install github.com/bigdra50/doc@latest

# Or build from source
git clone https://github.com/bigdra50/doc.git
cd doc
go build -o doc .
```

## Build and Test Commands

```bash
# Build the application
go build -o doc .

# Run tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linter
golangci-lint run

# Clean build artifacts
rm -f doc
```

## Usage

```bash
# Basic translation
cat document.md | doc ja

# Translation with verbose logging
cat document.html | doc -v ru

# Show supported language codes
doc --list

# Translation with custom instruction
cat spec.md | doc ja "convert from technical spec to user guide"
```

## LLM Provider System

The tool supports multiple LLM providers through environment variable configuration:

### Supported Providers

1. **Claude Code CLI** (default)
   - Uses installed Claude Code SDK
   - No API key required
   - Command: `claude -p`

2. **OpenAI API**
   - Requires OpenAI API key
   - Direct prompting for translation (no function calling)
   - Default Model: gpt-4o-mini (configurable)

3. **Anthropic Claude API** 
   - Requires Anthropic API key
   - (Not yet implemented - future enhancement)

### Configuration

The tool supports configuration through multiple methods:

#### 1. Configuration File (Recommended)

```bash
# Initialize configuration file
doc --init-config

# Show current configuration
doc --config

# Set configuration values
doc --set provider=openai
doc --set openai_api_key=sk-your-key
doc --set openai_model=gpt-4o
```

Configuration file location (follows XDG Base Directory spec):
- `$XDG_CONFIG_HOME/bigdra50/doc/config.toml`
- `~/.config/bigdra50/doc/config.toml` (fallback)

#### 2. Environment Variables

```bash
# Provider Selection (default: claude-code)
export LLM_PROVIDER=claude-code  # or 'openai' or 'anthropic'

# API Keys (required for respective providers)
export OPENAI_API_KEY=sk-your-openai-api-key
export ANTHROPIC_API_KEY=sk-ant-your-anthropic-api-key

# Model Selection
export OPENAI_MODEL=gpt-4o-mini
export ANTHROPIC_MODEL=claude-3-5-haiku-20241022
export CLAUDE_MODEL=sonnet

# Optional: Custom Claude Code CLI path
export CLAUDE_CODE_PATH=/custom/path/to/claude
```

#### 3. .env File

```bash
# Create .env file in current directory
echo "LLM_PROVIDER=openai" > .env
echo "OPENAI_API_KEY=sk-your-key" >> .env
echo "OPENAI_MODEL=gpt-4o-mini" >> .env
```

Priority: Environment Variables > Config File > .env File > Defaults

### Usage Examples

```bash
# Default Claude Code provider
cat document.md | doc ja

# OpenAI provider
LLM_PROVIDER=openai cat document.md | doc ja

# With environment file
echo "LLM_PROVIDER=openai" > .env
source .env
cat document.md | doc ja
```

## Features

- **Language Code Validation**: Supports 35+ language codes (ja, en, ru, zh, etc.) with prefix-based suggestions
- **Format Preservation**: Maintains original document structure (Markdown, HTML, plain text) exactly
- **Structured Response Processing**: JSON-based communication with Claude for robust error handling
- **Intelligent Status Detection**: Automatically detects same-language documents, untranslatable content, etc.
- **Progress Indication**: Animated spinner with elapsed time during translation
- **Comprehensive Error Handling**: Different exit codes for specific error conditions
- **Verbose Logging**: Debug mode with `-v` flag showing internal processing steps
- **Shell-Friendly**: Proper stdin/stdout handling for UNIX pipelines
- **Configuration Management**: Persistent settings with XDG Base Directory support
- **Model Selection**: Choose from multiple models per provider
- **CI/CD Integration**: GitHub Actions workflows for testing and releases

## Architecture

### Core Files
- `main.go`: Core application logic with provider abstraction
- `provider.go`: LLMProvider interface and configuration management
- `claude_provider.go`: Claude Code CLI implementation
- `openai_provider.go`: OpenAI API implementation
- `anthropic_provider.go`: Anthropic API implementation (placeholder)
- `models.go`: Model catalog with cost information
- `cli.go`: Command-line argument parsing and help
- `language.go`: Language code validation and suggestions
- `translation.go`: Translation orchestration logic
- `ui.go`: Terminal UI components (spinner, logging)
- `internal/config/`: Configuration management with TOML support
- `internal/utils/`: Utility functions

### Key Components
- **LLMProvider Interface**: Unified abstraction for all translation providers
- **TranslationResponse struct**: JSON response format with status, message, content, and error codes
- **ProviderConfig**: Environment-based configuration system
- **Language validation**: Comprehensive validation with prefix-based suggestion system
- **Spinner animation**: Progress indication for long-running translation operations
- **Provider Factory**: Dynamic provider creation based on environment configuration
- **Error Handler**: Status-based routing with specific exit codes

## Error Handling

### System Errors (Exit Code 1)
- Empty stdin input
- Claude command not found  
- Invalid language code with prefix-based suggestions
- Command execution failures
- Malformed JSON responses

### Translation-Specific Errors
- **Exit Code 2**: SAME_LANGUAGE - Document already in target language
- **Exit Code 3**: UNTRANSLATABLE - Content cannot be translated (code/data)
- **Exit Code 4**: FORMAT_ERROR - Document format too complex to preserve
- **Exit Code 5**: CONTENT_ERROR - Document content corrupted or unreadable

### Response Status Types
- `success`: Translation completed successfully
- `no_change_needed`: Document already in target language
- `error`: Translation failed with specific error code

## Dependencies

- Go 1.21+ (tested with 1.21, 1.22, 1.23)
- **Claude Code Provider**: Claude Code CLI (`npm install -g @anthropic-ai/claude-code`)
- **OpenAI Provider**: Valid OPENAI_API_KEY
- **Anthropic Provider**: Valid ANTHROPIC_API_KEY (not yet implemented)
- **External Dependencies**:
  - `github.com/BurntSushi/toml`: TOML configuration file support

## Testing

The tool includes comprehensive tests for:
- Language code validation with edge cases
- JSON response parsing and error handling
- Duration formatting for progress display
- Supported language map integrity
- Structured response processing
- Error code mapping and exit status

### Test Examples
```bash
# Test language code validation
go test -run TestValidateLanguageCode

# Test all functionality
go test -v

# Test translation scenarios
echo "Hello" | ./doc ja     # Basic translation
echo "こんにちは" | ./doc ja  # Same language detection
echo "print('hi')" | ./doc ja # Code handling
```

## CLI Best Practices Implemented

- **Progress Feedback**: Animated Unicode spinner with elapsed time display
- **Terminal Detection**: Automatic fallback to simple logs in non-terminal environments
- **Proper Stream Usage**: Errors to stderr, content to stdout
- **Meaningful Exit Codes**: Different codes for different error types (1-5)
- **User-Friendly Messages**: Clear error descriptions with helpful hints
- **Input Validation**: Language code validation with prefix-based suggestions
- **Robust Parsing**: JSON response parsing with graceful fallback
- **Help System**: `--list` for supported languages, usage examples
- **Performance Optimization**: Concise prompts for faster response times

## JSON Response Format

The tool uses structured JSON communication with Claude:

```json
{
  "status": "success|error|no_change_needed",
  "message": "Human-readable explanation",
  "content": "Actual translated document content",
  "error_code": "SAME_LANGUAGE|UNTRANSLATABLE|FORMAT_ERROR|CONTENT_ERROR"
}
```

This enables intelligent handling of edge cases like same-language documents and untranslatable content.

## Model Selection

### OpenAI Models
```bash
# List available models
doc --list-models openai

# Set model via config
doc --set openai_model=gpt-4o

# Set model via environment
export OPENAI_MODEL=gpt-4o-mini
```

### Available Models by Provider

#### OpenAI
- **gpt-4o**: High capability, multimodal
- **gpt-4o-mini**: Cost-effective, fast (default)
- **gpt-4**: Classic high performance
- **gpt-3.5-turbo**: Fast, lower cost

#### Claude Code
- **opus**: High capability
- **sonnet**: Balanced (default)
- **haiku**: Fast, efficient

#### Anthropic (planned)
- **claude-3-5-sonnet**: High capability
- **claude-3-5-haiku**: Fast, efficient (default)

## CI/CD

The project uses GitHub Actions for continuous integration:

- **Testing**: Multi-version Go testing (1.21, 1.22, 1.23)
- **Linting**: golangci-lint for code quality
- **Cross-platform builds**: Linux, macOS, Windows
- **Automated releases**: GoReleaser on tag push
- **Code coverage**: Automated coverage reports

### Release Process
```bash
# Create and push a tag
git tag v0.1.1
git push origin v0.1.1
```

This triggers automated:
- Binary builds for all platforms
- GitHub release creation
- Changelog generation