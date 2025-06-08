package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MarkdownFile represents a markdown file with metadata
type MarkdownFile struct {
	Path    string
	Name    string
	ModTime time.Time
	Size    int64
}

// FileScanner handles scanning directories for markdown files
type FileScanner struct {
	Directory       string
	Recursive       bool
	IncludePatterns []string
	ExcludePatterns []string
}

// ScanMarkdownFiles scans the directory and returns markdown files
func (fs *FileScanner) ScanMarkdownFiles() ([]MarkdownFile, error) {
	// Check if directory exists
	if _, err := os.Stat(fs.Directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", fs.Directory)
	}

	var files []MarkdownFile

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// If not recursive, skip subdirectories
			if !fs.Recursive && path != fs.Directory {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if it's a markdown file
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return nil
		}

		// Apply include patterns
		if len(fs.IncludePatterns) > 0 {
			matched := false
			for _, pattern := range fs.IncludePatterns {
				if matchPattern(info.Name(), pattern) {
					matched = true
					break
				}
			}
			if !matched {
				return nil
			}
		}

		// Apply exclude patterns
		for _, pattern := range fs.ExcludePatterns {
			if matchPattern(info.Name(), pattern) {
				return nil
			}
		}

		files = append(files, MarkdownFile{
			Path:    path,
			Name:    info.Name(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})

		return nil
	}

	if err := filepath.Walk(fs.Directory, walkFunc); err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return files, nil
}

// SortMarkdownFiles sorts markdown files based on the specified order
func SortMarkdownFiles(files []MarkdownFile, order string) []MarkdownFile {
	sorted := make([]MarkdownFile, len(files))
	copy(sorted, files)

	switch order {
	case "filename":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Name < sorted[j].Name
		})
	case "modified":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].ModTime.Before(sorted[j].ModTime)
		})
	case "size":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Size < sorted[j].Size
		})
	case "custom":
		// TODO: Implement custom ordering based on .docorder file
		// For now, fallback to filename
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Name < sorted[j].Name
		})
	default:
		// Default to filename ordering
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Name < sorted[j].Name
		})
	}

	return sorted
}

// matchPattern matches a filename against a pattern
func matchPattern(filename, pattern string) bool {
	matched, err := filepath.Match(pattern, filename)
	if err != nil {
		// If pattern is invalid, treat as non-match
		return false
	}
	return matched
}