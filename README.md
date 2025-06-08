# doc - Document Ordering Clerk

[![Go](https://github.com/bigdra50/doc/actions/workflows/go.yml/badge.svg)](https://github.com/bigdra50/doc/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bigdra50/doc)](https://goreportcard.com/report/github.com/bigdra50/doc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A versatile command-line tool for document operations: **translation** and **markdown file merging** while preserving format and structure using multiple LLM providers.

## Features

### 🌐 Translation

- Translate documents in any format while maintaining structure
- Multiple LLM providers (Claude Code, OpenAI, Anthropic)
- Intelligent response handling with structured JSON processing
- 35+ language codes support with validation

### 📚 Markdown File Merging

- Merge multiple markdown files into a single, well-structured document
- Automatic table of contents generation
- Flexible file ordering (filename, modified date, size, custom)
- Include/exclude pattern filtering
- Header level adjustment for consistent document hierarchy
- Metadata insertion with source tracking

## Installation

```bash
# Install from GitHub
go install github.com/bigdra50/doc@latest

# Or build from source
git clone https://github.com/bigdra50/doc.git
cd doc
go build -o doc .
```

## Quick Start

### Translation

```bash
# Basic translation
cat document.md | doc ja

# With verbose logging
cat document.html | doc -v ru

# Show supported languages
doc --list
```

### Markdown File Merging

```bash
# Basic merge - creates a document with automatic title and TOC
doc merge ./docs/

# Custom output file
doc merge ./docs/ my-book.md

# Preview without writing
doc merge ./docs/ --dry-run
```

## Markdown File Merging - Detailed Usage

### Basic Commands

```bash
# Merge all .md files in current directory
doc merge .

# Merge with custom output file
doc merge ./chapters/ book.md

# Merge with verbose output
doc -v merge ./docs/ guide.md
```

### File Ordering Options

```bash
# Sort by filename (default)
doc merge ./docs/ --order filename

# Sort by modification date (oldest first)
doc merge ./docs/ --order modified

# Sort by file size (smallest first)
doc merge ./docs/ --order size

# Use custom order file (.docorder)
doc merge ./docs/ --order custom
```

### Filtering Options

```bash
# Include only specific patterns
doc merge ./docs/ --include "chapter*.md"

# Exclude specific files
doc merge ./docs/ --exclude "README.md" --exclude "CHANGELOG.md"

# Combine include/exclude (multiple patterns)
doc merge ./docs/ --include "*.md" --exclude "draft_*" --exclude "*_backup.md"

# Recursive directory scanning
doc merge ./project/ -r --include "docs/*.md"
```

### Document Structure Control

```bash
# Disable table of contents
doc merge ./docs/ --no-toc

# Custom TOC depth (1-6 levels)
doc merge ./docs/ --toc-depth 2

# Disable automatic header adjustment (keep original levels)
doc merge ./docs/ --adjust-headers=false

# Custom base header level (useful for embedding in larger documents)
doc merge ./docs/ --base-level 3
```

### Metadata and Formatting

```bash
# Include metadata comments (source files, generation time)
doc merge ./docs/ --include-meta

# Custom file separator
doc merge ./docs/ --separator "\\n\\n***\\n\\n"

# Combine multiple options
doc merge ./docs/ book.md --include-meta --toc-depth 2 --order modified
```

### Advanced Use Cases

#### 📖 Creating a Book from Chapters

```bash
# Organize chapters with proper hierarchy
doc merge ./chapters/ my-book.md \\
  --order filename \\
  --include "chapter*.md" \\
  --include-meta \\
  --toc-depth 3
```

#### 📝 API Documentation

```bash
# Merge API docs with custom structure
doc merge ./api-docs/ api-reference.md \\
  --order custom \\
  --base-level 2 \\
  --separator "\\n\\n---\\n\\n" \\
  --exclude "internal_*"
```

#### 🎓 Course Materials

```bash
# Create course handbook from lessons
doc merge ./lessons/ course-handbook.md \\
  -r \\
  --include "lesson*.md" \\
  --include "exercise*.md" \\
  --order modified \\
  --include-meta
```

#### 📊 Project Documentation

```bash
# Merge project docs excluding drafts
doc merge ./project-docs/ project-guide.md \\
  -r \\
  --exclude "draft_*" \\
  --exclude "*_wip.md" \\
  --exclude "README.md" \\
  --toc-depth 4 \\
  --include-meta
```

#### 🔬 Research Papers Collection

```bash
# Merge research notes by date
doc merge ./research/ research-compilation.md \\
  --order modified \\
  --include "*.md" \\
  --exclude "template*" \\
  --base-level 2 \\
  --separator "\\n\\n---\\n\\n"
```

### Custom Order File (.docorder)

Create a `.docorder` file in your source directory to specify custom ordering:

```
# .docorder example
introduction.md
chapter-01-basics.md
chapter-02-advanced.md
chapter-03-examples.md
appendix.md
references.md
```

Then use:

```bash
doc merge ./docs/ --order custom
```

### Default Behavior

The merge command uses these intelligent defaults:

- **Document Title**: Auto-generated from output filename

  - `book.md` → `# Book`
  - `user-guide.md` → `# User Guide`
  - `api_reference.md` → `# Api Reference`

- **Header Hierarchy**: Automatic adjustment for clean structure

  - Original `# Chapter 1` → `## Chapter 1` (H2)
  - Original `## Section` → `### Section` (H3)
  - And so on...

- **Table of Contents**: Generated at H2 level with 3-level depth
- **File Separator**: Clean `---` dividers between files

### Example Output Structure

```markdown
# My Book

## Table of Contents

- [Chapter 1](#chapter-1)
  - [Introduction](#introduction)
  - [Getting Started](#getting-started)
- [Chapter 2](#chapter-2)
  - [Advanced Topics](#advanced-topics)

<!-- Source: chapter1.md -->

## Chapter 1

### Introduction

Content from the first chapter...

### Getting Started

Step-by-step instructions...

---

<!-- Source: chapter2.md -->

## Chapter 2

### Advanced Topics

Advanced content here...
```

## Translation Usage

### Basic Translation

```bash
# Translate to Japanese
cat document.md | doc ja

# Translate with custom instruction
cat spec.md | doc ja "convert from technical spec to user guide"

# Show all supported languages
doc --list
```

### LLM Provider Configuration

#### Environment Variables

```bash
# Use OpenAI (requires API key)
export LLM_PROVIDER=openai
export OPENAI_API_KEY=sk-your-key
cat document.md | doc ja

# Use Anthropic Claude (requires API key)
export LLM_PROVIDER=anthropic
export ANTHROPIC_API_KEY=sk-ant-your-key
cat document.md | doc ja
```

#### Configuration File

```bash
# Initialize config
doc --init-config

# Set provider and API keys
doc --set provider=openai
doc --set openai_api_key=sk-your-key
doc --set openai_model=gpt-4o

# View current config
doc --config
```

## Build and Test

```bash
# Build
go build -o doc .

# Run tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

## Configuration

### Supported Providers

1. **Claude Code CLI** (default)

   - No API key required
   - Uses installed Claude Code SDK

2. **OpenAI API**

   - Requires `OPENAI_API_KEY`
   - Default model: `gpt-4o-mini`

3. **Anthropic Claude API**
   - Requires `ANTHROPIC_API_KEY`
   - Default model: `claude-3-5-haiku-20241022`

### Configuration Locations

- Config file: `~/.config/bigdra50/doc/config.toml`
- Environment variables (override config file)
- `.env` file in current directory

## Error Handling

### Exit Codes

- **0**: Success
- **1**: General errors (invalid arguments, file not found, etc.)
- **2**: SAME_LANGUAGE - Document already in target language
- **3**: UNTRANSLATABLE - Content cannot be translated
- **4**: FORMAT_ERROR - Document format too complex
- **5**: CONTENT_ERROR - Document content corrupted

## Examples

### Complete Workflow Examples

#### Technical Documentation

```bash
# Create comprehensive technical docs
doc merge ./tech-docs/ technical-guide.md \\
  -r \\
  --include "*.md" \\
  --exclude "draft_*" \\
  --order filename \\
  --include-meta \\
  --toc-depth 4

# Then translate to Japanese
cat technical-guide.md | doc ja > technical-guide-ja.md
```

#### Multi-language Book

```bash
# Merge chapters into book
doc merge ./chapters/ book-en.md --order filename --include-meta

# Translate to multiple languages
cat book-en.md | doc ja > book-ja.md
cat book-en.md | doc zh > book-zh.md
cat book-en.md | doc es > book-es.md
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with Go for performance and reliability
- Supports multiple LLM providers for flexibility
- Follows UNIX philosophy: do one thing well
- Comprehensive test coverage with TDD approach

