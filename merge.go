package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// runMerge executes the merge command
func runMerge(cliArgs *CLIArgs) error {
	if cliArgs.Verbose {
		log("Starting merge operation")
		log("Directory: %s", cliArgs.MergeDirectory)
		log("Output file: %s", cliArgs.MergeOutputFile)
		log("Order: %s", cliArgs.MergeOrder)
		log("Recursive: %v", cliArgs.MergeRecursive)
	}

	// Create file scanner
	scanner := &FileScanner{
		Directory:       cliArgs.MergeDirectory,
		Recursive:       cliArgs.MergeRecursive,
		IncludePatterns: cliArgs.MergeIncludePatterns,
		ExcludePatterns: cliArgs.MergeExcludePatterns,
	}

	// Scan for markdown files
	log("Scanning directory: %s", cliArgs.MergeDirectory)
	files, err := scanner.ScanMarkdownFiles()
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no markdown files found in directory: %s", cliArgs.MergeDirectory)
	}

	log("Found %d markdown files", len(files))

	// Sort files
	sortedFiles := SortMarkdownFiles(files, cliArgs.MergeOrder)

	if cliArgs.Verbose {
		log("Files to merge (in order):")
		for i, file := range sortedFiles {
			relPath, _ := filepath.Rel(cliArgs.MergeDirectory, file.Path)
			log("  %d. %s (%d bytes)", i+1, relPath, file.Size)
		}
	}

	// Dry run mode
	if cliArgs.MergeDryRun {
		return runDryMode(cliArgs, sortedFiles)
	}

	// Merge files
	return mergeFiles(cliArgs, sortedFiles)
}

// runDryMode shows what would be merged without actually doing it
func runDryMode(cliArgs *CLIArgs, files []MarkdownFile) error {
	fmt.Printf("[DRY RUN] Would process the following files:\n")
	
	totalSize := int64(0)
	for i, file := range files {
		relPath, _ := filepath.Rel(cliArgs.MergeDirectory, file.Path)
		size := formatFileSize(file.Size)
		fmt.Printf("  %d. %s (%s)\n", i+1, relPath, size)
		totalSize += file.Size
	}
	
	fmt.Printf("[DRY RUN] Output file: %s\n", cliArgs.MergeOutputFile)
	fmt.Printf("[DRY RUN] Total size: %s\n", formatFileSize(totalSize))
	
	return nil
}

// mergeFiles merges the markdown files into a single output file
func mergeFiles(cliArgs *CLIArgs, files []MarkdownFile) error {
	// Create output file
	outputFile, err := os.Create(cliArgs.MergeOutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Start progress indication
	spinner := NewSpinner(fmt.Sprintf("Merging files... (0/%d)", len(files)))
	spinner.Start()

	// Write document title and metadata
	if err := writeDocumentHeader(outputFile, cliArgs, files); err != nil {
		spinner.Stop("Merge failed")
		return fmt.Errorf("failed to write document header: %w", err)
	}

	// Write table of contents if requested
	if cliArgs.MergeGenerateTOC {
		if err := writeTOC(outputFile, cliArgs, files); err != nil {
			spinner.Stop("Merge failed")
			return fmt.Errorf("failed to write table of contents: %w", err)
		}
	}

	// Merge files
	for i, file := range files {
		spinner.Stop("")
		spinner = NewSpinner(fmt.Sprintf("Processing files... (%d/%d) - %s", i+1, len(files), file.Name))
		spinner.Start()

		if err := mergeFile(outputFile, file, cliArgs); err != nil {
			spinner.Stop("Merge failed")
			return fmt.Errorf("failed to merge file %s: %w", file.Name, err)
		}

		// Add separator between files (except for the last one)
		if i < len(files)-1 {
			if _, err := outputFile.WriteString(cliArgs.MergeSeparator); err != nil {
				spinner.Stop("Merge failed")
				return fmt.Errorf("failed to write separator: %w", err)
			}
		}
	}
	
	// Calculate total size
	stat, err := os.Stat(cliArgs.MergeOutputFile)
	if err != nil {
		spinner.Stop("Merge failed")
		return fmt.Errorf("failed to get output file stats: %w", err)
	}

	finalMessage := fmt.Sprintf("Merge completed - Output: %s (%s)", cliArgs.MergeOutputFile, formatFileSize(stat.Size()))
	spinner.Stop(finalMessage)
	
	return nil
}

// writeDocumentHeader writes the document title and optional metadata
func writeDocumentHeader(file *os.File, cliArgs *CLIArgs, files []MarkdownFile) error {
	// Generate document title from output filename
	title := generateDocumentTitle(cliArgs.MergeOutputFile)
	
	// Write document title (H1)
	if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", title)); err != nil {
		return err
	}
	
	// Write metadata if requested
	if cliArgs.MergeIncludeMeta {
		header := fmt.Sprintf(`<!-- Generated by doc merge at %s -->
<!-- Source directory: %s -->
<!-- Files merged: %d -->
<!-- Command: doc merge %s -->

`, time.Now().Format("2006-01-02 15:04:05"), cliArgs.MergeDirectory, len(files), cliArgs.MergeDirectory)
		
		if _, err := file.WriteString(header); err != nil {
			return err
		}
	}
	
	return nil
}

