# doc - Document Ordering Clerk (D.O.C.)

A simple command-line tool for translating documents while preserving the original format using multiple LLM providers (Claude Code, OpenAI, Anthropic).

## Features

- **Format Preservation**: Completely maintains structures like Markdown, HTML, plain text, etc.
- **Multiple Provider Support**: Claude Code CLI (default), OpenAI API, Anthropic API
- **35+ Language Support**: A wide range of languages including Japanese, English, Chinese, Russian, etc.
- **Intelligent Response Handling**: JSON structured responses, error handling
- **Progress Display**: Animated spinner and elapsed time display
- **Shell Integration**: Full support for UNIX pipelines

## Installation

```bash
# Install from GitHub
go install github.com/bigdra50/doc@latest

# Or build from source
git clone https://github.com/bigdra50/doc.git
cd doc
go build -o doc .
```

## Usage

### Basic Translation

```bash
# Translate from standard input to Japanese
cat document.md | doc ja

# Translate from a file to Russian (with detailed logs)
cat spec.html | doc -v ru

# Display list of supported language codes
doc --list

# Translation with custom instructions
cat technical_doc.md | doc ja "Convert technical specifications to user guide"
```

### Provider Configuration

#### 1. Claude Code CLI (default)

```bash
# Installation (npm required)
npm install -g @anthropic-ai/claude-code

# Example usage
cat document.md | doc ja
```

#### 2. OpenAI API

```bash
# Set API key
export OPENAI_API_KEY=sk-your-openai-api-key
export LLM_PROVIDER=openai

# Example usage
cat document.md | doc ja
```

#### 3. Anthropic API

```bash
# Set API key
export ANTHROPIC_API_KEY=sk-ant-your-anthropic-api-key
export LLM_PROVIDER=anthropic

# Example usage (coming soon)
cat document.md | doc ja
```

## Configuration

### Configuration File

The tool supports persistent configuration via TOML file:

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

Configuration file location follows XDG Base Directory specification:

- `$XDG_CONFIG_HOME/bigdra50/doc/config.toml`
- `~/.config/bigdra50/doc/config.toml` (fallback)

### Environment Variables

Environment variables override configuration file settings:

#### Environment Configuration File

```bash
# Create .env file
echo "LLM_PROVIDER=openai" > .env
echo "OPENAI_API_KEY=sk-your-api-key" >> .env
echo "OPENAI_MODEL=gpt-4o-mini" >> .env
```

### Model Selection

```bash
# List available models
doc --list-models

# Display models by provider
doc --list-models openai
doc --list-models anthropic
```

## Supported Languages

| Code | Language   | Code | Language | Code | Language |
| ---- | ---------- | ---- | -------- | ---- | -------- |
| ja   | Japanese   | en   | English  | ko   | Korean   |
| zh   | Chinese    | ru   | Russian  | es   | Spanish  |
| fr   | French     | de   | German   | it   | Italian  |
| pt   | Portuguese | ar   | Arabic   | hi   | Hindi    |

You can check the complete list with `doc --list`.

## Environment Variables

| Variable Name       | Description        | Default Value               |
| ------------------- | ------------------ | --------------------------- |
| `LLM_PROVIDER`      | Provider selection | `claude-code`               |
| `OPENAI_API_KEY`    | OpenAI API key     | -                           |
| `ANTHROPIC_API_KEY` | Anthropic API key  | -                           |
| `OPENAI_MODEL`      | OpenAI model       | `gpt-4o-mini`               |
| `ANTHROPIC_MODEL`   | Anthropic model    | `claude-3-5-haiku-20241022` |
| `CLAUDE_MODEL`      | Claude Code model  | `sonnet`                    |

## Error Codes

| Code | Description                                         |
| ---- | --------------------------------------------------- |
| 0    | Success                                             |
| 1    | System error (no input, configuration issues, etc.) |
| 2    | Same language (already in target language)          |
| 3    | Translation not possible (code, data, etc.)         |
| 4    | Format error                                        |
| 5    | Content error                                       |

## Example Execution

```bash
# Basic translation
echo "Hello World" | doc ja
# → こんにちは世界

# Same language detection
echo "こんにちは" | doc ja
# → Exit code 2

# Markdown format preservation
echo "# Title\n- List item" | doc ja
# → # タイトル\n- リスト項目

# Progress display
cat large_document.md | doc -v ja
# [INFO] Reading document...
# ⠋ Translating with Claude Code CLI... (2.3s)
# ✓ Translation completed (2.3s)
```

## Development & Testing

```bash
# Run tests
go test

# Build
go build -o doc .

# Cleanup
rm -f doc
```

## Architecture

### Core Components

- **main.go**: CLI processing, application logic
- **provider.go**: LLM provider interface
- **claude_provider.go**: Claude Code CLI implementation
- **openai_provider.go**: OpenAI API implementation
- **models.go**: Model catalog, cost calculation

### Design Principles

1. **Interface-Centric**: Unified LLMProvider abstraction
2. **Environment-Driven**: Support for `.env` files, environment variables
3. **Error Handling**: Detailed exit code system
4. **Format Preservation**: Complete maintenance of original document structure

## License

MIT License

