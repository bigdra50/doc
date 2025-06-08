# doc - Document Ordering Clerk (D.O.C.)

A simple command-line tool for translating documents while preserving their original format using multiple LLM providers (Claude Code, OpenAI, Anthropic) with intelligent response handling.

## Project Overview

This is a Go-based CLI tool that follows the UNIX philosophy of "do one thing well" - translate documents in any format while maintaining their structure perfectly. Features structured JSON response processing for robust error handling and intelligent translation status detection.

## Build and Test Commands

```bash
# Build the application
go build -o doc .

# Run tests
go test

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
   - Uses function calling for structured translation
   - Model: GPT-4

3. **Anthropic Claude API** 
   - Requires Anthropic API key
   - (Not yet implemented - future enhancement)

### Environment Configuration

```bash
# Provider Selection (default: claude-code)
export LLM_PROVIDER=claude-code  # or 'openai' or 'anthropic'

# API Keys (required for respective providers)
export OPENAI_API_KEY=sk-your-openai-api-key
export ANTHROPIC_API_KEY=sk-ant-your-anthropic-api-key

# Optional: Custom Claude Code CLI path
export CLAUDE_CODE_PATH=/custom/path/to/claude
```

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

## Architecture

### Core Files
- `main.go`: Core application logic with provider abstraction
- `provider.go`: LLMProvider interface and configuration management
- `claude_provider.go`: Claude Code CLI implementation
- `openai_provider.go`: OpenAI API implementation with function calling
- `anthropic_provider.go`: Anthropic API implementation (placeholder)
- `main_test.go`: Comprehensive unit tests for all functionality

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

- Go 1.19+
- **Claude Code Provider**: Claude Code SDK (pre-installed)
- **OpenAI Provider**: Valid OPENAI_API_KEY environment variable
- **Anthropic Provider**: Valid ANTHROPIC_API_KEY environment variable (not yet implemented)

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