// generateDocumentTitle creates a document title from the output filename
func generateDocumentTitle(outputFile string) string {
	// Extract filename without extension
	base := filepath.Base(outputFile)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	
	// Convert to title case
	if name == "merged" {
		return "Document"
	}
	
	// Replace underscores and hyphens with spaces, then title case
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	
	// Simple title case conversion
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	
	return strings.Join(words, " ")
}

// writeTOC writes the table of contents to the output file
func writeTOC(file *os.File, cliArgs *CLIArgs, files []MarkdownFile) error {
	_, err := file.WriteString("## Table of Contents\n\n")
	if err != nil {
		return err
	}

	for _, markdownFile := range files {
		// Read file to extract headers
		content, err := os.ReadFile(markdownFile.Path)
		if err != nil {
			continue
		}

		headers := extractHeaders(string(content), cliArgs.MergeTOCDepth)
		for _, header := range headers {
			// Adjust header level for TOC (since file headers will be adjusted)
			adjustedLevel := header.Level + cliArgs.MergeBaseLevel - 1
			if adjustedLevel > cliArgs.MergeTOCDepth + 1 { // +1 for the document title level
				continue
			}
			
			indent := strings.Repeat("  ", adjustedLevel-2) // -2 because TOC starts at level 2
			link := strings.ToLower(strings.ReplaceAll(header.Text, " ", "-"))
			// Remove non-alphanumeric characters from link
			link = strings.Map(func(r rune) rune {
				if r == ' ' || r == '-' {
					return '-'
				}
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
					return r
				}
				return -1
			}, link)
			
			_, err := file.WriteString(fmt.Sprintf("%s- [%s](#%s)\n", indent, header.Text, link))
			if err != nil {
				return err
			}
		}
	}

	_, err = file.WriteString("\n")
	return err
}

// mergeFile merges a single markdown file into the output
func mergeFile(outputFile *os.File, file MarkdownFile, cliArgs *CLIArgs) error {
	// Write file source comment if metadata is enabled
	if cliArgs.MergeIncludeMeta {
		relPath, _ := filepath.Rel(cliArgs.MergeDirectory, file.Path)
		comment := fmt.Sprintf("<!-- Source: %s -->\n", relPath)
		if _, err := outputFile.WriteString(comment); err != nil {
			return err
		}
	}

	// Read the file content
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fileContent := string(content)

	// Adjust header levels if requested
	if cliArgs.MergeAdjustHeaders {
		fileContent = adjustHeaderLevels(fileContent, cliArgs.MergeBaseLevel)
	}

	// Write the content
	if _, err := outputFile.WriteString(fileContent); err != nil {
		return err
	}

	// Ensure content ends with newline
	if !strings.HasSuffix(fileContent, "\n") {
		if _, err := outputFile.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}

// Header represents a markdown header
type Header struct {
	Level int
	Text  string
}

// extractHeaders extracts headers from markdown content up to maxDepth
func extractHeaders(content string, maxDepth int) []Header {
	var headers []Header
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			level := 0
			for i, char := range line {
				if char == '#' {
					level++
				} else {
					if level > 0 && level <= maxDepth {
						text := strings.TrimSpace(line[i:])
						headers = append(headers, Header{Level: level, Text: text})
					}
					break
				}
			}
		}
	}

	return headers
}

// adjustHeaderLevels adjusts header levels in markdown content
func adjustHeaderLevels(content string, baseLevel int) string {
	lines := strings.Split(content, "\n")
	
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			// Count existing header level
			level := 0
			for _, char := range line {
				if char == '#' {
					level++
				} else {
					break
				}
			}
			
			if level > 0 {
				// Calculate new level
				newLevel := baseLevel + level - 1
				if newLevel > 6 {
					newLevel = 6 // Markdown only supports up to 6 levels
				}
				
				// Replace with new header level
				headerPrefix := strings.Repeat("#", newLevel)
				headerText := strings.TrimSpace(line[level:])
				lines[i] = headerPrefix + " " + headerText
			}
		}
	}
	
	return strings.Join(lines, "\n")
}

// formatFileSize formats file size in human-readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}